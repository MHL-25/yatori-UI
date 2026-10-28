package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"yatori-UI/internal/activity"
	"yatori-UI/internal/checker"
	"yatori-UI/internal/config"
	"yatori-UI/internal/dao"
	"yatori-UI/internal/entity"
	"yatori-UI/internal/monitor"

	"github.com/google/uuid"
)

func intPtr(v int) *int { return &v }

func GetAccountList() ([]entity.UserPO, error) {
	users, err := dao.QueryAllUsers()
	if err != nil {
		return nil, err
	}
	for i := range users {
		uid := users[i].Uid
		act := activity.GetActivity(uid)
		if act != nil {
			users[i].IsRunning = act.IsRunning()
		}
	}
	return users, nil
}

func toYingHuaPlatforms(items []entity.YingHuaUrlItem) []config.YingHuaPlatform {
	if len(items) == 0 {
		return nil
	}
	result := make([]config.YingHuaPlatform, 0, len(items))
	for _, item := range items {
		url := strings.TrimSpace(item.URL)
		if url == "" {
			continue
		}
		result = append(result, config.YingHuaPlatform{URL: url, Password: strings.TrimSpace(item.Password)})
	}
	return result
}

func AddAccount(req entity.AddAccountRequest) (*entity.UserPO, error) {
	existing, err := dao.QueryUserByAccount(req.AccountType, req.Account)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrAccountExists
	}

	uuidV7, _ := uuid.NewV7()
	platforms := toYingHuaPlatforms(req.YingHuaUrls)
	rootUrl := req.Url
	if rootUrl == "" && len(platforms) > 0 {
		rootUrl = platforms[0].URL
	}
	userPO := entity.UserPO{
		Uid:         uuidV7.String(),
		AccountType: req.AccountType,
		Url:         rootUrl,
		RemarkName:  req.RemarkName,
		Account:     req.Account,
		Password:    req.Password,
	}

	custom := config.CoursesCustom{
		VideoModel:      1,
		AutoExam:        1,
		ExamAutoSubmit:  1,
		CxNode:          intPtr(3),
		CxChapterTestSw: intPtr(1),
		CxWorkSw:        intPtr(1),
		CxExamSw:        intPtr(1),
	}
	if req.CoursesCustom != nil {
		custom.VideoModel = req.CoursesCustom.VideoModel
		custom.AutoExam = req.CoursesCustom.AutoExam
		custom.ExamAutoSubmit = req.CoursesCustom.ExamAutoSubmit
		custom.CxNode = intPtr(req.CoursesCustom.CxNode)
		custom.CxChapterTestSw = intPtr(req.CoursesCustom.CxChapterTestSw)
		custom.CxWorkSw = intPtr(req.CoursesCustom.CxWorkSw)
		custom.CxExamSw = intPtr(req.CoursesCustom.CxExamSw)
		custom.ShuffleSw = req.CoursesCustom.ShuffleSw
		custom.StudyTime = req.CoursesCustom.StudyTime
		custom.IncludeCourses = req.CoursesCustom.IncludeCourses
		custom.ExcludeCourses = req.CoursesCustom.ExcludeCourses
		custom.SelfCheckSw = req.CoursesCustom.SelfCheckSw
		custom.SelfCheckStartDate = req.CoursesCustom.SelfCheckStartDate
		custom.SelfCheckEndDate = req.CoursesCustom.SelfCheckEndDate
		custom.CxExtendReadSw = req.CoursesCustom.CxExtendReadSw
		custom.CxBbsSw = req.CoursesCustom.CxBbsSw
		custom.CxWorkRetake = req.CoursesCustom.CxWorkRetake
		custom.CxRetakeRate = req.CoursesCustom.CxRetakeRate
		if custom.CxRetakeRate == 0 {
			custom.CxRetakeRate = 80
		}
		if *custom.CxNode == 0 { custom.CxNode = intPtr(3) }
	}

	userConfig := config.User{
		AccountType:   req.AccountType,
		URL:           rootUrl,
		RemarkName:    req.RemarkName,
		Account:       req.Account,
		Password:      req.Password,
		YingHuaUrls:   platforms,
		IsProxy:       req.IsProxy,
		InformEmails:  req.InformEmails,
		ShowDocSw:     req.ShowDocSw,
		ShowDocUrls:   req.ShowDocUrls,
		CoursesCustom: custom,
	}
	userConfigJson, err := json.Marshal(userConfig)
	if err != nil {
		return nil, err
	}
	userPO.UserConfigJson = string(userConfigJson)

	if err := dao.InsertUser(&userPO); err != nil {
		return nil, err
	}

	monitor.GlobalEventBus.InitProgress(userPO.Uid, userPO.DisplayName(), userPO.AccountType)
	return &userPO, nil
}

func DeleteAccount(uid string) error {
	act := activity.GetActivity(uid)
	if act != nil && act.IsRunning() {
		act.Stop()
	}
	activity.PutActivity(uid, nil)
	monitor.GlobalEventBus.RemoveProgress(uid)
	return dao.DeleteUser(uid)
}

func UpdateAccount(req entity.UpdateAccountRequest) error {
	updateData := make(map[string]interface{})
	if req.AccountType != "" {
		updateData["account_type"] = req.AccountType
	}
	if req.Url != "" {
		updateData["url"] = req.Url
	}
	if req.RemarkName != "" {
		updateData["remark_name"] = req.RemarkName
	}
	if req.Account != "" {
		updateData["account"] = req.Account
	}
	if req.Password != "" {
		updateData["password"] = req.Password
	}

	if len(updateData) == 0 && req.CoursesCustom == nil {
		return ErrNoFieldsToUpdate
	}

	user, err := dao.QueryUserByUid(req.Uid)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrAccountNotFound
	}

	var existingConfig config.User
	if err := json.Unmarshal([]byte(user.UserConfigJson), &existingConfig); err != nil {
		return fmt.Errorf("账号配置解析失败: %v", err)
	}
	if existingConfig.CoursesCustom.CxNode == nil {
		existingConfig.CoursesCustom.CxNode = intPtr(3)
	}

	if req.AccountType != "" { existingConfig.AccountType = req.AccountType }
	if req.Url != "" { existingConfig.URL = req.Url }
	if req.YingHuaUrls != nil {
		existingConfig.YingHuaUrls = toYingHuaPlatforms(req.YingHuaUrls)
		if existingConfig.URL == "" && len(existingConfig.YingHuaUrls) > 0 {
			existingConfig.URL = existingConfig.YingHuaUrls[0].URL
		}
	}
	if req.RemarkName != "" { existingConfig.RemarkName = req.RemarkName }
	if req.Account != "" { existingConfig.Account = req.Account }
	if req.Password != "" { existingConfig.Password = req.Password }
	existingConfig.IsProxy = req.IsProxy
	if req.InformEmails != nil { existingConfig.InformEmails = req.InformEmails }
	existingConfig.ShowDocSw = req.ShowDocSw
	if req.ShowDocUrls != nil { existingConfig.ShowDocUrls = req.ShowDocUrls }

	if req.CoursesCustom != nil {
		existingConfig.CoursesCustom.VideoModel = req.CoursesCustom.VideoModel
		existingConfig.CoursesCustom.AutoExam = req.CoursesCustom.AutoExam
		existingConfig.CoursesCustom.ExamAutoSubmit = req.CoursesCustom.ExamAutoSubmit
		existingConfig.CoursesCustom.CxNode = intPtr(req.CoursesCustom.CxNode)
		existingConfig.CoursesCustom.CxChapterTestSw = intPtr(req.CoursesCustom.CxChapterTestSw)
		existingConfig.CoursesCustom.CxWorkSw = intPtr(req.CoursesCustom.CxWorkSw)
		existingConfig.CoursesCustom.CxExamSw = intPtr(req.CoursesCustom.CxExamSw)
		existingConfig.CoursesCustom.ShuffleSw = req.CoursesCustom.ShuffleSw
		existingConfig.CoursesCustom.StudyTime = req.CoursesCustom.StudyTime
		existingConfig.CoursesCustom.IncludeCourses = req.CoursesCustom.IncludeCourses
		existingConfig.CoursesCustom.ExcludeCourses = req.CoursesCustom.ExcludeCourses
		existingConfig.CoursesCustom.SelfCheckSw = req.CoursesCustom.SelfCheckSw
		existingConfig.CoursesCustom.SelfCheckStartDate = req.CoursesCustom.SelfCheckStartDate
		existingConfig.CoursesCustom.SelfCheckEndDate = req.CoursesCustom.SelfCheckEndDate
		existingConfig.CoursesCustom.CxExtendReadSw = req.CoursesCustom.CxExtendReadSw
		existingConfig.CoursesCustom.CxBbsSw = req.CoursesCustom.CxBbsSw
		existingConfig.CoursesCustom.CxWorkRetake = req.CoursesCustom.CxWorkRetake
		existingConfig.CoursesCustom.CxRetakeRate = req.CoursesCustom.CxRetakeRate
		if *existingConfig.CoursesCustom.CxNode == 0 {
			existingConfig.CoursesCustom.CxNode = intPtr(3)
		}
	}

	userConfigJson, _ := json.Marshal(existingConfig)
	updateData["user_config_json"] = string(userConfigJson)

	if err := dao.UpdateUser(req.Uid, updateData); err != nil {
		return err
	}
	// 配置变更后丢弃旧会话，避免沿用旧的单URL/旧密码缓存
	if act := activity.GetActivity(req.Uid); act != nil && !act.IsRunning() {
		activity.PutActivity(req.Uid, nil)
	}
	return nil
}

func getSetting() config.Setting {
	cfg := config.EnsureConfigFile("./config.yaml")
	return cfg.Setting
}

func LoginCheck(uid string) error {
	user, err := dao.QueryUserByUid(uid)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrAccountNotFound
	}

	act := activity.GetActivity(uid)
	if act == nil {
		act = activity.BuildActivity(*user, getSetting())
		if act == nil {
			return ErrUnsupportedPlatform
		}
		activity.PutActivity(uid, act)
	}

	monitor.GlobalEventBus.UpdateStatus(uid, monitor.StatusLogging)
	monitor.GlobalEventBus.AddLog(uid, "正在验证登录...")
	err = act.Login()
	if err != nil {
		monitor.GlobalEventBus.SetError(uid, "登录失败: "+err.Error())
		return err
	}
	monitor.GlobalEventBus.UpdateStatus(uid, monitor.StatusIdle)
	monitor.GlobalEventBus.AddLog(uid, "登录验证成功")
	return nil
}

func StartBrush(uid string) error {
	user, err := dao.QueryUserByUid(uid)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrAccountNotFound
	}

	act := activity.GetActivity(uid)
	if act == nil {
		act = activity.BuildActivity(*user, getSetting())
		if act == nil {
			return ErrUnsupportedPlatform
		}
		activity.PutActivity(uid, act)
	}

	if act.IsRunning() {
		return ErrAlreadyRunning
	}

	return act.Start()
}

func StopBrush(uid string) error {
	act := activity.GetActivity(uid)
	if act == nil {
		return ErrAccountNotFound
	}
	if !act.IsRunning() {
		return ErrNotRunning
	}
	return act.Stop()
}

func PauseBrush(uid string) error {
	act := activity.GetActivity(uid)
	if act == nil {
		return ErrAccountNotFound
	}
	if !act.IsRunning() {
		return ErrNotRunning
	}
	return act.Pause()
}

func GetProgress(uid string) *monitor.TaskProgress {
	return monitor.GlobalEventBus.GetProgress(uid)
}

func GetAllProgress() []*monitor.TaskProgress {
	return monitor.GlobalEventBus.GetAllProgress()
}

func GetConfig() *config.JSONDataForConfig {
	return config.EnsureConfigFile("./config.yaml")
}

func SaveConfig(cfg *config.JSONDataForConfig) error {
	return config.SaveConfig(cfg, "./config.yaml")
}

func ImportConfig(filePath string) (*config.JSONDataForConfig, error) {
	return config.ImportConfigFile(filePath)
}

func ImportAndMergeConfig(filePath string) (int, error) {
	cfg, err := config.ImportConfigFile(filePath)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, user := range cfg.Users {
		existing, _ := dao.QueryUserByAccount(user.AccountType, user.Account)
		if existing != nil {
			continue
		}
		uuidV7, _ := uuid.NewV7()
		userConfigJson, _ := json.Marshal(user)
		userPO := &entity.UserPO{
			Uid:            uuidV7.String(),
			AccountType:    user.AccountType,
			Url:            user.URL,
			RemarkName:     user.RemarkName,
			Account:        user.Account,
			Password:       user.Password,
			UserConfigJson: string(userConfigJson),
		}
		if err := dao.InsertUser(userPO); err != nil {
			continue
		}
		monitor.GlobalEventBus.InitProgress(userPO.Uid, userPO.DisplayName(), userPO.AccountType)
		count++
	}
	if cfg.Setting != (config.Setting{}) {
		config.SaveConfig(cfg, "./config.yaml")
	}
	return count, nil
}

func GetPlatforms() []map[string]string {
	platforms := make([]map[string]string, 0)
	for k, v := range config.PlatformNames {
		platforms = append(platforms, map[string]string{
			"id":   k,
			"name": v,
		})
	}
	return platforms
}

func GetAiTypes() []map[string]string {
	types := make([]map[string]string, 0)
	for k, v := range config.AiTypeNames {
		types = append(types, map[string]string{
			"id":   string(k),
			"name": v,
		})
	}
	return types
}

func StartAllBrush() (int, int) {
	users, err := dao.QueryAllUsers()
	if err != nil {
		return 0, 0
	}
	success := 0
	total := 0
	for _, user := range users {
		act := activity.GetActivity(user.Uid)
		if act != nil && act.IsRunning() {
			continue
		}
		total++
		if act == nil {
			act = activity.BuildActivity(user, getSetting())
			if act == nil {
				continue
			}
			activity.PutActivity(user.Uid, act)
		}
		if err := act.Start(); err == nil {
			success++
		}
	}
	return success, total
}

func StartBatchBrush(uids []string) (int, int) {
	success := 0
	total := len(uids)
	for _, uid := range uids {
		user, err := dao.QueryUserByUid(uid)
		if err != nil || user == nil {
			continue
		}
		act := activity.GetActivity(uid)
		if act != nil && act.IsRunning() {
			continue
		}
		if act == nil {
			act = activity.BuildActivity(*user, getSetting())
			if act == nil {
				continue
			}
			activity.PutActivity(uid, act)
		}
		if err := act.Start(); err == nil {
			success++
		}
	}
	return success, total
}

// CheckAccount 单个账号课程质检
func CheckAccount(uid string) (*entity.CheckReport, error) {
	return checker.Check(uid, getSetting())
}

// CheckAccounts 批量账号课程质检
func CheckAccounts(uids []string) []*entity.CheckReport {
	return checker.CheckMany(uids, getSetting())
}

// CheckAllAccounts 全部账号课程质检
func CheckAllAccounts() []*entity.CheckReport {
	return checker.CheckAll(getSetting())
}

// GetStoredCheckReports 读取全部账号最近一次质检结果
func GetStoredCheckReports() []*entity.CheckReport {
	return checker.GetAllStored()
}

// GetCourseLocks 读取账号待处理的定时开放记录
func GetCourseLocks(uid string) ([]entity.CourseLockPO, error) {
	return dao.QueryPendingCourseLocks(uid)
}

// GetCourseList 拉取账号课程名称列表
func GetCourseList(uid string) ([]entity.CourseBrief, error) {
	return checker.CourseList(uid, getSetting())
}

// GetUserByUid 查询账号
func GetUserByUid(uid string) (*entity.UserPO, error) {
	return dao.QueryUserByUid(uid)
}
