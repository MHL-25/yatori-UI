package entity

import "time"

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
	// 账号级自检日期区间（YYYY-MM-DD，含首尾）
	SelfCheckSw        int    `json:"selfCheckSw,omitempty"`
	SelfCheckStartDate string `json:"selfCheckStartDate,omitempty"`
	SelfCheckEndDate   string `json:"selfCheckEndDate,omitempty"`
	CxExtendReadSw     int    `json:"cxExtendReadSw,omitempty"`
	CxBbsSw            int    `json:"cxBbsSw,omitempty"`
	CxWorkRetake       int    `json:"cxWorkRetake,omitempty"`
	CxRetakeRate       int    `json:"cxRetakeRate,omitempty"`
}

// CourseLockPO 记录有定时开放要求的课程节点，供启动/定时扫描自动续刷
type CourseLockPO struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Uid         string    `gorm:"index;not null" json:"uid"`
	AccountName string    `json:"accountName"`
	CourseID    string    `gorm:"column:course_id" json:"courseId"`
	CourseName  string    `json:"courseName"`
	NodeID      string    `gorm:"column:node_id" json:"nodeId"`
	NodeName    string    `json:"nodeName"`
	Label       string    `json:"label"`
	UnlockAt    time.Time `json:"unlockAt"`
	Reason      string    `json:"reason"` // time=定时开放 prereq=前置未完成
	Handled     bool      `json:"handled"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// AccountCheckPO 保存账号最近一次课程质检结果
type AccountCheckPO struct {
	Uid       string    `gorm:"primaryKey" json:"uid"`
	CheckedAt time.Time `json:"checkedAt"`
	Summary   string    `json:"summary"`
}

// CheckChapter 单章节质检结果
type CheckChapter struct {
	Label       string `json:"label"`
	Name        string `json:"name"`
	NodeID      string `json:"nodeId"`
	JobTotal    int    `json:"jobTotal"`
	JobFinished int    `json:"jobFinished"`
	Status      string `json:"status"`
	LockType    string `json:"lockType"` // open/time/prereq/expired/closed
	UnlockAt    string `json:"unlockAt,omitempty"`
	EndAt       string `json:"endAt,omitempty"`
	Message     string `json:"message,omitempty"`
}

// CheckCourse 单课程质检结果
type CheckCourse struct {
	CourseID    string         `json:"courseId"`
	CourseName  string         `json:"courseName"`
	JobTotal    int            `json:"jobTotal"`
	JobFinished int            `json:"jobFinished"`
	JobRate     float64        `json:"jobRate"`
	CanStudy    bool           `json:"canStudy"`
	Reason      string         `json:"reason,omitempty"`
	Chapters    []CheckChapter `json:"chapters"`
}

// CheckReport 账号课程质检报告
type CheckReport struct {
	Uid              string        `json:"uid"`
	AccountName      string        `json:"accountName"`
	CheckedAt        time.Time     `json:"checkedAt"`
	TotalCourses     int           `json:"totalCourses"`
	CompletedCourses int           `json:"completedCourses"`
	PendingCourses   int           `json:"pendingCourses"`
	LockedCourses    int           `json:"lockedCourses"`
	Courses          []CheckCourse `json:"courses"`
}

// CourseBrief 课程简要信息，用于“只刷/排除课程”一键选择
type CourseBrief struct {
	CourseID   string  `json:"courseId"`
	CourseName string  `json:"courseName"`
	JobRate    float64 `json:"jobRate"`
	IsStart    bool    `json:"isStart"`
}

type AddAccountRequest struct {
	AccountType    string            `json:"accountType"`
	Url            string            `json:"url"`
	RemarkName     string            `json:"remarkName"`
	Account        string            `json:"account"`
	Password       string            `json:"password"`
	IsProxy        int               `json:"isProxy"`
	InformEmails   []string          `json:"informEmails"`
	ShowDocSw      int               `json:"showDocSw"`
	ShowDocUrls    []string          `json:"showDocUrls"`
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
	ShowDocSw      int               `json:"showDocSw"`
	ShowDocUrls    []string          `json:"showDocUrls"`
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
