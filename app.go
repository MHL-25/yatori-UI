package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"path/filepath"
	"sync"
	"time"
	"yatori-UI/internal/config"
	"yatori-UI/internal/dao"
	"yatori-UI/internal/entity"
	"yatori-UI/internal/monitor"
	"yatori-UI/internal/service"

	"github.com/fsnotify/fsnotify"
	systray "github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed build/windows/icon.ico
var trayIconData []byte

type App struct {
	ctx        context.Context
	notified   map[string]time.Time
	notifiedMu sync.Mutex
	shouldQuit bool
}

func NewApp() *App {
	return &App{notified: make(map[string]time.Time)}
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
	dao.InitDB("yatori.db")

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
						req := entity.AddAccountRequest{
							AccountType:   user.AccountType,
							Url:           user.URL,
							RemarkName:    user.RemarkName,
							Account:       user.Account,
							Password:      user.Password,
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
