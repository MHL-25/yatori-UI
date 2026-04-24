package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

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
