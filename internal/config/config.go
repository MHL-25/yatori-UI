package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/yatori-dev/yatori-go-core/models/ctype"
	"go.yaml.in/yaml/v3"
)

type JSONDataForConfig struct {
	Setting Setting `json:"setting"`
	Users   []User  `json:"users"`
}

type EmailInform struct {
	Sw       int    `json:"sw"`
	SMTPHost string `json:"smtpHost" yaml:"SMTPHost"`
	SMTPPort int    `json:"smtpPort" yaml:"SMTPPort"`
	UserName string `json:"userName" yaml:"userName"`
	Password string `json:"password"`
}

type BasicSetting struct {
	CompletionTone int    `default:"1" json:"completionTone,omitempty" yaml:"completionTone"`
	ColorLog       int    `json:"colorLog,omitempty" yaml:"colorLog"`
	LogOutFileSw   int    `json:"logOutFileSw,omitempty" yaml:"logOutFileSw"`
	LogLevel       string `json:"logLevel,omitempty" yaml:"logLevel"`
	LogModel       int    `json:"logModel" yaml:"logModel"`
	WebModel       int    `json:"webModel" yaml:"webModel"`
	// 启动自动化：AutoCheckSw=1 启动后自动质检；AutoStartSw=1 质检后自动开始刷课
	AutoCheckSw int `json:"autoCheckSw" yaml:"autoCheckSw"`
	AutoStartSw int `json:"autoStartSw" yaml:"autoStartSw"`
}

type AiSetting struct {
	AiType ctype.AiType `json:"aiType" yaml:"aiType"`
	AiUrl  string       `json:"aiUrl" yaml:"aiUrl"`
	Model  string       `json:"model"`
	APIKEY string       `json:"API_KEY" yaml:"API_KEY" mapstructure:"API_KEY"`
}

type VisionAiSetting struct {
	AiType ctype.AiType `json:"aiType" yaml:"aiType"`
	AiUrl  string       `json:"aiUrl" yaml:"aiUrl"`
	Model  string       `json:"model"`
	APIKEY string       `json:"API_KEY" yaml:"API_KEY" mapstructure:"API_KEY"`
}

type ApiQueSetting struct {
	Url string `json:"url"`
}

type Setting struct {
	BasicSetting    BasicSetting    `json:"basicSetting" yaml:"basicSetting"`
	EmailInform     EmailInform     `json:"emailInform" yaml:"emailInform"`
	AiSetting       AiSetting       `json:"aiSetting" yaml:"aiSetting"`
	VisionAiSetting VisionAiSetting `json:"visionAiSetting" yaml:"visionAiSetting"`
	ApiQueSetting   ApiQueSetting   `json:"apiQueSetting" yaml:"apiQueSetting"`
}

type CoursesSettings struct {
	Name         string   `json:"name"`
	IncludeExams []string `json:"includeExams" yaml:"includeExams"`
	ExcludeExams []string `json:"excludeExams" yaml:"excludeExams"`
}

type CoursesCustom struct {
	StudyTime       string            `json:"studyTime" yaml:"studyTime"`
	CxNode          *int              `json:"cxNode" yaml:"cxNode"`
	CxChapterTestSw *int              `json:"cxChapterTestSw" yaml:"cxChapterTestSw"`
	CxWorkSw        *int              `json:"cxWorkSw" yaml:"cxWorkSw"`
	CxExamSw        *int              `json:"cxExamSw" yaml:"cxExamSw"`
	ShuffleSw       int               `json:"shuffleSw" yaml:"shuffleSw"`
	VideoModel      int               `json:"videoModel" yaml:"videoModel"`
	AutoExam        int               `json:"autoExam" yaml:"autoExam"`
	ExamAutoSubmit  int               `json:"examAutoSubmit" yaml:"examAutoSubmit"`
	ExcludeCourses  []string          `json:"excludeCourses" yaml:"excludeCourses"`
	IncludeCourses  []string          `json:"includeCourses" yaml:"includeCourses"`
	CoursesSettings []CoursesSettings `json:"coursesSettings" yaml:"coursesSettings"`
	// 账号级自检日期区间（YYYY-MM-DD，含首尾两天）
	SelfCheckSw        int    `json:"selfCheckSw" yaml:"selfCheckSw"`
	SelfCheckStartDate string `json:"selfCheckStartDate" yaml:"selfCheckStartDate"`
	SelfCheckEndDate   string `json:"selfCheckEndDate" yaml:"selfCheckEndDate"`
	// 学习通扩展任务：拓展阅读、讨论区全部话题、作业重做
	CxExtendReadSw int `json:"cxExtendReadSw" yaml:"cxExtendReadSw"`
	CxBbsSw        int `json:"cxBbsSw" yaml:"cxBbsSw"`
	CxWorkRetake   int `json:"cxWorkRetake" yaml:"cxWorkRetake"`
	CxRetakeRate   int `json:"cxRetakeRate" yaml:"cxRetakeRate"`
}

type User struct {
	AccountType   string        `json:"accountType" yaml:"accountType"`
	URL           string        `json:"url"`
	RemarkName    string        `json:"remarkName,omitempty" yaml:"remarkName,omitempty" mapstructure:"remarkName"`
	Account       string        `json:"account"`
	Password      string        `json:"password"`
	IsProxy       int           `json:"isProxy" yaml:"isProxy"`
	InformEmails  []string      `json:"informEmails" yaml:"informEmails"`
	CoursesCustom CoursesCustom `json:"coursesCustom" yaml:"coursesCustom"`
	// ShowDoc 推送（账号级独立开关，可多个地址）
	ShowDocSw   int      `json:"showDocSw" yaml:"showDocSw"`
	ShowDocUrls []string `json:"showDocUrls" yaml:"showDocUrls"`
}

func defaultValue(config *JSONDataForConfig) {
	for i := range config.Users {
		v := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		if config.Users[i].CoursesCustom.CxNode == nil {
			(&config.Users[i].CoursesCustom).CxNode = &v[3]
		}
		if config.Users[i].CoursesCustom.CxChapterTestSw == nil {
			(&config.Users[i].CoursesCustom).CxChapterTestSw = &v[1]
		}
		if config.Users[i].CoursesCustom.CxWorkSw == nil {
			(&config.Users[i].CoursesCustom).CxWorkSw = &v[1]
		}
		if config.Users[i].CoursesCustom.CxExamSw == nil {
			(&config.Users[i].CoursesCustom).CxExamSw = &v[1]
		}
	}
}

func ReadConfig(filePath string) JSONDataForConfig {
	var configJson JSONDataForConfig
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		return configJson
	}
	err = viper.Unmarshal(&configJson)
	defaultValue(&configJson)
	if err != nil {
		return configJson
	}
	// 旧配置文件没有启动自动化字段时，默认开启
	if !viper.IsSet("setting.basicSetting.autoCheckSw") {
		configJson.Setting.BasicSetting.AutoCheckSw = 1
	}
	if !viper.IsSet("setting.basicSetting.autoStartSw") {
		configJson.Setting.BasicSetting.AutoStartSw = 1
	}
	return configJson
}

func SaveConfig(config *JSONDataForConfig, filePath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

func EnsureConfigFile(filePath string) *JSONDataForConfig {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		cfg := &JSONDataForConfig{}
		cfg.Setting.BasicSetting.CompletionTone = 1
		cfg.Setting.BasicSetting.ColorLog = 1
		cfg.Setting.BasicSetting.LogOutFileSw = 1
		cfg.Setting.BasicSetting.LogLevel = "INFO"
		cfg.Setting.BasicSetting.LogModel = 0
		cfg.Setting.BasicSetting.AutoCheckSw = 1
		cfg.Setting.BasicSetting.AutoStartSw = 1
		cfg.Setting.AiSetting.AiType = "TONGYI"
		cfg.Setting.ApiQueSetting.Url = "http://localhost:8083"
		SaveConfig(cfg, filePath)
		return cfg
	}
	cfg := ReadConfig(filePath)
	return &cfg
}

var PlatformNames = map[string]string{
	"XUEXITONG": "学习通",
	"YINGHUA":   "英华学堂",
	"CANGHUI":   "仓辉实训",
	"ENAEA":     "学习公社",
	"CQIE":      "重庆工程学院",
	"KETANGX":   "码上研训",
	"ICVE":      "智慧职教",
	"QSXT":      "青书学堂",
	"WELEARN":   "WeLearn",
	"HQKJ":      "海旗科技",
	"GONGXUE":   "工学云",
	"WEIBAN":    "安全微伴",
}

var AiTypeNames = map[ctype.AiType]string{
	"TONGYI":    "通义千问",
	"CHATGLM":   "智谱ChatGLM",
	"XINGHUO":   "讯飞星火",
	"DOUBAO":    "豆包",
	"OPENAI":    "OpenAI",
	"DEEPSEEK":  "DeepSeek",
	"SILICON":   "硅基流动",
	"METAAI":    "秘塔AI",
	"OTHER":     "其他兼容",
}

func CmpCourse(name string, list []string) bool {
	for _, item := range list {
		if item == name {
			return true
		}
	}
	return false
}

// InSelfCheckWindow 判断当前系统日期是否在账号配置的自检日期区间内（含首尾）。
// 未开启或日期为空时视为不限日期，返回 true。
func InSelfCheckWindow(user User, now time.Time) (bool, string) {
	custom := user.CoursesCustom
	if custom.SelfCheckSw != 1 {
		return true, "未启用自检日期区间"
	}
	startStr := strings.TrimSpace(custom.SelfCheckStartDate)
	endStr := strings.TrimSpace(custom.SelfCheckEndDate)
	if startStr == "" && endStr == "" {
		return true, "未设置日期区间，按不限日期处理"
	}
	var start, end time.Time
	var err error
	if startStr != "" {
		start, err = time.ParseInLocation("2006-01-02", startStr, time.Local)
		if err != nil {
			return true, "自检开始日期格式错误，按不限日期处理"
		}
	}
	if endStr != "" {
		end, err = time.ParseInLocation("2006-01-02", endStr, time.Local)
		if err != nil {
			return true, "自检结束日期格式错误，按不限日期处理"
		}
	}
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	if !start.IsZero() && today.Before(start) {
		return false, "当前日期早于自检开始日期 " + startStr
	}
	if !end.IsZero() && today.After(end) {
		return false, "当前日期晚于自检结束日期 " + endStr
	}
	return true, fmt.Sprintf("当前日期在自检区间 %s ~ %s 内", orAll(startStr), orAll(endStr))
}

func orAll(v string) string {
	if v == "" {
		return "不限"
	}
	return v
}

func ImportConfigFile(filePath string) (*JSONDataForConfig, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".yaml", ".yml":
		return importYAML(filePath)
	case ".json":
		return importJSON(filePath)
	default:
		return importYAML(filePath)
	}
}

func importYAML(filePath string) (*JSONDataForConfig, error) {
	v := viper.New()
	v.SetConfigFile(filePath)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg JSONDataForConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	defaultValue(&cfg)
	return &cfg, nil
}

func importJSON(filePath string) (*JSONDataForConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var cfg JSONDataForConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	defaultValue(&cfg)
	return &cfg, nil
}
