package entity

type UserPO struct {
	Uid            string `gorm:"not null;primaryKey" json:"uid"`
	AccountType    string `gorm:"not null;column:account_type" json:"accountType"`
	Url            string `gorm:"column:url" json:"url"`
	RemarkName     string `gorm:"column:remark_name" json:"remarkName"`
	Account        string `gorm:"not null;column:account" json:"account"`
	Password       string `gorm:"not null;column:password" json:"password"`
	UserConfigJson string `gorm:"not null;column:user_config_json" json:"userConfigJson"`
	IsRunning      bool   `gorm:"-" json:"isRunning"`
}

func (u *UserPO) DisplayName() string {
	if u.RemarkName != "" {
		return u.RemarkName
	}
	return u.Account
}

type SettingPO struct {
	Key   string `gorm:"not null;primaryKey" json:"key"`
	Value string `gorm:"not null" json:"value"`
}

type CoursesCustomData struct {
	VideoModel      int      `json:"videoModel,omitempty"`
	AutoExam        int      `json:"autoExam,omitempty"`
	ExamAutoSubmit  int      `json:"examAutoSubmit,omitempty"`
	CxNode          int      `json:"cxNode,omitempty"`
	CxChapterTestSw int      `json:"cxChapterTestSw,omitempty"`
	CxWorkSw        int      `json:"cxWorkSw,omitempty"`
	CxExamSw        int      `json:"cxExamSw,omitempty"`
	ShuffleSw       int      `json:"shuffleSw,omitempty"`
	StudyTime       string   `json:"studyTime,omitempty"`
	IncludeCourses  []string `json:"includeCourses,omitempty"`
	ExcludeCourses  []string `json:"excludeCourses,omitempty"`
}

type AddAccountRequest struct {
	AccountType    string            `json:"accountType"`
	Url            string            `json:"url"`
	RemarkName     string            `json:"remarkName"`
	Account        string            `json:"account"`
	Password       string            `json:"password"`
	IsProxy        int               `json:"isProxy"`
	InformEmails   []string          `json:"informEmails"`
	CoursesCustom  *CoursesCustomData `json:"coursesCustom,omitempty"`
}

type UpdateAccountRequest struct {
	Uid            string            `json:"uid"`
	AccountType    string            `json:"accountType"`
	Url            string            `json:"url"`
	RemarkName     string            `json:"remarkName"`
	Account        string            `json:"account"`
	Password       string            `json:"password"`
	IsProxy        int               `json:"isProxy"`
	InformEmails   []string          `json:"informEmails"`
	CoursesCustom  *CoursesCustomData `json:"coursesCustom,omitempty"`
}

type CourseInfo struct {
	CourseId   string  `json:"courseId"`
	CourseName string  `json:"courseName"`
	Progress   float32 `json:"progress"`
	Instructor string  `json:"instructor"`
	Status     string  `json:"status"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(msg string, data interface{}) Response {
	return Response{Code: 200, Message: msg, Data: data}
}

func ErrorResponse(msg string) Response {
	return Response{Code: 400, Message: msg}
}
