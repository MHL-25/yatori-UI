package monitor

import (
	"sync"
)

type TaskStatus string

const (
	StatusIdle       TaskStatus = "idle"
	StatusLogging    TaskStatus = "logging"
	StatusRunning    TaskStatus = "running"
	StatusPaused     TaskStatus = "paused"
	StatusCompleted  TaskStatus = "completed"
	StatusError      TaskStatus = "error"
	StatusStopped    TaskStatus = "stopped"
)

type TaskProgress struct {
	Uid          string     `json:"uid"`
	AccountName  string     `json:"accountName"`
	Platform     string     `json:"platform"`
	PlatformName string     `json:"platformName"`
	Status       TaskStatus `json:"status"`
	Progress     float64    `json:"progress"`
	CurrentTask  string     `json:"currentTask"`
	TotalCourses int        `json:"totalCourses"`
	DoneCourses  int        `json:"doneCourses"`
	Logs         []string   `json:"logs"`
	ErrorMessage string     `json:"errorMessage,omitempty"`
}

type EventBus struct {
	mu        sync.RWMutex
	progress  map[string]*TaskProgress
	listeners []func(eventType string, data interface{})
}

var GlobalEventBus = NewEventBus()

func NewEventBus() *EventBus {
	return &EventBus{
		progress: make(map[string]*TaskProgress),
	}
}

func (eb *EventBus) AddListener(fn func(eventType string, data interface{})) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.listeners = append(eb.listeners, fn)
}

func (eb *EventBus) emit(eventType string, data interface{}) {
	eb.mu.RLock()
	listeners := make([]func(eventType string, data interface{}), len(eb.listeners))
	copy(listeners, eb.listeners)
	eb.mu.RUnlock()
	for _, fn := range listeners {
		go fn(eventType, data)
	}
}

func (eb *EventBus) InitProgress(uid, accountName, platform string) {
	platformName := getPlatformName(platform)
	p := &TaskProgress{
		Uid:          uid,
		AccountName:  accountName,
		Platform:     platform,
		PlatformName: platformName,
		Status:       StatusIdle,
		Progress:     0,
		CurrentTask:  "等待启动",
		Logs:         []string{},
	}
	eb.mu.Lock()
	eb.progress[uid] = p
	eb.mu.Unlock()
	eb.emit("progress_update", p)
}

func (eb *EventBus) UpdateStatus(uid string, status TaskStatus) {
	var p *TaskProgress
	eb.mu.Lock()
	if pp, ok := eb.progress[uid]; ok {
		pp.Status = status
		p = pp
	}
	eb.mu.Unlock()
	if p != nil {
		eb.emit("progress_update", p)
	}
}

func (eb *EventBus) UpdateProgress(uid string, progress float64, currentTask string) {
	var p *TaskProgress
	eb.mu.Lock()
	if pp, ok := eb.progress[uid]; ok {
		pp.Progress = progress
		pp.CurrentTask = currentTask
		p = pp
	}
	eb.mu.Unlock()
	if p != nil {
		eb.emit("progress_update", p)
	}
}

func (eb *EventBus) UpdateCourseProgress(uid string, done, total int) {
	var p *TaskProgress
	eb.mu.Lock()
	if pp, ok := eb.progress[uid]; ok {
		pp.DoneCourses = done
		pp.TotalCourses = total
		if total > 0 {
			pp.Progress = float64(done) / float64(total) * 100
		}
		p = pp
	}
	eb.mu.Unlock()
	if p != nil {
		eb.emit("progress_update", p)
	}
}

func (eb *EventBus) AddLog(uid string, log string) {
	var p *TaskProgress
	eb.mu.Lock()
	if pp, ok := eb.progress[uid]; ok {
		pp.Logs = append(pp.Logs, log)
		if len(pp.Logs) > 200 {
			pp.Logs = pp.Logs[len(pp.Logs)-200:]
		}
		p = pp
	}
	eb.mu.Unlock()
	if p != nil {
		eb.emit("log_update", map[string]interface{}{"uid": uid, "log": log})
	}
}

func (eb *EventBus) SetError(uid string, errMsg string) {
	var p *TaskProgress
	eb.mu.Lock()
	if pp, ok := eb.progress[uid]; ok {
		pp.Status = StatusError
		pp.ErrorMessage = errMsg
		p = pp
	}
	eb.mu.Unlock()
	if p != nil {
		eb.emit("progress_update", p)
	}
}

func (eb *EventBus) GetProgress(uid string) *TaskProgress {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	if p, ok := eb.progress[uid]; ok {
		cp := *p
		cp.Logs = make([]string, len(p.Logs))
		copy(cp.Logs, p.Logs)
		return &cp
	}
	return nil
}

func (eb *EventBus) GetAllProgress() []*TaskProgress {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	result := make([]*TaskProgress, 0, len(eb.progress))
	for _, p := range eb.progress {
		cp := *p
		cp.Logs = make([]string, len(p.Logs))
		copy(cp.Logs, p.Logs)
		result = append(result, &cp)
	}
	return result
}

func (eb *EventBus) RemoveProgress(uid string) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	delete(eb.progress, uid)
}

func getPlatformName(platform string) string {
	names := map[string]string{
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
	if name, ok := names[platform]; ok {
		return name
	}
	return platform
}
