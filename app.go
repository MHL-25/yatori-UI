package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
	"yatori-UI/internal/activity"
	"yatori-UI/internal/checker"
	"yatori-UI/internal/config"
	"yatori-UI/internal/dao"
	"yatori-UI/internal/entity"
	"yatori-UI/internal/monitor"
	"yatori-UI/internal/service"

	"github.com/fsnotify/fsnotify"
	systray "github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/yatori-dev/yatori-go-core/utils"
)

//go:embed build/windows/icon.ico
var trayIconData []byte

type App struct {
	ctx          context.Context
	notified     map[string]time.Time
	notifiedMu   sync.Mutex
	shouldQuit   bool
	donatePrompt bool
	launchCount  int
}

const AppVersion = "1.2.1"

func NewApp() *App {
	return &App{notified: make(map[string]time.Time)}
}

func (a *App) GetAppVersion() string {
	return AppVersion
}

func (a *App) beforeClose(ctx context.Context) bool {
	if a.shouldQuit {
		return false
	}
	runtime.WindowHide(ctx)
	return true
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		os.Chdir(exeDir)
		dbPath := filepath.Join(exeDir, "yatori.db")
		if err := dao.InitDB(dbPath); err != nil {
			log.Printf("数据库初始化失败: %v", err)
		}
	} else {
		log.Printf("获取可执行文件路径失败: %v", err)
		if err := dao.InitDB("yatori.db"); err != nil {
			log.Printf("数据库初始化失败(回退): %v", err)
		}
	}

	utils.YatoriCoreInit()

	a.recordLaunch()

	go a.initSystray()

	monitor.GlobalEventBus.AddListener(func(eventType string, data interface{}) {
		jsonData, _ := json.Marshal(map[string]interface{}{
			"event": eventType,
			"data":  data,
		})
		runtime.EventsEmit(ctx, "monitor:update", string(jsonData))

		if eventType == "progress_update" {
			if progress, ok := data.(*monitor.TaskProgress); ok {
				name := progress.AccountName
				if name == "" {
					name = progress.Uid
				}
				key := progress.Uid + "_" + string(progress.Status)
				a.notifiedMu.Lock()
				lastTime, exists := a.notified[key]
				if exists && time.Since(lastTime) < 5*time.Second {
					a.notifiedMu.Unlock()
					return
				}
				a.notified[key] = time.Now()
				a.notifiedMu.Unlock()

				switch progress.Status {
				case monitor.StatusCompleted:
					runtime.EventsEmit(ctx, "notification", map[string]interface{}{
						"type":    "success",
						"title":   "任务完成",
						"message": name + " 的刷课任务已完成！",
					})
				case monitor.StatusPaused:
					runtime.EventsEmit(ctx, "notification", map[string]interface{}{
						"type":    "success",
						"title":   "已暂停",
						"message": name + " 的任务已暂停",
					})
				case monitor.StatusStopped:
					runtime.EventsEmit(ctx, "notification", map[string]interface{}{
						"type":    "success",
						"title":   "已停止",
						"message": name + " 的任务已停止",
					})
				case monitor.StatusError:
					runtime.EventsEmit(ctx, "notification", map[string]interface{}{
						"type":    "error",
						"title":   "任务异常",
						"message": name + " 出现错误：" + progress.ErrorMessage,
					})
				}
			}
		}
	})

	go a.watchConfigFile()

	// 启动自动化：按账号自检时间窗决定是否质检/续刷
	go a.startupAutomation()
}

// startupAutomation 启动后按账号自检时间窗执行质检与自动刷课
func (a *App) startupAutomation() {
	time.Sleep(3 * time.Second)
	cfg := config.EnsureConfigFile("./config.yaml")
	users, err := dao.QueryAllUsers()
	if err != nil {
		return
	}
	autoCheck := cfg.Setting.BasicSetting.AutoCheckSw == 1
	autoStart := cfg.Setting.BasicSetting.AutoStartSw == 1
	if !autoCheck && !autoStart {
		return
	}
	for _, u := range users {
		var userCfg config.User
		_ = json.Unmarshal([]byte(u.UserConfigJson), &userCfg)
		ok, reason := config.InSelfCheckWindow(userCfg, time.Now())
		if !ok {
			monitor.GlobalEventBus.AddLog(u.Uid, "⏰ 跳过自动任务："+reason)
			continue
		}
		monitor.GlobalEventBus.AddLog(u.Uid, "⏰ 自检时间窗校验通过："+reason)

		var report *entity.CheckReport
		if autoCheck {
			report, err = checker.Check(u.Uid, cfg.Setting)
			if err != nil {
				monitor.GlobalEventBus.AddLog(u.Uid, "❌ 自动质检失败: "+err.Error())
				continue
			}
		}
		if !autoStart {
			continue
		}
		if report != nil && report.PendingCourses == 0 && report.CompletedCourses == report.TotalCourses {
			monitor.GlobalEventBus.AddLog(u.Uid, "✅ 所有课程已完成，无需自动刷课")
			continue
		}
		if err := service.StartBrush(u.Uid); err != nil {
			monitor.GlobalEventBus.AddLog(u.Uid, "自动刷课未启动: "+err.Error())
		}
	}
	go a.watchCourseLocks()
}

// watchCourseLocks 每分钟检查定时开放节点，到点自动续刷对应账号
func (a *App) watchCourseLocks() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			locks, err := dao.QueryDueCourseLocks(time.Now())
			if err != nil || len(locks) == 0 {
				continue
			}
			for _, lock := range locks {
				user, err := dao.QueryUserByUid(lock.Uid)
				if err != nil || user == nil {
					_ = dao.MarkCourseLockHandled(lock.ID)
					continue
				}
				var userCfg config.User
				_ = json.Unmarshal([]byte(user.UserConfigJson), &userCfg)
				if ok, reason := config.InSelfCheckWindow(userCfg, time.Now()); !ok {
					monitor.GlobalEventBus.AddLog(lock.Uid, "⏰ 定时开放已到达，但不在自检时间窗："+reason)
					continue
				}
				_ = dao.MarkCourseLockHandled(lock.ID)
				if act := activity.GetActivity(lock.Uid); act != nil && act.IsRunning() {
					continue
				}
				monitor.GlobalEventBus.AddLog(lock.Uid, "🕒 课程 ["+lock.CourseName+"] 节点 ["+lock.Label+" "+lock.NodeName+"] 已到开放时间，自动续刷")
				if err := service.StartBrush(lock.Uid); err != nil {
					monitor.GlobalEventBus.AddLog(lock.Uid, "自动续刷未启动: "+err.Error())
				}
			}
		}
	}
}

func (a *App) CheckAccount(uid string) string {
	report, err := service.CheckAccount(uid)
	if err != nil {
		jsonData, _ := json.Marshal(entity.ErrorResponse(err.Error()))
		return string(jsonData)
	}
	jsonData, _ := json.Marshal(entity.SuccessResponse("OK", report))
	return string(jsonData)
}

func (a *App) CheckAccounts(jsonUids string) string {
	var uids []string
	if err := json.Unmarshal([]byte(jsonUids), &uids); err != nil {
		jsonData, _ := json.Marshal(entity.ErrorResponse("Parse error: "+err.Error()))
		return string(jsonData)
	}
	reports := service.CheckAccounts(uids)
	jsonData, _ := json.Marshal(entity.SuccessResponse("OK", reports))
	return string(jsonData)
}

func (a *App) CheckAllAccounts() string {
	reports := service.CheckAllAccounts()
	jsonData, _ := json.Marshal(entity.SuccessResponse("OK", reports))
	return string(jsonData)
}

func (a *App) GetCheckReports() string {
	reports := service.GetStoredCheckReports()
	jsonData, _ := json.Marshal(entity.SuccessResponse("OK", reports))
	return string(jsonData)
}

func (a *App) GetCourseLocks(uid string) string {
	locks, err := service.GetCourseLocks(uid)
	if err != nil {
		jsonData, _ := json.Marshal(entity.ErrorResponse(err.Error()))
		return string(jsonData)
	}
	jsonData, _ := json.Marshal(entity.SuccessResponse("OK", locks))
	return string(jsonData)
}

func (a *App) GetCourseList(uid string) string {
	courses, err := service.GetCourseList(uid)
	if err != nil {
		jsonData, _ := json.Marshal(entity.ErrorResponse(err.Error()))
		return string(jsonData)
	}
	jsonData, _ := json.Marshal(entity.SuccessResponse("OK", courses))
	return string(jsonData)
}

// QueryCourseList 使用表单账号密码临时查询课程，不保存账号
func (a *App) QueryCourseList(jsonReq string) string {
	var req entity.CourseQueryRequest
	if err := json.Unmarshal([]byte(jsonReq), &req); err != nil {
		jsonData, _ := json.Marshal(entity.ErrorResponse("Request parse error: "+err.Error()))
		return string(jsonData)
	}
	courses, err := checker.CourseListByCredentials(req)
	if err != nil {
		jsonData, _ := json.Marshal(entity.ErrorResponse(err.Error()))
		return string(jsonData)
	}
	jsonData, _ := json.Marshal(entity.SuccessResponse("OK", courses))
	return string(jsonData)
}

// recordLaunch 记录软件打开次数，并按伪随机概率决定是否弹出赞助提示
func (a *App) recordLaunch() {
	count, _ := strconv.Atoi(dao.GetSetting("launch_count"))
	count++
	a.launchCount = count
	_ = dao.SetSetting("launch_count", strconv.Itoa(count))

	if dao.GetSetting("sponsor_dismissed") == "true" {
		return
	}
	last, _ := strconv.Atoi(dao.GetSetting("sponsor_last_prompt"))
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if count >= 5 && count-last >= 5 && r.Intn(100) < 20 {
		a.donatePrompt = true
		_ = dao.SetSetting("sponsor_last_prompt", strconv.Itoa(count))
	}
}

// GetDonatePrompt 返回本次启动是否展示赞助提示
func (a *App) GetDonatePrompt() string {
	resp := entity.SuccessResponse("OK", map[string]interface{}{
		"show":  a.donatePrompt,
		"count": a.launchCount,
	})
	jsonData, _ := json.Marshal(resp)
	return string(jsonData)
}

// DisableDonatePrompt 用户选择“不再提示”
func (a *App) DisableDonatePrompt() {
	_ = dao.SetSetting("sponsor_dismissed", "true")
	a.donatePrompt = false
}

func (a *App) initSystray() {
	systray.Run(func() {
		systray.SetIcon(trayIconData)
		systray.SetTooltip("Yatori-UI - 智能网课助手")

		mShow := systray.AddMenuItem("显示主窗口", "显示应用主窗口")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("退出程序", "退出 Yatori-UI")

		go func() {
			for {
				select {
				case <-mShow.ClickedCh:
					runtime.WindowShow(a.ctx)
				case <-mQuit.ClickedCh:
					a.shouldQuit = true
					systray.Quit()
					runtime.Quit(a.ctx)
				}
			}
		}()
	}, func() {
	})
}

func (a *App) watchConfigFile() {
	configPath := "config.yaml"
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return
	}
	configDir := filepath.Dir(absPath)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()

	watcher.Add(configDir)

	var lastReload time.Time
	debounce := 2 * time.Second

	for {
		select {
		case <-a.ctx.Done():
			return
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Name != absPath {
				continue
			}
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				if time.Since(lastReload) < debounce {
					continue
				}
				lastReload = time.Now()
				time.Sleep(500 * time.Millisecond)

				cfg := config.ReadConfig(configPath)
				for i := range cfg.Users {
					user := cfg.Users[i]
					existing, _ := dao.QueryUserByAccount(user.Account, user.AccountType)
					if existing == nil {
						ccData := &entity.CoursesCustomData{
						VideoModel:      user.CoursesCustom.VideoModel,
						AutoExam:        user.CoursesCustom.AutoExam,
						ExamAutoSubmit:  user.CoursesCustom.ExamAutoSubmit,
						ShuffleSw:       user.CoursesCustom.ShuffleSw,
						StudyTime:       user.CoursesCustom.StudyTime,
						IncludeCourses:  user.CoursesCustom.IncludeCourses,
						ExcludeCourses:  user.CoursesCustom.ExcludeCourses,
					}
					if user.CoursesCustom.CxNode != nil {
						ccData.CxNode = *user.CoursesCustom.CxNode
					}
					if user.CoursesCustom.CxChapterTestSw != nil {
						ccData.CxChapterTestSw = *user.CoursesCustom.CxChapterTestSw
					}
					if user.CoursesCustom.CxWorkSw != nil {
						ccData.CxWorkSw = *user.CoursesCustom.CxWorkSw
					}
					if user.CoursesCustom.CxExamSw != nil {
						ccData.CxExamSw = *user.CoursesCustom.CxExamSw
					}
					ccData.CxExtendReadSw = user.CoursesCustom.CxExtendReadSw
					ccData.CxBbsSw = user.CoursesCustom.CxBbsSw
					ccData.CxWorkRetake = user.CoursesCustom.CxWorkRetake
					ccData.CxRetakeRate = user.CoursesCustom.CxRetakeRate
						req := entity.AddAccountRequest{
							AccountType:   user.AccountType,
							Url:           user.URL,
							RemarkName:    user.RemarkName,
							Account:       user.Account,
							Password:      user.Password,
							ShowDocSw:     user.ShowDocSw,
							ShowDocUrls:   user.ShowDocUrls,
							CoursesCustom: ccData,
						}
						service.AddAccount(req)
					}
				}
				runtime.EventsEmit(a.ctx, "config:reloaded", nil)
			}
		case _, ok := <-watcher.Errors:
			if !ok {
				return
			}
		}
	}
}

func (a *App) GetAccountList() string {
	users, err := service.GetAccountList()
	if err != nil {
		resp := entity.ErrorResponse(err.Error())
		jsonData, _ := json.Marshal(resp)
		return string(jsonData)
	}
	resp := entity.SuccessResponse("OK", users)
	jsonData, _ := json.Marshal(resp)
	return string(jsonData)
}

func (a *App) AddAccount(jsonReq string) string {
	var req entity.AddAccountRequest
	if err := json.Unmarshal([]byte(jsonReq), &req); err != nil {
		resp := entity.ErrorResponse("Request parse error: " + err.Error())
		jsonData, _ := json.Marshal(resp)
		return string(jsonData)
	}
	user, err := service.AddAccount(req)
	if err != nil {
		resp := entity.ErrorResponse(err.Error())
		jsonData, _ := json.Marshal(resp)
		return string(jsonData)
	}
	resp := entity.SuccessResponse("OK", user)
	jsonData, _ := json.Marshal(resp)
	return string(jsonData)
}

func (a *App) DeleteAccount(uid string) string {
	err := service.DeleteAccount(uid)
	if err != nil {
		resp := entity.ErrorResponse(err.Error())
		jsonData, _ := json.Marshal(resp)
		return string(jsonData)
	}
	resp := entity.SuccessResponse("OK", nil)
	jsonData, _ := json.Marshal(resp)
	return string(jsonData)
}

func (a *App) UpdateAccount(jsonReq string) string {
	var req entity.UpdateAccountRequest
	if err := json.Unmarshal([]byte(jsonReq), &req); err != nil {
		resp := entity.ErrorResponse("Request parse error: " + err.Error())
		jsonData, _ := json.Marshal(resp)
		return string(jsonData)
	}
	err := service.UpdateAccount(req)
	if err != nil {
		resp := entity.ErrorResponse(err.Error())
		jsonData, _ := json.Marshal(resp)
		return string(jsonData)
	}
	resp := entity.SuccessResponse("OK", nil)
	jsonData, _ := json.Marshal(resp)
	return string(jsonData)
}

func (a *App) LoginCheck(uid string) string {
	err := service.LoginCheck(uid)
	if err != nil {
		resp := entity.ErrorResponse(err.Error())
		jsonData, _ := json.Marshal(resp)
		return string(jsonData)
	}
	resp := entity.SuccessResponse("OK", nil)
	jsonData, _ := json.Marshal(resp)
	return string(jsonData)
}

func (a *App) StartBrush(uid string) string {
	err := service.StartBrush(uid)
	if err != nil {
		resp := entity.ErrorResponse(err.Error())
		jsonData, _ := json.Marshal(resp)
		return string(jsonData)
	}
	resp := entity.SuccessResponse("OK", nil)
	jsonData, _ := json.Marshal(resp)
	return string(jsonData)
}

func (a *App) StopBrush(uid string) string {
	err := service.StopBrush(uid)
	if err != nil {
		resp := entity.ErrorResponse(err.Error())
		jsonData, _ := json.Marshal(resp)
		return string(jsonData)
	}
	resp := entity.SuccessResponse("OK", nil)
	jsonData, _ := json.Marshal(resp)
	return string(jsonData)
}

func (a *App) PauseBrush(uid string) string {
	err := service.PauseBrush(uid)
	if err != nil {
		resp := entity.ErrorResponse(err.Error())
		jsonData, _ := json.Marshal(resp)
		return string(jsonData)
	}
	resp := entity.SuccessResponse("OK", nil)
	jsonData, _ := json.Marshal(resp)
	return string(jsonData)
}

func (a *App) StartAllBrush() string {
	success, total := service.StartAllBrush()
	resp := entity.SuccessResponse("OK", map[string]interface{}{
		"success": success,
		"total":   total,
	})
	jsonData, _ := json.Marshal(resp)
	return string(jsonData)
}

func (a *App) StartBatchBrush(jsonUids string) string {
	var uids []string
	if err := json.Unmarshal([]byte(jsonUids), &uids); err != nil {
		resp := entity.ErrorResponse("Parse error: " + err.Error())
		jsonData, _ := json.Marshal(resp)
		return string(jsonData)
	}
	success, total := service.StartBatchBrush(uids)
	resp := entity.SuccessResponse("OK", map[string]interface{}{
		"success": success,
		"total":   total,
	})
	jsonData, _ := json.Marshal(resp)
	return string(jsonData)
}

func (a *App) GetAllProgress() string {
	progress := service.GetAllProgress()
	resp := entity.SuccessResponse("OK", progress)
	jsonData, _ := json.Marshal(resp)
	return string(jsonData)
}

func (a *App) GetProgress(uid string) string {
	progress := service.GetProgress(uid)
	resp := entity.SuccessResponse("OK", progress)
	jsonData, _ := json.Marshal(resp)
	return string(jsonData)
}

func (a *App) GetConfig() string {
	cfg := service.GetConfig()
	resp := entity.SuccessResponse("OK", cfg)
	jsonData, _ := json.Marshal(resp)
	return string(jsonData)
}

func (a *App) SaveConfig(jsonCfg string) string {
	var cfg config.JSONDataForConfig
	if err := json.Unmarshal([]byte(jsonCfg), &cfg); err != nil {
		resp := entity.ErrorResponse("Config parse error: " + err.Error())
		jsonData, _ := json.Marshal(resp)
		return string(jsonData)
	}
	err := service.SaveConfig(&cfg)
	if err != nil {
		resp := entity.ErrorResponse("Save failed: " + err.Error())
		jsonData, _ := json.Marshal(resp)
		return string(jsonData)
	}
	resp := entity.SuccessResponse("OK", nil)
	jsonData, _ := json.Marshal(resp)
	return string(jsonData)
}

func (a *App) GetPlatforms() string {
	platforms := service.GetPlatforms()
	resp := entity.SuccessResponse("OK", platforms)
	jsonData, _ := json.Marshal(resp)
	return string(jsonData)
}

func (a *App) GetAiTypes() string {
	types := service.GetAiTypes()
	resp := entity.SuccessResponse("OK", types)
	jsonData, _ := json.Marshal(resp)
	return string(jsonData)
}

func (a *App) MinimizeWindow() {
	runtime.WindowMinimise(a.ctx)
}

func (a *App) ToggleMaximizeWindow() {
	if runtime.WindowIsMaximised(a.ctx) {
		runtime.WindowUnmaximise(a.ctx)
	} else {
		runtime.WindowMaximise(a.ctx)
	}
}

func (a *App) CloseWindow() {
	a.shouldQuit = false
	runtime.WindowHide(a.ctx)
}

func (a *App) QuitApp() {
	a.shouldQuit = true
	runtime.Quit(a.ctx)
}

func (a *App) ShowWindow() {
	runtime.WindowShow(a.ctx)
}

func (a *App) IsDisclaimerAccepted() bool {
	return dao.GetSetting("disclaimer_accepted") == "true"
}

func (a *App) AcceptDisclaimer() {
	dao.SetSetting("disclaimer_accepted", "true")
}

func (a *App) ImportConfigFile() string {
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "导入配置文件",
		Filters: []runtime.FileFilter{
			{DisplayName: "配置文件 (*.yaml, *.yml, *.json)", Pattern: "*.yaml;*.yml;*.json"},
			{DisplayName: "YAML文件 (*.yaml, *.yml)", Pattern: "*.yaml;*.yml"},
			{DisplayName: "JSON文件 (*.json)", Pattern: "*.json"},
			{DisplayName: "所有文件 (*.*)", Pattern: "*.*"},
		},
	})
	if err != nil || filePath == "" {
		resp := entity.ErrorResponse("未选择文件")
		jsonData, _ := json.Marshal(resp)
		return string(jsonData)
	}
	count, err := service.ImportAndMergeConfig(filePath)
	if err != nil {
		resp := entity.ErrorResponse("导入失败: " + err.Error())
		jsonData, _ := json.Marshal(resp)
		return string(jsonData)
	}
	resp := entity.SuccessResponse("OK", map[string]interface{}{
		"count":    count,
		"filePath": filePath,
	})
	jsonData, _ := json.Marshal(resp)
	return string(jsonData)
}

func (a *App) OpenFileDialog(title string, filters string) string {
	filterList := []runtime.FileFilter{
		{DisplayName: "所有文件 (*.*)", Pattern: "*.*"},
	}
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   title,
		Filters: filterList,
	})
	if err != nil || filePath == "" {
		return ""
	}
	return filePath
}
