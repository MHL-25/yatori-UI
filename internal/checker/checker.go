package checker

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"yatori-UI/internal/activity"
	"yatori-UI/internal/config"
	"yatori-UI/internal/dao"
	"yatori-UI/internal/entity"
	"yatori-UI/internal/monitor"
	"yatori-UI/internal/studyplan"

	"github.com/yatori-dev/yatori-go-core/aggregation/xuexitong"
	"github.com/yatori-dev/yatori-go-core/aggregation/yinghua"
	xuexitongApi "github.com/yatori-dev/yatori-go-core/api/xuexitong"
	yinghuaApi "github.com/yatori-dev/yatori-go-core/api/yinghua"
)

var checkLocks sync.Map // map[string]*sync.Mutex

func lockOf(uid string) func() {
	m, _ := checkLocks.LoadOrStore(uid, &sync.Mutex{})
	mu := m.(*sync.Mutex)
	mu.Lock()
	return mu.Unlock
}

// Check 对单个账号执行课程质检
func Check(uid string, setting config.Setting) (*entity.CheckReport, error) {
	unlock := lockOf(uid)
	defer unlock()

	user, err := dao.QueryUserByUid(uid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("账号不存在")
	}
	report := &entity.CheckReport{
		Uid:         uid,
		AccountName: user.DisplayName(),
		CheckedAt:   time.Now(),
		Courses:     make([]entity.CheckCourse, 0),
	}

	monitor.GlobalEventBus.UpdateStatus(uid, monitor.StatusLogging)
	monitor.GlobalEventBus.AddLog(uid, "🔍 开始课程质检...")

	if user.AccountType == "YINGHUA" {
		checkYingHua(user, report, setting)
		save(report)
		monitor.GlobalEventBus.UpdateStatus(uid, monitor.StatusIdle)
		monitor.GlobalEventBus.AddLog(uid, fmt.Sprintf("✅ 质检完成: %d 门课程，完成 %d，待刷 %d，受限 %d",
			report.TotalCourses, report.CompletedCourses, report.PendingCourses, report.LockedCourses))
		return report, nil
	}

	if user.AccountType != "XUEXITONG" {
		report.Courses = append(report.Courses, entity.CheckCourse{
			CourseName: "当前平台暂不支持深度质检",
			Reason:     "仅学习通支持章节级质检",
		})
		save(report)
		monitor.GlobalEventBus.UpdateStatus(uid, monitor.StatusIdle)
		return report, nil
	}

	cache, err := ensureCache(user, setting)
	if err != nil {
		return nil, err
	}

	courses, err := xuexitong.XueXiTPullCourseAction(cache)
	if err != nil {
		monitor.GlobalEventBus.SetError(uid, "拉取课程失败: "+err.Error())
		return nil, err
	}
	monitor.GlobalEventBus.AddLog(uid, fmt.Sprintf("📚 质检课程数: %d", len(courses)))

	for i := range courses {
		course := courses[i]
		cc := entity.CheckCourse{
			CourseID:    course.CourseID,
			CourseName:  course.CourseName,
			JobTotal:    course.JobCount,
			JobFinished: course.JobFinishCount,
			JobRate:     course.JobRate,
			Chapters:    make([]entity.CheckChapter, 0),
		}
		if !course.IsStart {
			cc.Reason = "课程未开始"
		}

		key, _ := strconv.Atoi(course.Key)
		action, _, err := xuexitong.PullCourseChapterAction(cache, course.Cpi, key)
		if err != nil {
			cc.Reason = "拉取章节失败: " + err.Error()
			report.Courses = append(report.Courses, cc)
			monitor.GlobalEventBus.AddLog(uid, fmt.Sprintf("⚠️ [%s] %s", course.CourseName, cc.Reason))
			continue
		}

		var nodes []int
		for _, item := range action.Knowledge {
			if item.Layer == 2 && item.JobCount > 0 {
				nodes = append(nodes, item.ID)
			}
		}
		courseId, _ := strconv.Atoi(course.CourseID)
		userId, _ := strconv.Atoi(cache.UserID)
		pointAction := action
		if len(nodes) > 0 {
			if pa, err := xuexitong.ChapterFetchPointAction(cache, nodes, &action, key, userId, course.Cpi, courseId); err == nil {
				pointAction = pa
			} else {
				monitor.GlobalEventBus.AddLog(uid, fmt.Sprintf("⚠️ [%s] 任务点状态获取失败: %v", course.CourseName, err))
			}
		}

		plans := studyplan.Analyze(pointAction.Knowledge, time.Now())
		canStudy := false
		var timeLocks []entity.CourseLockPO
		lockCount := map[studyplan.LockType]int{}
		completedChapters := 0
		taskTotal, taskFinished := 0, 0
		for _, p := range plans {
			ch := entity.CheckChapter{
				Label:       p.Node.Label,
				Name:        p.Node.Name,
				NodeID:      strconv.Itoa(p.Node.ID),
				JobTotal:    p.Node.JobCount,
				JobFinished: p.Node.PointFinished,
				Status:      p.Node.Status,
				LockType:    string(p.Lock),
				Message:     p.Reason,
			}
			if p.Completed {
				ch.JobFinished = p.Node.JobCount
				completedChapters++
			}
			if ch.JobFinished > ch.JobTotal {
				ch.JobFinished = ch.JobTotal
			}
			taskTotal += ch.JobTotal
			taskFinished += ch.JobFinished
			if p.Node.EndTime > 0 {
				ch.EndAt = time.UnixMilli(p.Node.EndTime).Format("2006-01-02 15:04")
			}
			if p.Lock == studyplan.LockTime {
				ch.UnlockAt = p.UnlockAt.Format("2006-01-02 15:04")
				timeLocks = append(timeLocks, entity.CourseLockPO{
					Uid:         uid,
					AccountName: report.AccountName,
					CourseID:    course.CourseID,
					CourseName:  course.CourseName,
					NodeID:      ch.NodeID,
					NodeName:    p.Node.Name,
					Label:       p.Node.Label,
					UnlockAt:    p.UnlockAt,
					Reason:      "time",
				})
			}
			if !p.Completed && p.Lock == studyplan.LockOpen && p.Node.PointFinished < p.Node.JobCount {
				canStudy = true
			}
			if !p.Completed {
				lockCount[p.Lock]++
			}
			cc.Chapters = append(cc.Chapters, ch)
		}
		// 课程列表的任务点统计为 0 时，用章节汇总值兜底
		if course.JobCount > 0 {
			cc.JobTotal = course.JobCount
			cc.JobFinished = course.JobFinishCount
			cc.JobRate = course.JobRate
		} else if taskTotal > 0 {
			cc.JobTotal = taskTotal
			cc.JobFinished = taskFinished
			cc.JobRate = float64(taskFinished) / float64(taskTotal) * 100
		}
		if !course.IsStart {
			cc.CanStudy = false
			cc.Reason = "课程未开始"
		} else if course.State == 1 {
			cc.CanStudy = false
			cc.Reason = "课程已结束"
		} else {
			cc.CanStudy = canStudy
			if !canStudy && cc.Reason == "" {
				cc.Reason = fmt.Sprintf("已完成 %d 节 / 待刷 %d 节 / 定时 %d 节 / 前置 %d 节 / 已截止 %d 节 / 已关闭 %d 节",
					completedChapters, lockCount[studyplan.LockOpen], lockCount[studyplan.LockTime],
					lockCount[studyplan.LockPrereq], lockCount[studyplan.LockExpired], lockCount[studyplan.LockClosed])
			}
		}
		report.Courses = append(report.Courses, cc)

		if err := dao.ReplaceCourseLocks(uid, course.CourseID, timeLocks); err != nil {
			monitor.GlobalEventBus.AddLog(uid, "⚠️ 保存定时开放记录失败: "+err.Error())
		}
		monitor.GlobalEventBus.AddLog(uid, fmt.Sprintf("📗 [%s] 进度 %.1f%% | 可学习: %v", course.CourseName, course.JobRate, canStudy))
	}

	report.TotalCourses = len(report.Courses)
	for _, c := range report.Courses {
		if c.JobTotal > 0 && c.JobFinished >= c.JobTotal {
			report.CompletedCourses++
		} else if c.CanStudy {
			report.PendingCourses++
		} else {
			report.LockedCourses++
		}
	}
	save(report)
	monitor.GlobalEventBus.UpdateStatus(uid, monitor.StatusIdle)
	monitor.GlobalEventBus.AddLog(uid, fmt.Sprintf("✅ 质检完成: %d 门课程，完成 %d，待刷 %d，受限 %d",
		report.TotalCourses, report.CompletedCourses, report.PendingCourses, report.LockedCourses))
	return report, nil
}

// checkYingHua 英华学堂课程质检：课程列表 + 章节节点 + 锁定/进度
func checkYingHua(user *entity.UserPO, report *entity.CheckReport, setting config.Setting) {
	uid := user.Uid
	cache, err := loginYingHua(user.Url, user.Account, user.Password)
	if err != nil {
		monitor.GlobalEventBus.SetError(uid, "质检登录失败: "+err.Error())
		report.Courses = append(report.Courses, entity.CheckCourse{
			CourseName: "登录失败",
			Reason:     err.Error(),
		})
		report.TotalCourses = 1
		return
	}
	monitor.GlobalEventBus.AddLog(uid, "🔑 英华登录成功")

	courses, err := yinghua.CourseListAction(cache)
	if err != nil {
		monitor.GlobalEventBus.SetError(uid, "拉取课程失败: "+err.Error())
		report.Courses = append(report.Courses, entity.CheckCourse{
			CourseName: "拉取课程失败",
			Reason:     err.Error(),
		})
		report.TotalCourses = 1
		return
	}
	monitor.GlobalEventBus.AddLog(uid, fmt.Sprintf("📚 质检课程数: %d", len(courses)))

	now := time.Now()
	for _, course := range courses {
		cc := entity.CheckCourse{
			CourseID:   course.Id,
			CourseName: course.Name,
			JobRate:    normalizeYingHuaProgress(course.Progress),
			Chapters:   make([]entity.CheckChapter, 0),
		}

		nodes, nodeErr := yinghua.VideosListAction(cache, course)
		if nodeErr != nil {
			cc.Reason = "拉取章节失败: " + nodeErr.Error()
			cc.JobTotal = course.VideoCount
			cc.JobFinished = course.VideoLearned
			report.Courses = append(report.Courses, cc)
			monitor.GlobalEventBus.AddLog(uid, fmt.Sprintf("⚠️ [%s] %s", course.Name, cc.Reason))
			continue
		}

		var timeLocks []entity.CourseLockPO
		completedCount, openCount, lockedCount := 0, 0, 0
		for index, node := range nodes {
			completed := node.Progress >= 99.9 || node.State == 2
			jobTotal, jobFinished := 1, 0
			if !node.TabVideo && !node.TabWork && !node.TabExam && !node.TabFile {
				jobTotal = 0
			}
			if completed && jobTotal > 0 {
				jobFinished = 1
				completedCount++
			}
			lockType := string(studyplan.LockOpen)
			status := "open"
			message := ""
			unlockAt := ""
			if !completed {
				if node.NodeLock > 0 && node.UnlockTime.After(now) {
					lockType = string(studyplan.LockTime)
					status = "time"
					unlockAt = node.UnlockTime.Format("2006-01-02 15:04")
					message = "未到开放时间"
					timeLocks = append(timeLocks, entity.CourseLockPO{
						Uid:         uid,
						AccountName: report.AccountName,
						CourseID:    course.Id,
						CourseName:  course.Name,
						NodeID:      node.Id,
						NodeName:    node.Name,
						Label:       fmt.Sprintf("%d", index+1),
						UnlockAt:    node.UnlockTime,
						Reason:      "time",
					})
					lockedCount++
				} else if node.NodeLock > 0 {
					lockType = string(studyplan.LockClosed)
					status = "close"
					message = "节点未解锁"
					lockedCount++
				} else {
					openCount++
				}
			}
			cc.JobTotal += jobTotal
			cc.JobFinished += jobFinished
			cc.Chapters = append(cc.Chapters, entity.CheckChapter{
				Label:       fmt.Sprintf("%d", index+1),
				Name:        node.Name,
				NodeID:      node.Id,
				JobTotal:    jobTotal,
				JobFinished: jobFinished,
				Status:      status,
				LockType:    lockType,
				UnlockAt:    unlockAt,
				Message:     message,
			})
		}
		if cc.JobTotal == 0 {
			cc.JobTotal = course.VideoCount
			cc.JobFinished = course.VideoLearned
		}
		if cc.JobTotal > 0 {
			cc.JobRate = float64(cc.JobFinished) / float64(cc.JobTotal) * 100
		}

		switch {
		case cc.JobTotal > 0 && cc.JobFinished >= cc.JobTotal:
			cc.CanStudy = false
			cc.Reason = "课程已完成"
		case !course.EndDate.IsZero() && course.EndDate.Before(now):
			cc.CanStudy = false
			cc.Reason = "课程已结束"
		case openCount > 0:
			cc.CanStudy = true
		default:
			cc.CanStudy = false
			cc.Reason = fmt.Sprintf("已完成 %d 节 / 待刷 0 节 / 定时 %d 节 / 锁定 %d 节", completedCount, lockedCount, lockedCount)
		}
		report.Courses = append(report.Courses, cc)
		if err := dao.ReplaceCourseLocks(uid, course.Id, timeLocks); err != nil {
			monitor.GlobalEventBus.AddLog(uid, "⚠️ 保存定时开放记录失败: "+err.Error())
		}
		monitor.GlobalEventBus.AddLog(uid, fmt.Sprintf("📗 [%s] 进度 %.1f%% | 可学习: %v", course.Name, cc.JobRate, cc.CanStudy))
	}

	report.TotalCourses = len(report.Courses)
	for _, c := range report.Courses {
		if c.JobTotal > 0 && c.JobFinished >= c.JobTotal {
			report.CompletedCourses++
		} else if c.CanStudy {
			report.PendingCourses++
		} else {
			report.LockedCourses++
		}
	}
}

func normalizeYingHuaProgress(progress float64) float64 {
	if progress <= 1 {
		return progress * 100
	}
	return progress
}

func loginYingHua(url, account, password string) (*yinghuaApi.YingHuaUserCache, error) {
	cache := &yinghuaApi.YingHuaUserCache{
		PreUrl:   url,
		Account:  account,
		Password: password,
	}
	if yinghua.IsYingHuaCookieOrToken(password) {
		return cache, yinghua.YingHuaCookieLoginAction(cache)
	}
	return cache, yinghua.YingHuaLoginAction(cache)
}

func save(report *entity.CheckReport) {
	data, err := json.Marshal(report)
	if err != nil {
		return
	}
	_ = dao.SaveCheckReport(report, string(data))
}

// ensureCache 复用或创建账号会话
func ensureCache(user *entity.UserPO, setting config.Setting) (*xuexitongApi.XueXiTUserCache, error) {
	act := activity.GetActivity(user.Uid)
	if act == nil {
		act = activity.BuildActivity(*user, setting)
		if act == nil {
			return nil, fmt.Errorf("不支持的平台")
		}
		activity.PutActivity(user.Uid, act)
	}
	if act.GetUserCache() == nil {
		if err := act.Login(); err != nil {
			monitor.GlobalEventBus.SetError(user.Uid, "质检登录失败: "+err.Error())
			return nil, err
		}
	}
	cache, ok := act.GetUserCache().(*xuexitongApi.XueXiTUserCache)
	if !ok {
		return nil, fmt.Errorf("账号会话不可用，请先登录")
	}
	return cache, nil
}

// CourseList 拉取账号课程名称列表，供“只刷/排除课程”一键选择
func CourseList(uid string, setting config.Setting) ([]entity.CourseBrief, error) {
	user, err := dao.QueryUserByUid(uid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("账号不存在")
	}
	if user.AccountType == "YINGHUA" {
		cache, err := loginYingHua(user.Url, user.Account, user.Password)
		if err != nil {
			return nil, err
		}
		courses, err := yinghua.CourseListAction(cache)
		if err != nil {
			return nil, err
		}
		result := make([]entity.CourseBrief, 0, len(courses))
		for _, c := range courses {
			result = append(result, entity.CourseBrief{
				CourseID:   c.Id,
				CourseName: c.Name,
				JobRate:    normalizeYingHuaProgress(c.Progress),
				IsStart:    true,
			})
		}
		return result, nil
	}
	if user.AccountType != "XUEXITONG" {
		return nil, fmt.Errorf("当前平台暂不支持一键查询课程")
	}
	cache, err := ensureCache(user, setting)
	if err != nil {
		return nil, err
	}
	courses, err := xuexitong.XueXiTPullCourseAction(cache)
	if err != nil {
		return nil, err
	}
	result := make([]entity.CourseBrief, 0, len(courses))
	for _, c := range courses {
		result = append(result, entity.CourseBrief{
			CourseID:   c.CourseID,
			CourseName: c.CourseName,
			JobRate:    c.JobRate,
			IsStart:    c.IsStart,
		})
	}
	return result, nil
}

// CourseListByCredentials 使用表单中的账号密码临时登录并查询课程，不写入数据库。
func CourseListByCredentials(req entity.CourseQueryRequest) ([]entity.CourseBrief, error) {
	if req.AccountType == "YINGHUA" {
		if req.Account == "" || req.Password == "" {
			return nil, fmt.Errorf("请先填写账号和密码/Cookie")
		}
		cache, err := loginYingHua(req.URL, req.Account, req.Password)
		if err != nil {
			return nil, fmt.Errorf("登录失败: %w", err)
		}
		courses, err := yinghua.CourseListAction(cache)
		if err != nil {
			return nil, err
		}
		result := make([]entity.CourseBrief, 0, len(courses))
		for _, c := range courses {
			result = append(result, entity.CourseBrief{
				CourseID:   c.Id,
				CourseName: c.Name,
				JobRate:    normalizeYingHuaProgress(c.Progress),
				IsStart:    true,
			})
		}
		return result, nil
	}
	if req.AccountType != "XUEXITONG" {
		return nil, fmt.Errorf("当前平台暂不支持一键查询课程")
	}
	if req.Account == "" || req.Password == "" {
		return nil, fmt.Errorf("请先填写账号和密码")
	}
	cache := &xuexitongApi.XueXiTUserCache{Name: req.Account, Password: req.Password}
	var err error
	if len(req.Password) < 50 {
		err = xuexitong.XueXiTLoginAction(cache)
	} else {
		err = xuexitong.XueXiTCookieLoginAction(cache)
	}
	if err != nil {
		return nil, fmt.Errorf("登录失败: %w", err)
	}
	courses, err := xuexitong.XueXiTPullCourseAction(cache)
	if err != nil {
		return nil, err
	}
	result := make([]entity.CourseBrief, 0, len(courses))
	for _, c := range courses {
		result = append(result, entity.CourseBrief{
			CourseID:   c.CourseID,
			CourseName: c.CourseName,
			JobRate:    c.JobRate,
			IsStart:    c.IsStart,
		})
	}
	return result, nil
}

// CheckMany 批量质检指定账号
func CheckMany(uids []string, setting config.Setting) []*entity.CheckReport {
	reports := make([]*entity.CheckReport, 0, len(uids))
	for _, uid := range uids {
		r, err := Check(uid, setting)
		if err != nil {
			monitor.GlobalEventBus.AddLog(uid, "❌ 质检失败: "+err.Error())
			continue
		}
		reports = append(reports, r)
	}
	return reports
}

// CheckAll 质检全部账号
func CheckAll(setting config.Setting) []*entity.CheckReport {
	users, err := dao.QueryAllUsers()
	if err != nil {
		return nil
	}
	uids := make([]string, 0, len(users))
	for _, u := range users {
		uids = append(uids, u.Uid)
	}
	return CheckMany(uids, setting)
}

// GetStored 读取账号最近一次质检报告
func GetStored(uid string) *entity.CheckReport {
	raw, err := dao.GetCheckReport(uid)
	if err != nil || raw == "" {
		return nil
	}
	var report entity.CheckReport
	if json.Unmarshal([]byte(raw), &report) != nil {
		return nil
	}
	return &report
}

// GetAllStored 读取全部账号最近一次质检报告
func GetAllStored() []*entity.CheckReport {
	list, err := dao.GetAllCheckReports()
	if err != nil {
		return nil
	}
	reports := make([]*entity.CheckReport, 0, len(list))
	for _, po := range list {
		var report entity.CheckReport
		if json.Unmarshal([]byte(po.Summary), &report) == nil {
			reports = append(reports, &report)
		}
	}
	return reports
}
