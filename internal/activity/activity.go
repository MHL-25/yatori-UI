package activity

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"yatori-UI/internal/config"
	"yatori-UI/internal/entity"
	"yatori-UI/internal/email"
	"yatori-UI/internal/monitor"

	"github.com/thedevsaddam/gojsonq"
	"github.com/yatori-dev/yatori-go-core/aggregation/cqie"
	"github.com/yatori-dev/yatori-go-core/aggregation/enaea"
	"github.com/yatori-dev/yatori-go-core/aggregation/haiqikeji"
	"github.com/yatori-dev/yatori-go-core/aggregation/icve"
	"github.com/yatori-dev/yatori-go-core/aggregation/ketangx"
	"github.com/yatori-dev/yatori-go-core/aggregation/qingshuxuetang"
	"github.com/yatori-dev/yatori-go-core/aggregation/welearn"
	"github.com/yatori-dev/yatori-go-core/aggregation/xuexitong"
	"github.com/yatori-dev/yatori-go-core/aggregation/xuexitong/point"
	"github.com/yatori-dev/yatori-go-core/aggregation/yinghua"
	cqieApi "github.com/yatori-dev/yatori-go-core/api/cqie"
	enaeaApi "github.com/yatori-dev/yatori-go-core/api/enaea"
	hqkjApi "github.com/yatori-dev/yatori-go-core/api/haiqikeji"
	icveApi "github.com/yatori-dev/yatori-go-core/api/icve"
	ketangxApi "github.com/yatori-dev/yatori-go-core/api/ketangx"
	qsxtApi "github.com/yatori-dev/yatori-go-core/api/qingshuxuetang"
	welearnApi "github.com/yatori-dev/yatori-go-core/api/welearn"
	xuexitongApi "github.com/yatori-dev/yatori-go-core/api/xuexitong"
	yinghuaApi "github.com/yatori-dev/yatori-go-core/api/yinghua"
	"github.com/yatori-dev/yatori-go-core/models/ctype"
	"github.com/yatori-dev/yatori-go-core/que-core/aiq"
	"github.com/yatori-dev/yatori-go-core/que-core/external"
	"github.com/yatori-dev/yatori-go-core/utils"
	"github.com/yatori-dev/yatori-go-core/utils/qutils"
)

type Activity interface {
	Login() error
	Start() error
	Stop() error
	Pause() error
	GetUserCache() interface{}
	SetUser(user config.User)
	GetUser() config.User
	GetUid() string
	IsRunning() bool
	GetSetting() config.Setting
	SetSetting(setting config.Setting)
}

type UserActivityBase struct {
	Uid       string         `json:"uid"`
	User      config.User    `json:"user"`
	Setting   config.Setting `json:"setting"`
	Running   bool           `json:"running"`
	Stopped   bool           `json:"stopped"`
	Paused    bool           `json:"paused"`
	UserCache interface{}    `json:"-"`
	mu        sync.Mutex
}

func (u *UserActivityBase) SetUser(user config.User)    { u.User = user }
func (u *UserActivityBase) GetUser() config.User        { return u.User }
func (u *UserActivityBase) GetUid() string              { return u.Uid }
func (u *UserActivityBase) GetUserCache() interface{}   { return u.UserCache }
func (u *UserActivityBase) GetSetting() config.Setting  { return u.Setting }
func (u *UserActivityBase) SetSetting(s config.Setting) { u.Setting = s }

func (u *UserActivityBase) IsRunning() bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.Running
}

func (u *UserActivityBase) setRunning(running bool) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.Running = running
}

func (u *UserActivityBase) isStopped() bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.Stopped
}

func (u *UserActivityBase) setStopped(stopped bool) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.Stopped = stopped
}

func (u *UserActivityBase) isPaused() bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.Paused
}

func (u *UserActivityBase) setPaused(paused bool) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.Paused = paused
}

func (u *UserActivityBase) sendCompletionEmail() {
	setting := u.Setting
	user := u.User
	if setting.EmailInform.Sw == 1 && len(user.InformEmails) > 0 {
		platformName := config.PlatformNames[user.AccountType]
		content := fmt.Sprintf("账号：[%s]</br>平台：[%s]</br>通知：所有课程已执行完毕", user.Account, platformName)
		go email.SendMail(setting.EmailInform.SMTPHost, setting.EmailInform.SMTPPort, setting.EmailInform.UserName, setting.EmailInform.Password, user.InformEmails, content)
	}
}

var activityMap = struct {
	sync.RWMutex
	m map[string]Activity
}{m: make(map[string]Activity)}

func GetActivity(uid string) Activity {
	activityMap.RLock()
	defer activityMap.RUnlock()
	return activityMap.m[uid]
}

func PutActivity(uid string, act Activity) {
	activityMap.Lock()
	defer activityMap.Unlock()
	activityMap.m[uid] = act
}

func BuildActivity(userPO entity.UserPO, setting config.Setting) Activity {
	user := config.User{}
	err := json.Unmarshal([]byte(userPO.UserConfigJson), &user)
	if err != nil {
		return nil
	}

	base := UserActivityBase{
		Uid:     userPO.Uid,
		User:    user,
		Setting: setting,
	}

	monitor.GlobalEventBus.InitProgress(userPO.Uid, userPO.DisplayName(), user.AccountType)

	switch user.AccountType {
	case "XUEXITONG":
		return &XXTActivity{UserActivityBase: base}
	case "YINGHUA", "CANGHUI":
		return &YingHuaActivity{UserActivityBase: base}
	case "ENAEA":
		return &EnaeaActivity{UserActivityBase: base}
	case "CQIE":
		return &CqieActivity{UserActivityBase: base}
	case "ICVE":
		return &IcveActivity{UserActivityBase: base}
	case "QSXT":
		return &QsxtActivity{UserActivityBase: base}
	case "WELEARN":
		return &WeLearnActivity{UserActivityBase: base}
	case "KETANGX":
		return &KetangxActivity{UserActivityBase: base}
	case "HQKJ":
		return &HqkjActivity{UserActivityBase: base}
	default:
		return nil
	}
}

// ==================== XUEXITONG ====================

type XXTActivity struct {
	UserActivityBase
	model3Caches []xuexitongApi.XueXiTUserCache
}

func (a *XXTActivity) Login() error {
	cache := &xuexitongApi.XueXiTUserCache{Name: a.User.Account, Password: a.User.Password}
	if a.User.IsProxy == 1 {
		cache.IpProxySW = true
	}
	var loginError error
	if len(a.User.Password) < 50 {
		loginError = xuexitong.XueXiTLoginAction(cache)
	} else {
		loginError = xuexitong.XueXiTCookieLoginAction(cache)
	}
	if loginError != nil {
		monitor.GlobalEventBus.SetError(a.Uid, "❌ 登录失败: "+loginError.Error())
		return loginError
	}
	a.UserCache = cache
	monitor.GlobalEventBus.AddLog(a.Uid, "🔑 登录成功: "+cache.Name)
	return nil
}

func (a *XXTActivity) Start() error {
	a.setStopped(false)
	a.setPaused(false)
	a.setRunning(true)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusRunning)
	monitor.GlobalEventBus.AddLog(a.Uid, "🚀 开始刷课任务...")
	go a.run()
	return nil
}

func (a *XXTActivity) Stop() error {
	a.setStopped(true)
	a.setRunning(false)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusStopped)
	monitor.GlobalEventBus.AddLog(a.Uid, "⏹️ 任务已停止")
	return nil
}

func (a *XXTActivity) Pause() error {
	a.setRunning(false)
	a.setPaused(true)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusPaused)
	monitor.GlobalEventBus.AddLog(a.Uid, "⏸️ 任务已暂停")
	return nil
}

func (a *XXTActivity) run() {
	defer func() {
		a.setRunning(false)
		if r := recover(); r != nil {
			monitor.GlobalEventBus.SetError(a.Uid, fmt.Sprintf("Panic: %v", r))
			return
		}
		if a.isPaused() {
			return
		}
		if !a.isStopped() {
			monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusCompleted)
			a.sendCompletionEmail()
			monitor.GlobalEventBus.UpdateProgress(a.Uid, 100, "\u6240\u6709\u8bfe\u7a0b\u5b66\u4e60\u5b8c\u6bd5")
		}
	}()

	if a.UserCache == nil {
		monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusLogging)
		monitor.GlobalEventBus.AddLog(a.Uid, "🔑 正在登录...")
		if err := a.Login(); err != nil {
			return
		}
	}

	cache := a.UserCache.(*xuexitongApi.XueXiTUserCache)

	user := a.User
	if user.CoursesCustom.VideoModel == 3 {
		num := 3
		if user.CoursesCustom.CxNode != nil && *user.CoursesCustom.CxNode != 0 {
			num = *user.CoursesCustom.CxNode
		}
		if user.CoursesCustom.CxNode != nil && *user.CoursesCustom.CxNode == -1 {
			monitor.GlobalEventBus.AddLog(a.Uid, "⚠️ 警告: 无限制并发模式，可能触发封号!")
		} else {
			monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("🔄 多任务点模式: 同时 %d 个任务点并发", num))
		}
		a.model3Caches = make([]xuexitongApi.XueXiTUserCache, 0, num)
		for i := 0; i < num; i++ {
			a.model3Caches = append(a.model3Caches, *cache)
			if i > 0 {
				xuexitong.ReLogin(&a.model3Caches[i])
				monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("🔑 并发缓存登录 %d/%d", i+1, num))
				time.Sleep(1 * time.Second)
			}
		}
	}

	courseList, err := xuexitong.XueXiTPullCourseAction(cache)
	if err != nil {
		monitor.GlobalEventBus.SetError(a.Uid, "❌ 拉取课程列表失败: "+err.Error())
		return
	}

	totalCourses := len(courseList)
	monitor.GlobalEventBus.UpdateCourseProgress(a.Uid, 0, totalCourses)
	monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("📚 获取到 %d 门课程，开始学习", totalCourses))

	completedCourses := 0
	var nodesLock sync.WaitGroup
	for _, course := range courseList {
		if !a.IsRunning() {
			break
		}
		completedCourses++
		monitor.GlobalEventBus.UpdateCourseProgress(a.Uid, completedCourses, totalCourses)
		if user.CoursesCustom.VideoModel == 1 {
			a.courseStudy(&course)
		} else {
			nodesLock.Add(1)
			go func(course xuexitong.XueXiTCourse) {
				defer nodesLock.Done()
				a.courseStudy(&course)
			}(course)
		}
	}
	nodesLock.Wait()
	monitor.GlobalEventBus.AddLog(a.Uid, "✅ 所有待学习课程学习完毕")
}

func (a *XXTActivity) courseStudy(courseItem *xuexitong.XueXiTCourse) bool {
	if !a.IsRunning() {
		return false
	}
	user := a.User
	cache := a.UserCache.(*xuexitongApi.XueXiTUserCache)
	setting := a.Setting

	if len(user.CoursesCustom.ExcludeCourses) != 0 && config.CmpCourse(courseItem.CourseName, user.CoursesCustom.ExcludeCourses) {
		monitor.GlobalEventBus.AddLog(a.Uid, "⏭️ 跳过课程(已排除): "+courseItem.CourseName)
		return false
	}
	if len(user.CoursesCustom.IncludeCourses) != 0 && !config.CmpCourse(courseItem.CourseName, user.CoursesCustom.IncludeCourses) {
		monitor.GlobalEventBus.AddLog(a.Uid, "⏭️ 跳过课程(不在包含列表): "+courseItem.CourseName)
		return false
	}
	if !courseItem.IsStart {
		monitor.GlobalEventBus.AddLog(a.Uid, "⏭️ 跳过课程(未开始): "+courseItem.CourseName)
		return false
	}

	monitor.GlobalEventBus.AddLog(a.Uid, "📖 开始学习课程: "+courseItem.CourseName)
	monitor.GlobalEventBus.UpdateProgress(a.Uid, 0, "学习: "+courseItem.CourseName)

	a.chapterStudy(setting, cache, courseItem)
	a.writeCourseWorkAndExam(setting, cache, courseItem)
	monitor.GlobalEventBus.AddLog(a.Uid, "✅ 课程学习完毕: "+courseItem.CourseName)
	return true
}

func (a *XXTActivity) chapterStudy(setting config.Setting, cache *xuexitongApi.XueXiTUserCache, courseItem *xuexitong.XueXiTCourse) {
	if courseItem.JobRate >= 100 || courseItem.State == 1 {
		monitor.GlobalEventBus.AddLog(a.Uid, "⏭️ 课程已完成，跳过: "+courseItem.CourseName)
		return
	}
	user := a.User
	key, _ := strconv.Atoi(courseItem.Key)
	action, _, err := xuexitong.PullCourseChapterAction(cache, courseItem.Cpi, key)
	if err != nil {
		monitor.GlobalEventBus.AddLog(a.Uid, "⚠️ 拉取章节失败: "+err.Error())
		return
	}

	if user.CoursesCustom.ShuffleSw == 1 {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(action.Knowledge), func(i, j int) {
			action.Knowledge[i], action.Knowledge[j] = action.Knowledge[j], action.Knowledge[i]
		})
	}

	if len(action.Knowledge) == 0 {
		monitor.GlobalEventBus.AddLog(a.Uid, "⏭️ 该课程章节为空，已自动跳过")
		return
	}

	monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("📑 获取到 %d 个章节", len(action.Knowledge)))

	var nodes []int
	for _, item := range action.Knowledge {
		nodes = append(nodes, item.ID)
	}

	courseId, _ := strconv.Atoi(courseItem.CourseID)
	userId, _ := strconv.Atoi(cache.UserID)
	pointAction, err := xuexitong.ChapterFetchPointAction(cache, nodes, &action, key, userId, courseItem.Cpi, courseId)
	if err != nil {
		monitor.GlobalEventBus.AddLog(a.Uid, "⚠️ 获取任务点失败: "+err.Error())
		return
	}

	isFinished := func(index int) bool {
		if index < 0 || index >= len(pointAction.Knowledge) {
			return false
		}
		i := pointAction.Knowledge[index]
		if i.PointTotal == 0 && i.PointFinished == 0 {
			xuexitong.EnterChapterForwardCallAction(cache, strconv.Itoa(courseId), strconv.Itoa(key), strconv.Itoa(pointAction.Knowledge[index].ID), strconv.Itoa(courseItem.Cpi))
		}
		return i.PointTotal >= 0 && i.PointTotal == i.PointFinished
	}

	if user.CoursesCustom.VideoModel == 3 && len(a.model3Caches) > 0 {
		queue := make(chan int, len(a.model3Caches))
		for i := 0; i < len(a.model3Caches); i++ {
			queue <- i
		}
		var nodeLock sync.WaitGroup
		for index := range nodes {
			if !a.IsRunning() {
				break
			}
			if isFinished(index) {
				continue
			}
			if user.CoursesCustom.CxNode != nil && *user.CoursesCustom.CxNode == -1 {
				nodeLock.Add(1)
				go func(index int) {
					defer nodeLock.Done()
					resCache := *cache
					xuexitong.ReLogin(&resCache)
					a.nodeRun(setting, &resCache, courseItem, pointAction, action, nodes, index, key, courseId)
				}(index)
				time.Sleep(1 * time.Second)
			} else {
				idx := <-queue
				nodeLock.Add(1)
				go func(idx int, index int) {
					defer nodeLock.Done()
					defer func() { queue <- idx }()
					a.nodeRun(setting, &a.model3Caches[idx], courseItem, pointAction, action, nodes, index, key, courseId)
				}(idx, index)
			}
		}
		nodeLock.Wait()
	} else {
		for index := range nodes {
			if !a.IsRunning() {
				return
			}
			if isFinished(index) {
				continue
			}
			a.nodeRun(setting, cache, courseItem, pointAction, action, nodes, index, key, courseId)
		}
	}
}

func (a *XXTActivity) nodeRun(setting config.Setting, cache *xuexitongApi.XueXiTUserCache, courseItem *xuexitong.XueXiTCourse,
	pointAction xuexitong.ChaptersList, action xuexitong.ChaptersList, nodes []int, index int, key int, courseId int) {

	_, fetchCards, err1 := xuexitong.ChapterFetchCardsAction(cache, &action, nodes, index, courseId, key, courseItem.Cpi)
	if err1 != nil {
		monitor.GlobalEventBus.AddLog(a.Uid, "⚠️ 获取卡片失败: "+err1.Error())
		return
	}

	videoDTOs, workDTOs, documentDTOs, hyperlinkDTOs, liveDTOs, bbsDTOs := xuexitongApi.ParsePointDto(fetchCards)
	if videoDTOs == nil && workDTOs == nil && documentDTOs == nil && hyperlinkDTOs == nil && liveDTOs == nil && bbsDTOs == nil {
		return
	}

	user := a.User

	if videoDTOs != nil && user.CoursesCustom.VideoModel != 0 {
		for i := range videoDTOs {
			if !a.IsRunning() {
				return
			}
			videoDTO := &videoDTOs[i]
			card, enc, err2 := xuexitong.PageMobileChapterCardAction(cache, key, courseId, videoDTO.KnowledgeID, videoDTO.CardIndex, courseItem.Cpi)
			if err2 != nil {
				if strings.Contains(err2.Error(), "章节未开放") {
					monitor.GlobalEventBus.AddLog(a.Uid, "⏭️ 章节未开放，跳过: "+videoDTO.Title)
					break
				}
				if strings.Contains(err2.Error(), "没有历史人脸") {
					monitor.GlobalEventBus.AddLog(a.Uid, "🚫 过人脸失败，该账号可能从未进行过人脸识别，请先进行一次人脸识别后再试")
					return
				}
				if strings.Contains(err2.Error(), "活体检测失败") {
					monitor.GlobalEventBus.AddLog(a.Uid, "🚫 过人脸失败，该账号所录入的人脸可能并不规范")
					return
				}
				monitor.GlobalEventBus.AddLog(a.Uid, "⚠️ 视频卡片获取失败，跳过: "+videoDTO.Title+" - "+err2.Error())
				continue
			}
			videoDTO.AttachmentsDetection(card)
			if !videoDTO.IsJob {
				monitor.GlobalEventBus.AddLog(a.Uid, "⏭️ 非任务点视频，跳过: "+videoDTO.Title)
				continue
			}
			videoDTO.Enc = enc
			if videoDTO.IsPassed {
				monitor.GlobalEventBus.AddLog(a.Uid, "⏭️ 视频已完成，跳过: "+videoDTO.Title)
				continue
			}
			if videoDTO.Attachment == nil && videoDTO.JobID == "" && videoDTO.Duration <= videoDTO.PlayTime {
				monitor.GlobalEventBus.AddLog(a.Uid, "⏭️ 视频已播放完毕，跳过: "+videoDTO.Title)
				continue
			}
			monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("🎬 开始观看视频: %s (%d/%ds)", videoDTO.Title, videoDTO.PlayTime, videoDTO.Duration))
			if videoDTO.Type == ctype.Video {
				a.executeVideo(cache, courseItem, pointAction.Knowledge[index], videoDTO, key, courseItem.Cpi)
			} else if videoDTO.Type == ctype.InsertAudio {
				a.executeAudio(cache, courseItem, pointAction.Knowledge[index], videoDTO, key, courseItem.Cpi)
			}
			randSleepTime := rand.Intn(51) + 10
			time.Sleep(time.Duration(randSleepTime) * time.Second)
		}
	}

	if documentDTOs != nil && user.CoursesCustom.VideoModel != 0 {
		for i := range documentDTOs {
			if !a.IsRunning() {
				return
			}
			documentDTO := &documentDTOs[i]
			card, _, err2 := xuexitong.PageMobileChapterCardAction(cache, key, courseId, documentDTO.KnowledgeID, documentDTO.CardIndex, courseItem.Cpi)
			if err2 != nil {
				if strings.Contains(err2.Error(), "章节未开放") {
					monitor.GlobalEventBus.AddLog(a.Uid, "⏭️ 该章节未开放，跳过")
					break
				}
				if strings.Contains(err2.Error(), "没有历史人脸") {
					monitor.GlobalEventBus.AddLog(a.Uid, "🚫 过人脸失败，该账号可能从未进行过人脸识别")
					return
				}
				monitor.GlobalEventBus.AddLog(a.Uid, "⚠️ 文档卡片获取失败，跳过: "+err2.Error())
				continue
			}
			documentDTO.AttachmentsDetection(card)
			if !documentDTO.IsJob {
				continue
			}
			point.ExecuteDocument(cache, documentDTO)
			time.Sleep(5 * time.Second)
		}
	}

	if workDTOs != nil && user.CoursesCustom.AutoExam != 0 && user.CoursesCustom.CxChapterTestSw != nil && *user.CoursesCustom.CxChapterTestSw == 1 {
		for i := range workDTOs {
			if !a.IsRunning() {
				return
			}
			workDTO := &workDTOs[i]
			mobileCard, _, err2 := xuexitong.PageMobileChapterCardAction(cache, key, courseId, workDTO.KnowledgeID, workDTO.CardIndex, courseItem.Cpi)
			if err2 != nil {
				if strings.Contains(err2.Error(), "章节未开放") {
					monitor.GlobalEventBus.AddLog(a.Uid, "⏭️ 该章节未开放，跳过")
					continue
				}
				if strings.Contains(err2.Error(), "没有历史人脸") {
					monitor.GlobalEventBus.AddLog(a.Uid, "🚫 过人脸失败，该账号可能从未进行过人脸识别")
					return
				}
				monitor.GlobalEventBus.AddLog(a.Uid, "⚠️ 作业卡片获取失败，跳过: "+err2.Error())
				continue
			}
			flag, _ := workDTO.AttachmentsDetection(mobileCard)
			questionAction, err2 := xuexitong.ParseWorkQuestionAction(cache, workDTO)
			if err2 != nil {
				if strings.Contains(err2.Error(), "已截止，不能作答") {
					monitor.GlobalEventBus.AddLog(a.Uid, "⏭️ 作业已截止，跳过: "+questionAction.Title)
					continue
				}
				continue
			}
			if !flag {
				monitor.GlobalEventBus.AddLog(a.Uid, "⏭️ 作业已完成，跳过: "+questionAction.Title)
				continue
			}
			if len(questionAction.Short) == 0 && len(questionAction.Choice) == 0 &&
				len(questionAction.Judge) == 0 && len(questionAction.Fill) == 0 &&
				len(questionAction.TermExplanation) == 0 && len(questionAction.Essay) == 0 {
				monitor.GlobalEventBus.AddLog(a.Uid, "⏭️ 作业无题目，跳过: "+questionAction.Title)
				continue
			}
			a.chapterTestAction(cache, questionAction, courseItem, pointAction.Knowledge[index])
		}
	}

	if hyperlinkDTOs != nil && user.CoursesCustom.VideoModel != 0 {
		for i := range hyperlinkDTOs {
			if !a.IsRunning() {
				return
			}
			hyperlinkDTO := &hyperlinkDTOs[i]
			card, _, err2 := xuexitong.PageMobileChapterCardAction(cache, key, courseId, hyperlinkDTO.KnowledgeID, hyperlinkDTO.CardIndex, courseItem.Cpi)
			if err2 != nil {
				if strings.Contains(err2.Error(), "章节未开放") {
					monitor.GlobalEventBus.AddLog(a.Uid, "⏭️ 该章节未开放，跳过")
					continue
				}
				if strings.Contains(err2.Error(), "没有历史人脸") {
					monitor.GlobalEventBus.AddLog(a.Uid, "🚫 过人脸失败，该账号可能从未进行过人脸识别")
					return
				}
				monitor.GlobalEventBus.AddLog(a.Uid, "⚠️ 外链卡片获取失败，跳过: "+err2.Error())
				continue
			}
			hyperlinkDTO.AttachmentsDetection(card)
			point.ExecuteHyperlink(cache, hyperlinkDTO)
			time.Sleep(5 * time.Second)
		}
	}

	if liveDTOs != nil && user.CoursesCustom.VideoModel != 0 {
		for i := range liveDTOs {
			if !a.IsRunning() {
				return
			}
			liveDTO := &liveDTOs[i]
			card, _, err2 := xuexitong.PageMobileChapterCardAction(cache, key, courseId, liveDTO.KnowledgeID, liveDTO.CardIndex, courseItem.Cpi)
			if err2 != nil {
				if strings.Contains(err2.Error(), "章节未开放") {
					monitor.GlobalEventBus.AddLog(a.Uid, "⏭️ 该章节未开放，跳过")
					continue
				}
				if strings.Contains(err2.Error(), "没有历史人脸") {
					monitor.GlobalEventBus.AddLog(a.Uid, "🚫 过人脸失败，该账号可能从未进行过人脸识别")
					return
				}
				monitor.GlobalEventBus.AddLog(a.Uid, "⚠️ 直播卡片获取失败，跳过: "+err2.Error())
				continue
			}
			liveDTO.AttachmentsDetection(card)
			if !liveDTO.IsJob {
				continue
			}
			point.PullLiveInfoAction(cache, liveDTO)
			if liveDTO.LiveStatusCode == 0 {
				monitor.GlobalEventBus.AddLog(a.Uid, "⏭️ 直播未开播，跳过: "+liveDTO.Title)
				continue
			}
			relationReport, err2 := point.LiveCreateRelationAction(cache, liveDTO)
			if err2 != nil {
				monitor.GlobalEventBus.AddLog(a.Uid, "⚠️ 直播建立联系失败: "+relationReport+" "+err2.Error())
			} else {
				monitor.GlobalEventBus.AddLog(a.Uid, "📡 直播建立联系成功: "+relationReport)
			}
			for {
				if !a.IsRunning() {
					return
				}
				report, err := point.ExecuteLive(cache, liveDTO)
				point.PullLiveInfoAction(cache, liveDTO)
				if err != nil {
					monitor.GlobalEventBus.AddLog(a.Uid, "⚠️ 直播提交异常: "+report+" "+err.Error())
				}
				if strings.Contains(report, "@success") {
					monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("直播任务点状态: %s 进度: %.2f%%", report, liveDTO.VideoCompletePercent))
				}
				if liveDTO.VideoCompletePercent >= 90 {
					monitor.GlobalEventBus.AddLog(a.Uid, "✅ 直播任务点已完成: "+liveDTO.Title)
					break
				}
				time.Sleep(30 * time.Second)
			}
		}
	}

	if bbsDTOs != nil && user.CoursesCustom.AutoExam != 0 {
		for i := range bbsDTOs {
			if !a.IsRunning() {
				return
			}
			bbsDTO := &bbsDTOs[i]
			card, _, err2 := xuexitong.PageMobileChapterCardAction(cache, key, courseId, bbsDTO.KnowledgeID, bbsDTO.CardIndex, courseItem.Cpi)
			if err2 != nil {
				if strings.Contains(err2.Error(), "章节未开放") {
					monitor.GlobalEventBus.AddLog(a.Uid, "⏭️ 该章节未开放，跳过")
					continue
				}
				if strings.Contains(err2.Error(), "没有历史人脸") {
					monitor.GlobalEventBus.AddLog(a.Uid, "🚫 过人脸失败，该账号可能从未进行过人脸识别")
					return
				}
				monitor.GlobalEventBus.AddLog(a.Uid, "⚠️ 讨论卡片获取失败，跳过: "+err2.Error())
				continue
			}
			bbsDTO.AttachmentsDetection(card)
			if !bbsDTO.IsJob {
				continue
			}
			bbsTopic, err1 := point.PullPhoneBbsInfoAction(cache, bbsDTO)
			if bbsTopic == nil {
				monitor.GlobalEventBus.AddLog(a.Uid, "⏭️ 无法拉取讨论主题，跳过: "+bbsDTO.Title)
				continue
			}
			if err1 != nil {
				monitor.GlobalEventBus.AddLog(a.Uid, "⚠️ 讨论拉取异常: "+err1.Error())
			}
			var report string
			var err error
			if user.CoursesCustom.AutoExam == 1 {
				report, err = bbsTopic.AIAnswer(cache, bbsDTO, setting.AiSetting.AiUrl, setting.AiSetting.Model, setting.AiSetting.AiType, setting.AiSetting.APIKEY)
			} else if user.CoursesCustom.AutoExam == 2 {
				report, err = bbsTopic.ExternalAnswer(cache, bbsDTO, setting.ApiQueSetting.Url)
			} else if user.CoursesCustom.AutoExam == 3 {
				report, err = bbsTopic.XXTAIAnswer(cache, bbsDTO)
			}
			if err != nil {
				monitor.GlobalEventBus.AddLog(a.Uid, "⚠️ 讨论提交异常: "+report+" "+err.Error())
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func (a *XXTActivity) executeVideo(cache *xuexitongApi.XueXiTUserCache, courseItem *xuexitong.XueXiTCourse, knowledgeItem xuexitong.KnowledgeItem, p *xuexitongApi.PointVideoDto, key, courseCpi int) {
	if state, _ := xuexitong.VideoDtoFetchAction(cache, p); state {
		playingTime := p.PlayTime
		if !p.IsPassed && p.PlayTime == p.Duration {
			playingTime = 0
		}
		overTime := 0
		selectSec := 58
		extendSec := 5
		limitTime := max(500, p.Duration/2)
		mode := 1
		for {
			if !a.IsRunning() {
				return
			}
			var playReport string
			var err error
			if playingTime != p.Duration {
				if playingTime == p.PlayTime {
					playReport, err = xuexitong.VideoSubmitStudyTimeAction(cache, p, playingTime, mode, 3)
				} else {
					playReport, err = xuexitong.VideoSubmitStudyTimeAction(cache, p, playingTime, mode, 0)
				}
			} else {
				playReport, err = xuexitong.VideoSubmitStudyTimeAction(cache, p, playingTime, mode, 0)
			}
			if err != nil {
				if strings.Contains(err.Error(), "failed to fetch video, status code: 500") {
					monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频触发风控500，跳过: %s - %s", p.Title, err.Error()))
					break
				}
				if strings.Contains(err.Error(), "failed to fetch video, status code: 403") {
					if mode == 1 {
						mode = 0
						monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频触发403，切换为Web端: %s", p.Title))
						continue
					}
					monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频触发403(Web端)，尝试绕过人脸: %s", p.Title))
					pullJson, img, err2 := cache.GetHistoryFaceImg("")
					if err2 != nil {
						monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("上传人脸失败，跳过该视频: %s - %s %s", p.Title, pullJson, err2))
						return
					}
					disturbImage := utils.ImageRGBDisturb(img)
					_, _, _, _, errPass := xuexitong.PassFacePCAction(cache, p.CourseID, p.ClassID, p.Cpi, fmt.Sprintf("%d", p.KnowledgeID), p.Enc, p.JobID, p.ObjectID, p.Mid, p.RandomCaptureTime, disturbImage)
					if errPass != nil {
						monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("绕过人脸失败: %s - %s", p.Title, errPass.Error()))
					} else {
						monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("绕过人脸成功: %s", p.Title))
					}
					time.Sleep(5 * time.Second)
					continue
				}
				if strings.Contains(err.Error(), "failed to fetch video, status code: 404") {
					time.Sleep(10 * time.Second)
					continue
				}
			}
			if err != nil {
				monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频提交学时异常，跳过: %s - %s", p.Title, err.Error()))
				break
			}
			if gojsonq.New().JSONString(playReport).Find("isPassed") == nil {
				monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频提交学时返回异常，跳过: %s - %s", p.Title, playReport))
				break
			}
			outTimeMsg := gojsonq.New().JSONString(playReport).Find("OutTimeMsg")
			if outTimeMsg != nil {
				if msg, ok := outTimeMsg.(string); ok && msg == "观看时长超过阈值" {
					monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频超时完成: %s %d/%d %.2f%%", p.Title, p.Duration, p.Duration, float32(p.Duration)/float32(p.Duration)*100))
					break
				}
			}
			isPassed, ok := gojsonq.New().JSONString(playReport).Find("isPassed").(bool)
			if ok && isPassed && playingTime >= p.Duration {
				if overTime == 0 {
					monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频完成: %s %d/%d %.2f%%", p.Title, p.Duration, p.Duration, float32(p.Duration)/float32(p.Duration)*100))
				} else {
					monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频过超提交成功: %s %d/%d 过超:%d/%d %.2f%%", p.Title, p.Duration, p.Duration, overTime, limitTime, float32(p.Duration)/float32(p.Duration)*100))
				}
				break
			}
			if overTime == 0 {
				monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频: %s %d/%d %.2f%%", p.Title, playingTime, p.Duration, float32(playingTime)/float32(p.Duration)*100))
			} else {
				monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频(过超): %s %d/%d 过超:%d/%d %.2f%%", p.Title, playingTime, p.Duration, overTime, limitTime, float32(playingTime)/float32(p.Duration)*100))
			}
			if overTime >= limitTime {
				monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频过超提交失败，跳过: %s", p.Title))
				break
			}
			if p.Duration-playingTime < selectSec && p.Duration != playingTime {
				remainingTime := p.Duration - playingTime
				playingTime = p.Duration
				time.Sleep(time.Duration(remainingTime) * time.Second)
			} else if p.Duration == playingTime {
				if p.JobID == "" && p.Attachment == nil {
					monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("非任务点视频播放完毕: %s", p.Title))
					break
				}
				overTime += extendSec
				time.Sleep(time.Duration(extendSec) * time.Second)
			} else {
				playingTime = playingTime + selectSec
				time.Sleep(time.Duration(selectSec) * time.Second)
			}
		}
	} else {
		monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频任务点解析失败，可能是视频本身问题，跳过: %s", p.Title))
	}
}

func (a *XXTActivity) executeAudio(cache *xuexitongApi.XueXiTUserCache, courseItem *xuexitong.XueXiTCourse, knowledgeItem xuexitong.KnowledgeItem, p *xuexitongApi.PointVideoDto, key, courseCpi int) {
	if state, _ := xuexitong.VideoDtoFetchAction(cache, p); state {
		playingTime := p.PlayTime
		if !p.IsPassed && p.PlayTime == p.Duration {
			playingTime = 0
		}
		overTime := 0
		selectSec := 58
		extendSec := 5
		limitTime := max(500, p.Duration/2)
		mode := 1
		for {
			if !a.IsRunning() {
				return
			}
			var playReport string
			var err error
			if playingTime != p.Duration {
				if playingTime == p.PlayTime {
					playReport, err = xuexitong.AudioSubmitStudyTimeAction(cache, p, playingTime, mode, 3)
				} else {
					playReport, err = xuexitong.AudioSubmitStudyTimeAction(cache, p, playingTime, mode, 0)
				}
			} else {
				playReport, err = xuexitong.AudioSubmitStudyTimeAction(cache, p, playingTime, mode, 0)
			}
			if err != nil {
				if strings.Contains(err.Error(), "failed to fetch video, status code: 500") {
					monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("音频触发风控500，跳过: %s - %s", p.Title, err.Error()))
					break
				}
				if strings.Contains(err.Error(), "failed to fetch video, status code: 403") {
					if mode == 1 {
						mode = 0
						monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("音频触发403，切换为Web端: %s", p.Title))
						continue
					}
					monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("音频触发403(Web端)，尝试绕过人脸: %s", p.Title))
					pullJson, img, err2 := cache.GetHistoryFaceImg("")
					if err2 != nil {
						monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("上传人脸失败，跳过该音频: %s - %s %s", p.Title, pullJson, err2))
						return
					}
					disturbImage := utils.ImageRGBDisturb(img)
					_, _, _, _, errPass := xuexitong.PassFacePCAction(cache, p.CourseID, p.ClassID, p.Cpi, fmt.Sprintf("%d", p.KnowledgeID), p.Enc, p.JobID, p.ObjectID, p.Mid, p.RandomCaptureTime, disturbImage)
					if errPass != nil {
						monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("绕过人脸失败: %s - %s", p.Title, errPass.Error()))
					} else {
						monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("绕过人脸成功: %s", p.Title))
					}
					time.Sleep(5 * time.Second)
					continue
				}
				if strings.Contains(err.Error(), "failed to fetch video, status code: 404") {
					time.Sleep(10 * time.Second)
					continue
				}
			}
			if err != nil {
				monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("音频提交学时异常，跳过: %s - %s", p.Title, err.Error()))
				break
			}
			if gojsonq.New().JSONString(playReport).Find("isPassed") == nil {
				monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("音频提交学时返回异常，跳过: %s - %s", p.Title, playReport))
				break
			}
			outTimeMsg := gojsonq.New().JSONString(playReport).Find("OutTimeMsg")
			if outTimeMsg != nil {
				if msg, ok := outTimeMsg.(string); ok && msg == "观看时长超过阈值" {
					monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("音频超时完成: %s %d/%d %.2f%%", p.Title, p.Duration, p.Duration, float32(p.Duration)/float32(p.Duration)*100))
					break
				}
			}
			isPassed, ok := gojsonq.New().JSONString(playReport).Find("isPassed").(bool)
			if ok && isPassed && playingTime >= p.Duration {
				if overTime == 0 {
					monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("音频完成: %s %d/%d %.2f%%", p.Title, p.Duration, p.Duration, float32(p.Duration)/float32(p.Duration)*100))
				} else {
					monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("音频过超提交成功: %s %d/%d 过超:%d/%d %.2f%%", p.Title, p.Duration, p.Duration, overTime, limitTime, float32(p.Duration)/float32(p.Duration)*100))
				}
				break
			}
			if overTime == 0 {
				monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("音频: %s %d/%d %.2f%%", p.Title, playingTime, p.Duration, float32(playingTime)/float32(p.Duration)*100))
			} else {
				monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("音频(过超): %s %d/%d 过超:%d/%d %.2f%%", p.Title, playingTime, p.Duration, overTime, limitTime, float32(playingTime)/float32(p.Duration)*100))
			}
			if overTime >= limitTime {
				monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("音频过超提交失败，跳过: %s", p.Title))
				break
			}
			if p.Duration-playingTime < selectSec && p.Duration != playingTime {
				remainingTime := p.Duration - playingTime
				playingTime = p.Duration
				time.Sleep(time.Duration(remainingTime) * time.Second)
			} else if p.Duration == playingTime {
				if p.JobID == "" && p.Attachment == nil {
					monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("非任务点音频播放完毕: %s", p.Title))
					break
				}
				overTime += extendSec
				time.Sleep(time.Duration(extendSec) * time.Second)
			} else {
				playingTime = playingTime + selectSec
				time.Sleep(time.Duration(selectSec) * time.Second)
			}
		}
	} else {
		monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("音频任务点解析失败，可能是音频本身问题，跳过: %s", p.Title))
	}
}

func (a *XXTActivity) writeCourseWorkAndExam(setting config.Setting, cache *xuexitongApi.XueXiTUserCache, courseItem *xuexitong.XueXiTCourse) {
	user := a.User
	if user.CoursesCustom.AutoExam == 0 {
		return
	}

	if user.CoursesCustom.AutoExam == 1 {
		if err := aiq.AICheck(setting.AiSetting.AiUrl, setting.AiSetting.Model, setting.AiSetting.APIKEY, setting.AiSetting.AiType); err != nil {
			monitor.GlobalEventBus.AddLog(a.Uid, "⚠️ AI检查失败: "+err.Error())
			return
		}
	} else if user.CoursesCustom.AutoExam == 2 {
		if err := external.CheckApiQueRequest(setting.ApiQueSetting.Url, 5, nil); err != nil {
			monitor.GlobalEventBus.AddLog(a.Uid, "⚠️ 外置题库检查失败: "+err.Error())
			return
		}
	}

	if user.CoursesCustom.CxWorkSw != nil && *user.CoursesCustom.CxWorkSw == 1 {
		workList, err1 := xuexitong.PullWorkListAction(cache, *courseItem)
		if err1 == nil {
			for _, work := range workList {
				if !a.IsRunning() {
					return
				}
				if work.Status != "待做" && work.Status != "未交" && work.Status != "待重做" {
					continue
				}
				err2 := xuexitong.EnterWorkAction(cache, &work)
				if err2 != nil {
					continue
				}
				a.workAction(cache, work, courseItem)
			}
		}
	}

	if user.CoursesCustom.CxExamSw != nil && *user.CoursesCustom.CxExamSw == 1 {
		examList, err1 := xuexitong.PullExamListAction(cache, *courseItem)
		if err1 == nil {
			for _, exam := range examList {
				if !a.IsRunning() {
					return
				}
				if exam.Status != "待做" && exam.Status != "待重考" {
					continue
				}
				err2 := xuexitong.EnterExamAction(cache, &exam)
				if err2 != nil {
					continue
				}
				a.examAction(cache, exam, courseItem)
			}
		}
	}
}

func (a *XXTActivity) chapterTestAction(cache *xuexitongApi.XueXiTUserCache, questionAction xuexitongApi.Question, courseItem *xuexitong.XueXiTCourse, knowledgeItem xuexitong.KnowledgeItem) {
	user := a.User
	setting := a.Setting
	stopStart := 1
	stopEnd := 2

	for i := range questionAction.Choice {
		q := &questionAction.Choice[i]
		switch user.CoursesCustom.AutoExam {
		case 1:
			message := xuexitong.AIProblemMessage(questionAction.Title, q.Type.String(), xuexitongApi.ExamTurn{XueXChoiceQue: *q})
			q.AnswerAIGet(cache.UserID, setting.AiSetting.AiUrl, setting.AiSetting.Model, setting.AiSetting.AiType, message, setting.AiSetting.APIKEY)
		case 2:
			q.AnswerExternalGet(setting.ApiQueSetting.Url)
		case 3:
			message := xuexitong.AIProblemMessage(questionAction.Title, q.Type.String(), xuexitongApi.ExamTurn{XueXChoiceQue: *q})
			q.AnswerXXTAIGet(cache, questionAction.ClassId, questionAction.CourseId, questionAction.Cpi, message)
		}
		time.Sleep(time.Duration(rand.Intn(stopEnd-stopStart)+stopStart) * time.Second)
	}
	for i := range questionAction.Judge {
		q := &questionAction.Judge[i]
		switch user.CoursesCustom.AutoExam {
		case 1:
			message := xuexitong.AIProblemMessage(questionAction.Title, q.Type.String(), xuexitongApi.ExamTurn{XueXJudgeQue: *q})
			q.AnswerAIGet(cache.UserID, setting.AiSetting.AiUrl, setting.AiSetting.Model, setting.AiSetting.AiType, message, setting.AiSetting.APIKEY)
		case 2:
			q.AnswerExternalGet(setting.ApiQueSetting.Url)
		case 3:
			message := xuexitong.AIProblemMessage(questionAction.Title, q.Type.String(), xuexitongApi.ExamTurn{XueXJudgeQue: *q})
			q.AnswerXXTAIGet(cache, questionAction.ClassId, questionAction.CourseId, questionAction.Cpi, message)
		}
		time.Sleep(time.Duration(rand.Intn(stopEnd-stopStart)+stopStart) * time.Second)
	}
	for i := range questionAction.Fill {
		q := &questionAction.Fill[i]
		switch user.CoursesCustom.AutoExam {
		case 1:
			message := xuexitong.AIProblemMessage(questionAction.Title, q.Type.String(), xuexitongApi.ExamTurn{XueXFillQue: *q})
			q.AnswerAIGet(cache.UserID, setting.AiSetting.AiUrl, setting.AiSetting.Model, setting.AiSetting.AiType, message, setting.AiSetting.APIKEY)
		case 2:
			q.AnswerExternalGet(setting.ApiQueSetting.Url)
		case 3:
			message := xuexitong.AIProblemMessage(questionAction.Title, q.Type.String(), xuexitongApi.ExamTurn{XueXFillQue: *q})
			q.AnswerXXTAIGet(cache, questionAction.ClassId, questionAction.CourseId, questionAction.Cpi, message)
		}
		time.Sleep(time.Duration(rand.Intn(stopEnd-stopStart)+stopStart) * time.Second)
	}
	for i := range questionAction.Short {
		q := &questionAction.Short[i]
		switch user.CoursesCustom.AutoExam {
		case 1:
			message := xuexitong.AIProblemMessage(questionAction.Title, q.Type.String(), xuexitongApi.ExamTurn{XueXShortQue: *q})
			q.AnswerAIGet(cache.UserID, setting.AiSetting.AiUrl, setting.AiSetting.Model, setting.AiSetting.AiType, message, setting.AiSetting.APIKEY)
		case 2:
			q.AnswerExternalGet(setting.ApiQueSetting.Url)
		case 3:
			message := xuexitong.AIProblemMessage(questionAction.Title, q.Type.String(), xuexitongApi.ExamTurn{XueXShortQue: *q})
			q.AnswerXXTAIGet(cache, questionAction.ClassId, questionAction.CourseId, questionAction.Cpi, message)
		}
		time.Sleep(time.Duration(rand.Intn(stopEnd-stopStart)+stopStart) * time.Second)
	}
	for i := range questionAction.TermExplanation {
		q := &questionAction.TermExplanation[i]
		switch user.CoursesCustom.AutoExam {
		case 1:
			message := xuexitong.AIProblemMessage(questionAction.Title, q.Type.String(), xuexitongApi.ExamTurn{XueXTermExplanationQue: *q})
			q.AnswerAIGet(cache.UserID, setting.AiSetting.AiUrl, setting.AiSetting.Model, setting.AiSetting.AiType, message, setting.AiSetting.APIKEY)
		case 2:
			q.AnswerExternalGet(setting.ApiQueSetting.Url)
		case 3:
			message := xuexitong.AIProblemMessage(questionAction.Title, q.Type.String(), xuexitongApi.ExamTurn{XueXTermExplanationQue: *q})
			q.AnswerXXTAIGet(cache, questionAction.ClassId, questionAction.CourseId, questionAction.Cpi, message)
		}
		time.Sleep(time.Duration(rand.Intn(stopEnd-stopStart)+stopStart) * time.Second)
	}
	for i := range questionAction.Essay {
		q := &questionAction.Essay[i]
		switch user.CoursesCustom.AutoExam {
		case 1:
			message := xuexitong.AIProblemMessage(questionAction.Title, q.Type.String(), xuexitongApi.ExamTurn{XueXEssayQue: *q})
			q.AnswerAIGet(cache.UserID, setting.AiSetting.AiUrl, setting.AiSetting.Model, setting.AiSetting.AiType, message, setting.AiSetting.APIKEY)
		case 2:
			q.AnswerExternalGet(setting.ApiQueSetting.Url)
		case 3:
			message := xuexitong.AIProblemMessage(questionAction.Title, q.Type.String(), xuexitongApi.ExamTurn{XueXEssayQue: *q})
			q.AnswerXXTAIGet(cache, questionAction.ClassId, questionAction.CourseId, questionAction.Cpi, message)
		}
		time.Sleep(time.Duration(rand.Intn(stopEnd-stopStart)+stopStart) * time.Second)
	}
	for i := range questionAction.Matching {
		q := &questionAction.Matching[i]
		switch user.CoursesCustom.AutoExam {
		case 1:
			message := xuexitong.AIProblemMessage(questionAction.Title, q.Type.String(), xuexitongApi.ExamTurn{XueXMatchingQue: *q})
			q.AnswerAIGet(cache.UserID, setting.AiSetting.AiUrl, setting.AiSetting.Model, setting.AiSetting.AiType, message, setting.AiSetting.APIKEY)
		case 2:
			q.AnswerExternalGet(setting.ApiQueSetting.Url)
		case 3:
			message := xuexitong.AIProblemMessage(questionAction.Title, q.Type.String(), xuexitongApi.ExamTurn{XueXMatchingQue: *q})
			q.AnswerXXTAIGet(cache, questionAction.ClassId, questionAction.CourseId, questionAction.Cpi, message)
		}
		time.Sleep(time.Duration(rand.Intn(stopEnd-stopStart)+stopStart) * time.Second)
	}

	xuexitong.AnswerFixedPattern(questionAction.Choice, questionAction.Judge)
	if user.CoursesCustom.ExamAutoSubmit == 0 {
		xuexitong.WorkNewSubmitAnswerAction(cache, questionAction, false)
	} else if user.CoursesCustom.ExamAutoSubmit == 1 {
		xuexitong.WorkNewSubmitAnswerAction(cache, questionAction, true)
	} else if user.CoursesCustom.ExamAutoSubmit == 2 {
		if checkAnswerIsAvoid(questionAction.Choice, questionAction.Judge, questionAction.Fill, questionAction.Short) {
			xuexitong.WorkNewSubmitAnswerAction(cache, questionAction, false)
		} else {
			xuexitong.WorkNewSubmitAnswerAction(cache, questionAction, true)
		}
	}
	monitor.GlobalEventBus.AddLog(a.Uid, "✅ 章节测试完成: "+questionAction.Title)
}

func (a *XXTActivity) workAction(cache *xuexitongApi.XueXiTUserCache, work xuexitong.XXTWork, courseItem *xuexitong.XueXiTCourse) {
	user := a.User
	setting := a.Setting
	monitor.GlobalEventBus.AddLog(a.Uid, "📝 正在做作业: "+work.Name)
	for i := range work.QuestionTotal {
		if !a.IsRunning() {
			return
		}
		question, err2 := work.PullWorkQuestionAction(cache, i)
		if err2 != nil {
			continue
		}
		if user.CoursesCustom.AutoExam == 1 {
			question.WriteQuestionForAIAction(cache, setting.AiSetting.AiUrl, setting.AiSetting.Model, setting.AiSetting.AiType, setting.AiSetting.APIKEY, setting.VisionAiSetting.AiUrl, setting.VisionAiSetting.Model, setting.VisionAiSetting.AiType, setting.VisionAiSetting.APIKEY)
		} else if user.CoursesCustom.AutoExam == 2 {
			question.WriteQuestionForExternalAction(setting.ApiQueSetting.Url)
		} else if user.CoursesCustom.AutoExam == 3 {
			question.WriteQuestionForXXTAIAction(cache, question.ClassId, question.CourseId, question.Cpi)
		}
		submitResult, err3 := question.SubmitWorkAnswerAction(cache, (user.CoursesCustom.ExamAutoSubmit == 1 || user.CoursesCustom.ExamAutoSubmit == 2) && work.QuestionTotal == i+1)
		if err3 != nil {
			monitor.GlobalEventBus.AddLog(a.Uid, "作业提交失败: "+err3.Error())
		} else {
			monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("作业提交成功 Q%d: %s", i+1, submitResult))
		}
	}
	monitor.GlobalEventBus.AddLog(a.Uid, "✅ 作业完成: "+work.Name)
}

func (a *XXTActivity) examAction(cache *xuexitongApi.XueXiTUserCache, exam xuexitong.XXTExam, courseItem *xuexitong.XueXiTCourse) {
	user := a.User
	setting := a.Setting
	monitor.GlobalEventBus.AddLog(a.Uid, "正在考试: "+exam.Name)
	for i := range exam.QuestionTotal {
		if !a.IsRunning() {
			return
		}
		question, err2 := exam.PullExamQuestionAction(cache, i)
		if err2 != nil {
			continue
		}
		if user.CoursesCustom.AutoExam == 1 {
			question.WriteQuestionForAIAction(cache, setting.AiSetting.AiUrl, setting.AiSetting.Model, setting.AiSetting.AiType, setting.AiSetting.APIKEY, setting.VisionAiSetting.AiUrl, setting.VisionAiSetting.Model, setting.VisionAiSetting.AiType, setting.VisionAiSetting.APIKEY)
		} else if user.CoursesCustom.AutoExam == 2 {
			question.WriteQuestionForExternalAction(setting.ApiQueSetting.Url)
		} else if user.CoursesCustom.AutoExam == 3 {
			question.WriteQuestionForXXTAIAction(cache, question.ClassId, question.CourseId, question.Cpi)
		}
		isSubmit := false
		if (user.CoursesCustom.ExamAutoSubmit == 1 || user.CoursesCustom.ExamAutoSubmit == 2) && exam.QuestionTotal == i+1 {
			isSubmit = true
		}
		submitResult, err3 := question.SubmitExamAnswerAction(cache, isSubmit)
		if err3 != nil {
			monitor.GlobalEventBus.AddLog(a.Uid, "考试提交失败: "+err3.Error())
		} else {
			monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("考试提交成功 Q%d", i+1))
		}
		re := regexp.MustCompile(`考试(\d+)分钟内不允许提交考试`)
		matches := re.FindStringSubmatch(submitResult)
		if len(matches) > 1 {
			minSubmitTime, _ := strconv.Atoi(matches[1])
			monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("检测到该考试限制开考%d分钟内不允许提交，已自动延时%d分钟...", minSubmitTime, minSubmitTime))
			time.Sleep(time.Duration(minSubmitTime) * time.Minute)
			submitResult, err3 = question.SubmitExamAnswerAction(cache, isSubmit)
		}
		if strings.Contains(submitResult, "考试时间已用完") || strings.Contains(submitResult, "time_used") || strings.Contains(submitResult, "expired") {
			monitor.GlobalEventBus.AddLog(a.Uid, "考试时间已用完，跳过: "+exam.Name)
			break
		}
	}
	monitor.GlobalEventBus.AddLog(a.Uid, "✅ 考试完成: "+exam.Name)
}

func checkAnswerIsAvoid(choices []xuexitongApi.ChoiceQue, judges []xuexitongApi.JudgeQue, fills []xuexitongApi.FillQue, shorts []xuexitongApi.ShortQue) bool {
	for _, choice := range choices {
		resStatus := true
		if choice.Answers != nil {
			candidateSelects := []string{}
			for _, option := range choice.Options {
				candidateSelects = append(candidateSelects, option)
			}
			for _, answer := range choice.Answers {
				var sortArray []qutils.Co = qutils.SimilarityArrayAndSort(answer, candidateSelects)
				if sortArray[0].Score >= 0.9 {
					resStatus = false
				}
			}
			if resStatus {
				return true
			}
		} else {
			return true
		}
	}
	for _, judge := range judges {
		resStatus := true
		if judge.Answers != nil {
			for _, answer := range judge.Answers {
				for _, option := range judge.Options {
					if answer == option || answer == "\u9519\u8bef" || answer == "\u6b63\u786e" {
						resStatus = false
					}
				}
			}
			if resStatus {
				return true
			}
		} else {
			return true
		}
	}
	for _, fill := range fills {
		if fill.OpFromAnswer == nil || len(fill.OpFromAnswer) <= 0 {
			return true
		}
	}
	for _, short := range shorts {
		if short.OpFromAnswer == nil || len(short.OpFromAnswer) <= 0 {
			return true
		}
	}
	return false
}

// ==================== YINGHUA ====================

type YingHuaActivity struct {
	UserActivityBase
}

func (a *YingHuaActivity) Login() error {
	cache := &yinghuaApi.YingHuaUserCache{PreUrl: a.User.URL, Account: a.User.Account, Password: a.User.Password}
	if a.User.IsProxy == 1 {
		cache.IpProxySW = true
	}
	err1 := yinghua.YingHuaLoginAction(cache)
	if err1 != nil {
		monitor.GlobalEventBus.SetError(a.Uid, "❌ 登录失败: "+err1.Error())
		return err1
	}
	a.UserCache = cache
	monitor.GlobalEventBus.AddLog(a.Uid, "🔑 登录成功: "+cache.Account)
	return nil
}

func (a *YingHuaActivity) Start() error {
	a.setStopped(false)
	a.setPaused(false)
	a.setRunning(true)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusRunning)
	monitor.GlobalEventBus.AddLog(a.Uid, "🚀 开始刷课任务...")
	go a.run()
	return nil
}

func (a *YingHuaActivity) Stop() error {
	a.setStopped(true)
	a.setRunning(false)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusStopped)
	monitor.GlobalEventBus.AddLog(a.Uid, "⏹️ 任务已停止")
	return nil
}

func (a *YingHuaActivity) Pause() error {
	a.setRunning(false)
	a.setPaused(true)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusPaused)
	monitor.GlobalEventBus.AddLog(a.Uid, "⏸️ 任务已暂停")
	return nil
}

func (a *YingHuaActivity) run() {
	defer func() {
		a.setRunning(false)
		if r := recover(); r != nil {
			monitor.GlobalEventBus.SetError(a.Uid, fmt.Sprintf("Panic: %v", r))
			return
		}
		if !a.isStopped() && !a.isPaused() {
			monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusCompleted)
			a.sendCompletionEmail()
			monitor.GlobalEventBus.UpdateProgress(a.Uid, 100, "\u6240\u6709\u8bfe\u7a0b\u5b66\u4e60\u5b8c\u6bd5")
		}
	}()

	if a.UserCache == nil {
		monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusLogging)
		monitor.GlobalEventBus.AddLog(a.Uid, "🔑 正在登录...")
		if err := a.Login(); err != nil {
			return
		}
	}

	cache := a.UserCache.(*yinghuaApi.YingHuaUserCache)

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if !a.IsRunning() {
				return
			}
			yinghuaApi.KeepAliveApi(*cache, 8)
		}
	}()

	list, _ := yinghua.CourseListAction(cache)
	totalCourses := len(list)
	monitor.GlobalEventBus.UpdateCourseProgress(a.Uid, 0, totalCourses)
	monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("获取到 %d 门课程", totalCourses))

	for i, course := range list {
		if !a.IsRunning() {
			return
		}
		monitor.GlobalEventBus.UpdateProgress(a.Uid, float64(i)/float64(totalCourses)*100, "正在学习: "+course.Name)
		monitor.GlobalEventBus.UpdateCourseProgress(a.Uid, i, totalCourses)
		a.nodeListStudy(cache, &course)
		monitor.GlobalEventBus.UpdateCourseProgress(a.Uid, i+1, totalCourses)
	}
}

func (a *YingHuaActivity) nodeListStudy(cache *yinghuaApi.YingHuaUserCache, course *yinghua.YingHuaCourse) {
	user := a.User
	_ = a.Setting

	if len(user.CoursesCustom.ExcludeCourses) != 0 && config.CmpCourse(course.Name, user.CoursesCustom.ExcludeCourses) {
		monitor.GlobalEventBus.AddLog(a.Uid, "跳过(已排除): "+course.Name)
		return
	}
	if len(user.CoursesCustom.IncludeCourses) != 0 && !config.CmpCourse(course.Name, user.CoursesCustom.IncludeCourses) {
		monitor.GlobalEventBus.AddLog(a.Uid, "跳过(不在包含列表): "+course.Name)
		return
	}
	if time.Now().Before(course.StartDate) {
		monitor.GlobalEventBus.AddLog(a.Uid, "跳过(未开始): "+course.Name)
		return
	}

	monitor.GlobalEventBus.AddLog(a.Uid, "正在学习: "+course.Name)
	nodeList, err := yinghua.VideosListAction(cache, *course)
	if err != nil {
		monitor.GlobalEventBus.AddLog(a.Uid, "拉取视频失败: "+err.Error())
		return
	}

	for _, node := range nodeList {
		if !a.IsRunning() {
			return
		}
		if node.TabVideo && int(node.Progress) != 100 {
			switch user.CoursesCustom.VideoModel {
			case 1:
				a.videoAction(cache, course, node)
			case 2:
				a.videoViolenceAction(cache, course, node)
			case 3:
				a.videoBadRedAction(cache, course, node)
			}
		}
		if user.CoursesCustom.AutoExam != 0 {
			if node.TabWork {
				a.workAction(cache, course, node)
			}
			if node.TabExam {
				a.examAction(cache, course, node)
			}
		}
	}
	monitor.GlobalEventBus.AddLog(a.Uid, "课程完成: "+course.Name)
}

func (a *YingHuaActivity) videoAction(cache *yinghuaApi.YingHuaUserCache, course *yinghua.YingHuaCourse, node yinghua.YingHuaNode) {
	monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频(普通): %s - %s", course.Name, node.Name))
	viewTime := node.ViewedDuration
	studyId := "0"
	for {
		if !a.IsRunning() {
			return
		}
		viewTime += 5
		if node.Progress == 100 {
			break
		}
		sub, err := yinghua.SubmitStudyTimeAction(cache, node.Id, studyId, viewTime)
		if err != nil {
			monitor.GlobalEventBus.AddLog(a.Uid, "提交时间错误: "+err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		yinghua.LoginTimeoutAfreshAction(cache, sub)
		msgVal := gojsonq.New().JSONString(sub).Find("msg")
		msg, ok := msgVal.(string)
		if !ok || msg == "" {
			time.Sleep(10 * time.Second)
			continue
		}
		if msg != "提交学时成功!" {
			time.Sleep(10 * time.Second)
			continue
		}
		studyIdVal := gojsonq.New().JSONString(sub).Find("result.data.studyId")
		if idFloat, ok := studyIdVal.(float64); ok {
			studyId = strconv.Itoa(int(idFloat))
		}
		monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频: %s %d/%d %.1f%%", node.Name, viewTime, node.VideoDuration, float32(viewTime)/float32(node.VideoDuration)*100))
		time.Sleep(5 * time.Second)
		if viewTime >= node.VideoDuration {
			break
		}
	}
}

func (a *YingHuaActivity) videoViolenceAction(cache *yinghuaApi.YingHuaUserCache, course *yinghua.YingHuaCourse, node yinghua.YingHuaNode) {
	monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频(暴力): %s - %s", course.Name, node.Name))
	viewTime := node.ViewedDuration
	studyId := "0"
	for {
		if !a.IsRunning() {
			return
		}
		viewTime += 8
		if node.Progress == 100 {
			break
		}
		sub, err := yinghua.SubmitStudyTimeAction(cache, node.Id, studyId, viewTime)
		if err != nil {
			monitor.GlobalEventBus.AddLog(a.Uid, "提交时间错误: "+err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		yinghua.LoginTimeoutAfreshAction(cache, sub)
		msgVal := gojsonq.New().JSONString(sub).Find("msg")
		msg, ok := msgVal.(string)
		if !ok || msg == "" {
			time.Sleep(10 * time.Second)
			continue
		}
		if msg != "提交学时成功!" {
			time.Sleep(10 * time.Second)
			continue
		}
		studyIdVal := gojsonq.New().JSONString(sub).Find("result.data.studyId")
		if idFloat, ok := studyIdVal.(float64); ok {
			studyId = strconv.Itoa(int(idFloat))
		}
		monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频: %s %d/%d %.1f%%", node.Name, viewTime, node.VideoDuration, float32(viewTime)/float32(node.VideoDuration)*100))
		time.Sleep(8 * time.Second)
		if viewTime >= node.VideoDuration {
			break
		}
	}
}

func (a *YingHuaActivity) videoBadRedAction(cache *yinghuaApi.YingHuaUserCache, course *yinghua.YingHuaCourse, node yinghua.YingHuaNode) {
	monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频(去红): %s - %s", course.Name, node.Name))
	viewTime := node.ViewedDuration
	studyId := "0"
	for {
		if !a.IsRunning() {
			return
		}
		viewTime += 8
		if node.Progress == 100 {
			break
		}
		sub, err := yinghua.SubmitStudyTimeAction(cache, node.Id, studyId, viewTime)
		if err != nil {
			monitor.GlobalEventBus.AddLog(a.Uid, "提交时间错误: "+err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		yinghua.LoginTimeoutAfreshAction(cache, sub)
		msgVal := gojsonq.New().JSONString(sub).Find("msg")
		msg, ok := msgVal.(string)
		if !ok || msg == "" {
			time.Sleep(10 * time.Second)
			continue
		}
		if msg != "提交学时成功!" {
			if strings.Contains(msg, "检测到可能使用并行播放刷课") {
				monitor.GlobalEventBus.AddLog(a.Uid, "🔴 检测到标红，尝试去红: "+node.Name)
			}
			time.Sleep(10 * time.Second)
			continue
		}
		studyIdVal := gojsonq.New().JSONString(sub).Find("result.data.studyId")
		if idFloat, ok := studyIdVal.(float64); ok {
			studyId = strconv.Itoa(int(idFloat))
		}
		monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频: %s %d/%d %.1f%%", node.Name, viewTime, node.VideoDuration, float32(viewTime)/float32(node.VideoDuration)*100))
		time.Sleep(8 * time.Second)
		if viewTime >= node.VideoDuration {
			break
		}
	}
}

func (a *YingHuaActivity) workAction(cache *yinghuaApi.YingHuaUserCache, course *yinghua.YingHuaCourse, node yinghua.YingHuaNode) {
	user := a.User
	setting := a.Setting
	detailAction, _ := yinghua.WorkDetailAction(cache, node.Id)
	if len(detailAction) == 0 {
		return
	}
	monitor.GlobalEventBus.AddLog(a.Uid, "📝 正在做作业: "+node.Name)
	for _, work := range detailAction {
		var err error
		switch user.CoursesCustom.AutoExam {
		case 1:
			err = yinghua.StartWorkAction(cache, work, setting.AiSetting.AiUrl, setting.AiSetting.Model, setting.AiSetting.APIKEY, setting.AiSetting.AiType, user.CoursesCustom.ExamAutoSubmit, setting.VisionAiSetting.AiUrl, setting.VisionAiSetting.Model, setting.VisionAiSetting.AiType, setting.VisionAiSetting.APIKEY)
		case 2:
			err = yinghua.StartWorkForExternalAction(cache, setting.ApiQueSetting.Url, work, user.CoursesCustom.ExamAutoSubmit)
		}
		if err != nil {
			monitor.GlobalEventBus.AddLog(a.Uid, "作业错误: "+err.Error())
		}
	}
}

func (a *YingHuaActivity) examAction(cache *yinghuaApi.YingHuaUserCache, course *yinghua.YingHuaCourse, node yinghua.YingHuaNode) {
	user := a.User
	setting := a.Setting
	detailAction, _ := yinghua.ExamDetailAction(cache, node.Id)
	if len(detailAction) == 0 {
		return
	}
	monitor.GlobalEventBus.AddLog(a.Uid, "正在考试: "+node.Name)
	for _, exam := range detailAction {
		var err error
		switch user.CoursesCustom.AutoExam {
		case 1:
			err = yinghua.StartExamAction(cache, exam, setting.AiSetting.AiUrl, setting.AiSetting.Model, setting.AiSetting.APIKEY, setting.AiSetting.AiType, user.CoursesCustom.ExamAutoSubmit, setting.VisionAiSetting.AiUrl, setting.VisionAiSetting.Model, setting.VisionAiSetting.AiType, setting.VisionAiSetting.APIKEY)
		case 2:
			err = yinghua.StartExamForExternalAction(cache, exam, setting.ApiQueSetting.Url, user.CoursesCustom.ExamAutoSubmit)
		}
		if err != nil {
			monitor.GlobalEventBus.AddLog(a.Uid, "考试错误: "+err.Error())
		}
	}
}

// ==================== ENAEA ====================

type EnaeaActivity struct{ UserActivityBase }

func (a *EnaeaActivity) Login() error {
	cache := &enaeaApi.EnaeaUserCache{Account: a.User.Account, Password: a.User.Password}
	_, err := enaea.EnaeaLoginAction(cache)
	if err != nil {
		monitor.GlobalEventBus.SetError(a.Uid, "❌ 登录失败: "+err.Error())
		return err
	}
	a.UserCache = cache
	monitor.GlobalEventBus.AddLog(a.Uid, "🔑 登录成功: "+a.User.Account)
	return nil
}

func (a *EnaeaActivity) Start() error {
	a.setStopped(false)
	a.setPaused(false)
	a.setRunning(true)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusRunning)
	monitor.GlobalEventBus.AddLog(a.Uid, "🚀 开始刷课任务...")
	go a.run()
	return nil
}

func (a *EnaeaActivity) Stop() error {
	a.setStopped(true)
	a.setRunning(false)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusStopped)
	return nil
}

func (a *EnaeaActivity) Pause() error {
	a.setRunning(false)
	a.setPaused(true)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusPaused)
	return nil
}

func (a *EnaeaActivity) run() {
	defer func() {
		a.setRunning(false)
		if r := recover(); r != nil {
			monitor.GlobalEventBus.SetError(a.Uid, fmt.Sprintf("Panic: %v", r))
			return
		}
		if a.isPaused() {
			return
		}
		if !a.isStopped() {
			monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusCompleted)
			a.sendCompletionEmail()
			monitor.GlobalEventBus.UpdateProgress(a.Uid, 100, "\u6240\u6709\u8bfe\u7a0b\u5b66\u4e60\u5b8c\u6bd5")
		}
	}()
	if a.UserCache == nil {
		monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusLogging)
		monitor.GlobalEventBus.AddLog(a.Uid, "🔑 正在登录...")
		if err := a.Login(); err != nil {
			return
		}
	}
	user := a.User
	cache := a.UserCache.(*enaeaApi.EnaeaUserCache)
	projectList, err := enaea.ProjectListAction(cache)
	if err != nil {
		monitor.GlobalEventBus.SetError(a.Uid, "拉取项目列表失败: "+err.Error())
		return
	}
	totalProjects := len(projectList)
	monitor.GlobalEventBus.UpdateCourseProgress(a.Uid, 0, totalProjects)
	monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("获取到 %d 个项目", totalProjects))

	excludeProjects := []string{}
	includeProjects := []string{}
	excludeTitleTags := []string{}
	includeTitleTags := []string{}
	for _, cours := range user.CoursesCustom.ExcludeCourses {
		split := strings.Split(cours, "-->")
		if len(split) >= 1 {
			excludeProjects = append(excludeProjects, split[0])
		}
		if len(split) >= 2 {
			excludeTitleTags = append(excludeTitleTags, split[1])
		}
	}
	for _, cours := range user.CoursesCustom.IncludeCourses {
		split := strings.Split(cours, "-->")
		if len(split) >= 1 {
			includeProjects = append(includeProjects, split[0])
		}
		if len(split) >= 2 {
			includeTitleTags = append(includeTitleTags, split[1])
		}
	}

	for i, project := range projectList {
		if !a.IsRunning() {
			return
		}
		if len(excludeProjects) != 0 && config.CmpCourse(project.ClusterName, excludeProjects) {
			monitor.GlobalEventBus.AddLog(a.Uid, "跳过项目(已排除): "+project.ClusterName)
			continue
		}
		if len(includeProjects) != 0 && !config.CmpCourse(project.ClusterName, includeProjects) {
			monitor.GlobalEventBus.AddLog(a.Uid, "跳过项目(不在包含列表): "+project.ClusterName)
			continue
		}
		monitor.GlobalEventBus.UpdateProgress(a.Uid, float64(i)/float64(totalProjects)*100, "正在学习: "+project.CircleName)
		monitor.GlobalEventBus.UpdateCourseProgress(a.Uid, i, totalProjects)
		courseList, err2 := enaea.CourseListAction(cache, project.CircleId)
		if err2 != nil {
			continue
		}
		for _, course := range courseList {
			if !a.IsRunning() {
				return
			}
			if len(excludeTitleTags) != 0 && config.CmpCourse(course.TitleTag, excludeTitleTags) {
				monitor.GlobalEventBus.AddLog(a.Uid, "跳过课程(已排除): "+course.TitleTag)
				continue
			}
			if len(includeTitleTags) != 0 && !config.CmpCourse(course.TitleTag, includeTitleTags) {
				monitor.GlobalEventBus.AddLog(a.Uid, "跳过课程(不在包含列表): "+course.TitleTag)
				continue
			}
			a.nodeListStudy(cache, &course)
		}
		monitor.GlobalEventBus.AddLog(a.Uid, "项目完成: "+project.CircleName)
	}
	monitor.GlobalEventBus.AddLog(a.Uid, "所有待学习课程学习完毕")
}

func (a *EnaeaActivity) nodeListStudy(cache *enaeaApi.EnaeaUserCache, course *enaea.EnaeaCourse) {
	monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("正在学习: 【%s】【%s】", course.TitleTag, course.CourseTitle))
	videoList, err := enaea.VideoListAction(cache, course)
	for err != nil {
		enaea.LoginTimeoutAfreshAction(cache, err)
		videoList, err = enaea.VideoListAction(cache, course)
	}
	for _, video := range videoList {
		if !a.IsRunning() {
			return
		}
		a.videoAction(cache, &video)
	}
	monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("课程完成: 【%s】【%s】", course.TitleTag, course.CourseTitle))
}

func (a *EnaeaActivity) videoAction(cache *enaeaApi.EnaeaUserCache, node *enaea.EnaeaVideo) {
	if a.User.CoursesCustom.VideoModel == 0 {
		return
	}
	monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频: 【%s】【%s】", node.CourseName, node.CourseContentStr))
	err := enaea.StatisticTicForCCVideAction(cache, node)
	if err != nil {
		monitor.GlobalEventBus.AddLog(a.Uid, "预提交异常: "+err.Error())
	}
	for {
		if !a.IsRunning() {
			return
		}
		if node.StudyProgress >= 100 {
			monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频完成: 【%s】【%s】", node.CourseName, node.CourseContentStr))
			break
		}
		var submitErr error
		if a.User.CoursesCustom.VideoModel == 1 {
			submitErr = enaea.SubmitStudyTimeAction(cache, node, time.Now().UnixMilli(), 0)
		} else if a.User.CoursesCustom.VideoModel == 2 {
			submitErr = enaea.SubmitStudyTimeAction(cache, node, 60, 1)
		}
		if submitErr != nil {
			if submitErr.Error() != "request frequently!" {
				monitor.GlobalEventBus.AddLog(a.Uid, "提交学时异常: "+submitErr.Error())
			}
		}
		enaea.LoginTimeoutAfreshAction(cache, submitErr)
		monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频: 【%s】【%s】进度: %.2f%%", node.CourseName, node.CourseContentStr, node.StudyProgress))
		time.Sleep(25 * time.Second)
		if node.StudyProgress >= 100 {
			break
		}
	}
}

// ==================== CQIE ====================

type CqieActivity struct{ UserActivityBase }

func (a *CqieActivity) Login() error {
	cache := &cqieApi.CqieUserCache{Account: a.User.Account, Password: a.User.Password}
	err := cqie.CqieLoginAction(cache)
	if err != nil {
		monitor.GlobalEventBus.SetError(a.Uid, "❌ 登录失败: "+err.Error())
		return err
	}
	a.UserCache = cache
	monitor.GlobalEventBus.AddLog(a.Uid, "🔑 登录成功: "+a.User.Account)
	return nil
}

func (a *CqieActivity) Start() error {
	a.setStopped(false)
	a.setPaused(false)
	a.setRunning(true)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusRunning)
	go a.run()
	return nil
}

func (a *CqieActivity) Stop() error {
	a.setStopped(true)
	a.setRunning(false)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusStopped)
	return nil
}

func (a *CqieActivity) Pause() error {
	a.setRunning(false)
	a.setPaused(true)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusPaused)
	return nil
}

func (a *CqieActivity) run() {
	defer func() {
		a.setRunning(false)
		if r := recover(); r != nil {
			monitor.GlobalEventBus.SetError(a.Uid, fmt.Sprintf("Panic: %v", r))
			return
		}
		if a.isPaused() {
			return
		}
		if !a.isStopped() {
			monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusCompleted)
			a.sendCompletionEmail()
			monitor.GlobalEventBus.UpdateProgress(a.Uid, 100, "\u6240\u6709\u8bfe\u7a0b\u5b66\u4e60\u5b8c\u6bd5")
		}
	}()

	if a.UserCache == nil {
		monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusLogging)
		monitor.GlobalEventBus.AddLog(a.Uid, "🔑 正在登录...")
		if err := a.Login(); err != nil {
			return
		}
	}
	user := a.User
	cache := a.UserCache.(*cqieApi.CqieUserCache)
	courseList, err := cqie.CqiePullCourseListAction(cache)
	if err != nil {
		monitor.GlobalEventBus.SetError(a.Uid, "❌ 拉取课程列表失败: "+err.Error())
		return
	}
	totalCourses := len(courseList)
	monitor.GlobalEventBus.UpdateCourseProgress(a.Uid, 0, totalCourses)
	monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("获取到 %d 门课程", totalCourses))
	for i, course := range courseList {
		if !a.IsRunning() {
			return
		}
		if len(user.CoursesCustom.ExcludeCourses) != 0 && config.CmpCourse(course.CourseName, user.CoursesCustom.ExcludeCourses) {
			monitor.GlobalEventBus.AddLog(a.Uid, "跳过(已排除): "+course.CourseName)
			continue
		}
		if len(user.CoursesCustom.IncludeCourses) != 0 && !config.CmpCourse(course.CourseName, user.CoursesCustom.IncludeCourses) {
			monitor.GlobalEventBus.AddLog(a.Uid, "跳过(不在包含列表): "+course.CourseName)
			continue
		}
		monitor.GlobalEventBus.UpdateProgress(a.Uid, float64(i)/float64(totalCourses)*100, "正在学习: "+course.CourseName)
		monitor.GlobalEventBus.UpdateCourseProgress(a.Uid, i, totalCourses)
		monitor.GlobalEventBus.AddLog(a.Uid, "正在学习: "+course.CourseName)
		videoList, err2 := cqie.PullCourseVideoListAndProgress(cache, &course)
		if err2 != nil {
			continue
		}
		for _, video := range videoList {
			if !a.IsRunning() {
				return
			}
			switch user.CoursesCustom.VideoModel {
			case 1:
				a.videoAction(cache, &video)
			case 2:
				a.videoSpeedAction(cache, &video)
			}
		}
		monitor.GlobalEventBus.AddLog(a.Uid, "课程完成: "+course.CourseName)
	}
	monitor.GlobalEventBus.AddLog(a.Uid, "所有待学习课程学习完毕")
}

func (a *CqieActivity) videoAction(cache *cqieApi.CqieUserCache, video *cqie.CqieVideo) {
	monitor.GlobalEventBus.AddLog(a.Uid, "视频(常规): "+video.VideoName)
	cqie.StartStudyVideoAction(cache, video)
	startPos := video.StudyTime
	stopPos := video.StudyTime
	maxPos := video.StudyTime
	err := cqie.SaveVideoStudyTimeAction(cache, video, startPos, stopPos)
	if err != nil {
		monitor.GlobalEventBus.AddLog(a.Uid, "保存学习点异常: "+err.Error())
	}
	for {
		if !a.IsRunning() {
			return
		}
		if maxPos >= video.TimeLength+3 {
			startPos = video.TimeLength
			stopPos = video.TimeLength
			maxPos = video.TimeLength
			break
		}
		if stopPos >= maxPos {
			maxPos = startPos + 3
		}
		err := cqie.SubmitStudyTimeAction(cache, video, time.Now(), startPos, stopPos, maxPos)
		if err != nil {
			monitor.GlobalEventBus.AddLog(a.Uid, "提交学时异常: "+err.Error())
		}
		monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频: %s 进度: %.2f%%", video.VideoName, float32(video.StudyTime)/float32(video.TimeLength)*100))
		startPos = startPos + 3
		stopPos = stopPos + 3
		time.Sleep(3 * time.Second)
	}
	err = cqie.SaveVideoStudyTimeAction(cache, video, startPos, stopPos)
	if err != nil {
		monitor.GlobalEventBus.AddLog(a.Uid, "保存学习点异常: "+err.Error())
	}
	monitor.GlobalEventBus.AddLog(a.Uid, "视频完成: "+video.VideoName)
}

func (a *CqieActivity) videoSpeedAction(cache *cqieApi.CqieUserCache, video *cqie.CqieVideo) {
	monitor.GlobalEventBus.AddLog(a.Uid, "视频(秒刷): "+video.VideoName)
	cqie.StartStudyVideoAction(cache, video)
	startPos := video.StudyTime
	stopPos := video.StudyTime
	maxPos := video.StudyTime
	err := cqie.SaveVideoStudyTimeAction(cache, video, startPos, stopPos)
	if err != nil {
		monitor.GlobalEventBus.AddLog(a.Uid, "保存学习点异常: "+err.Error())
	}
	err1 := cqie.SubmitStudyTimeAction(cache, video, time.Now(), startPos, stopPos, maxPos)
	if err1 != nil {
		monitor.GlobalEventBus.AddLog(a.Uid, "提交学时异常: "+err1.Error())
	}
	err = cqie.SaveVideoStudyTimeAction(cache, video, startPos, stopPos)
	if err != nil {
		monitor.GlobalEventBus.AddLog(a.Uid, "保存学习点异常: "+err.Error())
	}
	monitor.GlobalEventBus.AddLog(a.Uid, "视频完成(秒刷): "+video.VideoName)
}

// ==================== ICVE ====================

type IcveActivity struct{ UserActivityBase }

func (a *IcveActivity) Login() error {
	cache := &icveApi.IcveUserCache{Account: a.User.Account, Password: a.User.Password}
	err := icve.IcveLoginAction(cache)
	if err != nil {
		monitor.GlobalEventBus.SetError(a.Uid, "❌ 登录失败: "+err.Error())
		return err
	}
	a.UserCache = cache
	monitor.GlobalEventBus.AddLog(a.Uid, "🔑 登录成功: "+a.User.Account)
	return nil
}

func (a *IcveActivity) Start() error {
	a.setStopped(false)
	a.setPaused(false)
	a.setRunning(true)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusRunning)
	go a.run()
	return nil
}

func (a *IcveActivity) Stop() error {
	a.setStopped(true)
	a.setRunning(false)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusStopped)
	return nil
}

func (a *IcveActivity) Pause() error {
	a.setRunning(false)
	a.setPaused(true)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusPaused)
	return nil
}

func (a *IcveActivity) run() {
	defer func() {
		a.setRunning(false)
		if r := recover(); r != nil {
			monitor.GlobalEventBus.SetError(a.Uid, fmt.Sprintf("Panic: %v", r))
			return
		}
		if !a.isStopped() && !a.isPaused() {
			monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusCompleted)
			a.sendCompletionEmail()
			monitor.GlobalEventBus.UpdateProgress(a.Uid, 100, "\u6240\u6709\u8bfe\u7a0b\u5b66\u4e60\u5b8c\u6bd5")
		}
	}()
	if a.UserCache == nil {
		monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusLogging)
		monitor.GlobalEventBus.AddLog(a.Uid, "🔑 正在登录...")
		if err := a.Login(); err != nil {
			return
		}
	}
	user := a.User
	cache := a.UserCache.(*icveApi.IcveUserCache)
	courseList, err := icve.PullZYKCourseAction(cache)
	if err != nil {
		monitor.GlobalEventBus.SetError(a.Uid, "❌ 拉取课程列表失败: "+err.Error())
		return
	}
	totalCourses := len(courseList)
	monitor.GlobalEventBus.UpdateCourseProgress(a.Uid, 0, totalCourses)
	monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("获取到 %d 门课程", totalCourses))
	for i, course := range courseList {
		if !a.IsRunning() {
			return
		}
		if len(user.CoursesCustom.ExcludeCourses) != 0 && config.CmpCourse(course.CourseName, user.CoursesCustom.ExcludeCourses) {
			monitor.GlobalEventBus.AddLog(a.Uid, "跳过(已排除): "+course.CourseName)
			continue
		}
		if len(user.CoursesCustom.IncludeCourses) != 0 && !config.CmpCourse(course.CourseName, user.CoursesCustom.IncludeCourses) {
			monitor.GlobalEventBus.AddLog(a.Uid, "跳过(不在包含列表): "+course.CourseName)
			continue
		}
		if course.Status == "3" {
			monitor.GlobalEventBus.AddLog(a.Uid, "跳过(已结束): "+course.CourseName)
			continue
		}
		monitor.GlobalEventBus.UpdateProgress(a.Uid, float64(i)/float64(totalCourses)*100, "正在学习: "+course.CourseName)
		monitor.GlobalEventBus.UpdateCourseProgress(a.Uid, i, totalCourses)
		monitor.GlobalEventBus.AddLog(a.Uid, "正在学习: "+course.CourseName)
		nodeList, err2 := icve.PullZYKCourseNodeAction(cache, course)
		if err2 != nil {
			continue
		}
		for _, node := range nodeList {
			if !a.IsRunning() {
				return
			}
			if user.CoursesCustom.VideoModel != 0 && !node.IsLook {
				submitResult, err3 := icve.SubmitZYKStudyTimeAction(cache, node)
				if err3 != nil {
					monitor.GlobalEventBus.AddLog(a.Uid, "学习异常: "+node.Name+" - "+err3.Error())
				} else {
					monitor.GlobalEventBus.AddLog(a.Uid, "学习完毕: "+node.Name+" - "+submitResult)
				}
			}
		}
		monitor.GlobalEventBus.AddLog(a.Uid, "课程完成: "+course.CourseName)
	}
	monitor.GlobalEventBus.AddLog(a.Uid, "所有待学习课程学习完毕")
}

// ==================== QSXT ====================

type QsxtActivity struct{ UserActivityBase }

func (a *QsxtActivity) Login() error {
	cache := &qsxtApi.QsxtUserCache{Account: a.User.Account, Password: a.User.Password}
	_, err := qingshuxuetang.QsxtLoginAction(cache)
	if err != nil {
		monitor.GlobalEventBus.SetError(a.Uid, "❌ 登录失败: "+err.Error())
		return err
	}
	a.UserCache = cache
	monitor.GlobalEventBus.AddLog(a.Uid, "🔑 登录成功: "+a.User.Account)
	return nil
}

func (a *QsxtActivity) Start() error {
	a.setStopped(false)
	a.setPaused(false)
	a.setRunning(true)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusRunning)
	go a.run()
	return nil
}

func (a *QsxtActivity) Stop() error {
	a.setStopped(true)
	a.setRunning(false)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusStopped)
	return nil
}

func (a *QsxtActivity) Pause() error {
	a.setRunning(false)
	a.setPaused(true)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusPaused)
	return nil
}

func (a *QsxtActivity) run() {
	defer func() {
		a.setRunning(false)
		if r := recover(); r != nil {
			monitor.GlobalEventBus.SetError(a.Uid, fmt.Sprintf("Panic: %v", r))
			return
		}
		if !a.isStopped() && !a.isPaused() {
			monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusCompleted)
			a.sendCompletionEmail()
			monitor.GlobalEventBus.UpdateProgress(a.Uid, 100, "\u6240\u6709\u8bfe\u7a0b\u5b66\u4e60\u5b8c\u6bd5")
		}
	}()
	if a.UserCache == nil {
		monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusLogging)
		monitor.GlobalEventBus.AddLog(a.Uid, "🔑 正在登录...")
		if err := a.Login(); err != nil {
			return
		}
	}
	user := a.User
	cache := a.UserCache.(*qsxtApi.QsxtUserCache)
	courseList, err := qingshuxuetang.PullCourseListAction(cache)
	if err != nil {
		monitor.GlobalEventBus.SetError(a.Uid, "❌ 拉取课程列表失败: "+err.Error())
		return
	}
	totalCourses := len(courseList)
	monitor.GlobalEventBus.UpdateCourseProgress(a.Uid, 0, totalCourses)
	monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("获取到 %d 门课程", totalCourses))
	for i, course := range courseList {
		if !a.IsRunning() {
			return
		}
		if len(user.CoursesCustom.ExcludeCourses) != 0 && config.CmpCourse(course.CourseName, user.CoursesCustom.ExcludeCourses) {
			monitor.GlobalEventBus.AddLog(a.Uid, "跳过(已排除): "+course.CourseName)
			continue
		}
		if len(user.CoursesCustom.IncludeCourses) != 0 && !config.CmpCourse(course.CourseName, user.CoursesCustom.IncludeCourses) {
			monitor.GlobalEventBus.AddLog(a.Uid, "跳过(不在包含列表): "+course.CourseName)
			continue
		}
		if course.StudyStatusName != "在修" {
			monitor.GlobalEventBus.AddLog(a.Uid, "跳过(非在修): "+course.CourseName)
			continue
		}
		monitor.GlobalEventBus.UpdateProgress(a.Uid, float64(i)/float64(totalCourses)*100, "正在学习: "+course.CourseName)
		monitor.GlobalEventBus.UpdateCourseProgress(a.Uid, i, totalCourses)
		monitor.GlobalEventBus.AddLog(a.Uid, "正在学习: "+course.CourseName)

		if user.CoursesCustom.VideoModel != 0 && course.CoursewareLearnGainScore < course.CoursewareLearnTotalScore {
			nodeList, err2 := qingshuxuetang.PullCourseNodeListAction(cache, course)
			if err2 == nil {
				for _, node := range nodeList {
					if !a.IsRunning() {
						return
					}
					if course.CoursewareLearnGainScore >= course.CoursewareLearnTotalScore {
						break
					}
					if user.CoursesCustom.VideoModel == 1 {
						a.nodeSubmitTimeAction(cache, &course, node)
					}
				}
			}
		}

		if user.CoursesCustom.VideoModel != 0 && course.CourseMaterialsLearnGainCore < course.CourseMaterialsLearnTotalCore {
			materialList, err2 := qingshuxuetang.PullCourseMaterialListAction(cache, course)
			if err2 == nil {
				for _, material := range materialList {
					if !a.IsRunning() {
						return
					}
					if course.CourseMaterialsLearnGainCore >= course.CourseMaterialsLearnTotalCore {
						break
					}
					if user.CoursesCustom.VideoModel == 1 {
						a.materialSubmitTimeAction(cache, &course, material)
					}
				}
			}
		}

		if user.CoursesCustom.AutoExam != 0 && course.CourseWorkGainScore < course.CourseWorkTotalScore {
			workList, err2 := qingshuxuetang.PullWorkListAction(cache, course)
			if err2 == nil {
				for _, work := range workList {
					if !a.IsRunning() {
						return
					}
					a.workAction(cache, &course, &work)
				}
			}
		}
		monitor.GlobalEventBus.AddLog(a.Uid, "课程完成: "+course.CourseName)
	}
	monitor.GlobalEventBus.AddLog(a.Uid, "所有待学习课程学习完毕")
}

func (a *QsxtActivity) nodeSubmitTimeAction(cache *qsxtApi.QsxtUserCache, course *qingshuxuetang.QsxtCourse, node qingshuxuetang.QsxtNode) {
	if node.NodeType == "chapter" {
		return
	}
	endTime := int(node.Duration / 1000)
	if endTime == 0 {
		endTime = 350
	}
	totalTime := node.TotalStudyDuration
	if totalTime > endTime {
		return
	}
	startId, err2 := node.StartStudyTimeAction(cache)
	if err2 != nil {
		monitor.GlobalEventBus.AddLog(a.Uid, "开始学习异常: "+node.NodeName+" - "+err2.Error())
		return
	}
	studyTime := 0
	maxTime := endTime - totalTime
	monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频: %s 进度: %d/%d", node.NodeName, studyTime+totalTime, endTime))
	for {
		if !a.IsRunning() {
			return
		}
		time.Sleep(60 * time.Second)
		_, err3 := node.SubmitStudyTimeAction(cache, startId, false)
		if err3 != nil {
			monitor.GlobalEventBus.AddLog(a.Uid, "学习异常: "+node.NodeName+" - "+err3.Error())
		} else {
			monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("学时提交成功: %s 进度: %d/%d", node.NodeName, studyTime+totalTime, endTime))
		}
		studyTime += 60
		if studyTime >= maxTime {
			break
		}
	}
	node.SubmitStudyTimeAction(cache, startId, true)
	qingshuxuetang.UpdateCourseScore(cache, course)
	monitor.GlobalEventBus.AddLog(a.Uid, "视频完成: "+node.NodeName)
}

func (a *QsxtActivity) materialSubmitTimeAction(cache *qsxtApi.QsxtUserCache, course *qingshuxuetang.QsxtCourse, material qingshuxuetang.QsxtCourseMaterial) {
	monitor.GlobalEventBus.AddLog(a.Uid, "资料学习: "+material.Name)
	time.Sleep(60 * time.Second)
	qingshuxuetang.UpdateCourseScore(cache, course)
}

func (a *QsxtActivity) workAction(cache *qsxtApi.QsxtUserCache, course *qingshuxuetang.QsxtCourse, work *qingshuxuetang.QsxtWork) {
	monitor.GlobalEventBus.AddLog(a.Uid, "📝 正在做作业: "+work.Title)
	isSubmit := a.User.CoursesCustom.ExamAutoSubmit == 1
	qingshuxuetang.WriteWorkAction(cache, *work, isSubmit)
	monitor.GlobalEventBus.AddLog(a.Uid, "✅ 作业完成: "+work.Title)
}

// ==================== WELEARN ====================

type WeLearnActivity struct{ UserActivityBase }

func (a *WeLearnActivity) Login() error {
	cache := &welearnApi.WeLearnUserCache{Account: a.User.Account, Password: a.User.Password}
	err := welearn.WeLearnLoginAction(cache)
	if err != nil {
		monitor.GlobalEventBus.SetError(a.Uid, "❌ 登录失败: "+err.Error())
		return err
	}
	a.UserCache = cache
	monitor.GlobalEventBus.AddLog(a.Uid, "🔑 登录成功: "+a.User.Account)
	return nil
}

func (a *WeLearnActivity) Start() error {
	a.setStopped(false)
	a.setPaused(false)
	a.setRunning(true)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusRunning)
	go a.run()
	return nil
}

func (a *WeLearnActivity) Stop() error {
	a.setStopped(true)
	a.setRunning(false)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusStopped)
	return nil
}

func (a *WeLearnActivity) Pause() error {
	a.setRunning(false)
	a.setPaused(true)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusPaused)
	return nil
}

func (a *WeLearnActivity) run() {
	defer func() {
		a.setRunning(false)
		if r := recover(); r != nil {
			monitor.GlobalEventBus.SetError(a.Uid, fmt.Sprintf("Panic: %v", r))
			return
		}
		if !a.isStopped() && !a.isPaused() {
			monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusCompleted)
			a.sendCompletionEmail()
			monitor.GlobalEventBus.UpdateProgress(a.Uid, 100, "\u6240\u6709\u8bfe\u7a0b\u5b66\u4e60\u5b8c\u6bd5")
		}
	}()
	if a.UserCache == nil {
		monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusLogging)
		monitor.GlobalEventBus.AddLog(a.Uid, "🔑 正在登录...")
		if err := a.Login(); err != nil {
			return
		}
	}
	user := a.User
	cache := a.UserCache.(*welearnApi.WeLearnUserCache)
	courseList, err := welearn.WeLearnPullCourseListAction(cache)
	if err != nil {
		monitor.GlobalEventBus.SetError(a.Uid, "❌ 拉取课程列表失败: "+err.Error())
		return
	}
	totalCourses := len(courseList)
	monitor.GlobalEventBus.UpdateCourseProgress(a.Uid, 0, totalCourses)
	monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("获取到 %d 门课程", totalCourses))
	for i, course := range courseList {
		if !a.IsRunning() {
			return
		}
		if len(user.CoursesCustom.ExcludeCourses) != 0 && config.CmpCourse(course.Name, user.CoursesCustom.ExcludeCourses) {
			monitor.GlobalEventBus.AddLog(a.Uid, "跳过(已排除): "+course.Name)
			continue
		}
		if len(user.CoursesCustom.IncludeCourses) != 0 && !config.CmpCourse(course.Name, user.CoursesCustom.IncludeCourses) {
			monitor.GlobalEventBus.AddLog(a.Uid, "跳过(不在包含列表): "+course.Name)
			continue
		}
		monitor.GlobalEventBus.UpdateProgress(a.Uid, float64(i)/float64(totalCourses)*100, "正在学习: "+course.Name)
		monitor.GlobalEventBus.UpdateCourseProgress(a.Uid, i, totalCourses)
		monitor.GlobalEventBus.AddLog(a.Uid, "正在学习: "+course.Name)
		chapters, err2 := welearn.WeLearnPullCourseChapterAction(cache, course)
		if err2 != nil {
			continue
		}
		for _, chapter := range chapters {
			if !a.IsRunning() {
				return
			}
			points, err3 := welearn.WeLearnPullChapterPointAction(cache, course, chapter)
			if err3 != nil {
				continue
			}
			for _, pt := range points {
				if !a.IsRunning() {
					return
				}
				switch user.CoursesCustom.VideoModel {
				case 1:
					a.nodeSubmitTimeAction(cache, course, pt)
				case 2:
					a.nodeCompleteAction(cache, course, pt)
				}
			}
		}
		monitor.GlobalEventBus.AddLog(a.Uid, "课程完成: "+course.Name)
	}
	monitor.GlobalEventBus.AddLog(a.Uid, "所有待学习课程学习完毕")
}

func (a *WeLearnActivity) nodeSubmitTimeAction(cache *welearnApi.WeLearnUserCache, course welearn.WeLearnCourse, node welearn.WeLearnPoint) {
	user := a.User
	if user.CoursesCustom.VideoModel == 0 {
		return
	}
	if !node.IsVisible {
		monitor.GlobalEventBus.AddLog(a.Uid, "跳过(未解锁): "+node.Location)
		return
	}
	_, progressMeasure, sessionTime, totalTime, scaled, err := welearn.WeLearnSubmitStudyTimeAction(cache, course, node)
	if err != nil {
		monitor.GlobalEventBus.AddLog(a.Uid, "提交学时异常: "+err.Error())
		return
	}
	endTime := 1600
	learnTime := user.CoursesCustom.StudyTime
	if learnTime != "" {
		parts := strings.Split(learnTime, "-")
		if len(parts) == 2 {
			minVal, err1 := strconv.Atoi(parts[0])
			maxVal, err2 := strconv.Atoi(parts[1])
			if err1 == nil && err2 == nil && maxVal > minVal {
				endTime = (rand.Intn(maxVal-minVal+1) + minVal) * 60
			}
		}
	}
	if totalTime > endTime {
		return
	}
	for {
		if !a.IsRunning() {
			return
		}
		_, err1 := cache.KeepPointSessionPlan1Api(course.Cid, node.Id, course.Uid, course.ClassId, sessionTime, totalTime, 3, nil)
		if err1 != nil {
			monitor.GlobalEventBus.AddLog(a.Uid, "学时提交异常: "+err1.Error())
		} else {
			monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("学时提交成功: %s 进度: %d/%d", node.Location, totalTime, endTime))
		}
		if sessionTime >= endTime {
			break
		}
		sessionTime += 60
		totalTime += 60
		time.Sleep(60 * time.Second)
	}
	cache.SubmitStudyPlan2Api(course.Cid, node.Id, course.Uid, scaled, course.ClassId, progressMeasure, "completed", 3, nil)
	monitor.GlobalEventBus.AddLog(a.Uid, "✅ 任务点完成: "+node.Location)
}

func (a *WeLearnActivity) nodeCompleteAction(cache *welearnApi.WeLearnUserCache, course welearn.WeLearnCourse, node welearn.WeLearnPoint) {
	user := a.User
	if user.CoursesCustom.VideoModel == 0 {
		return
	}
	if !node.IsVisible {
		monitor.GlobalEventBus.AddLog(a.Uid, "跳过(未解锁): "+node.Location)
		return
	}
	if node.IsComplete == "completed" || node.IsComplete == "已完成" {
		return
	}
	err := welearn.WeLearnCompletePointAction(cache, course, node)
	if err != nil {
		monitor.GlobalEventBus.AddLog(a.Uid, "学习异常: "+node.Location+" - "+err.Error())
		return
	}
	monitor.GlobalEventBus.AddLog(a.Uid, "✅ 任务点完成: "+node.Location)
}

// ==================== KETANGX ====================

type KetangxActivity struct{ UserActivityBase }

func (a *KetangxActivity) Login() error {
	cache := &ketangxApi.KetangxUserCache{Account: a.User.Account, Password: a.User.Password}
	err := ketangx.LoginAction(cache)
	if err != nil {
		monitor.GlobalEventBus.SetError(a.Uid, "❌ 登录失败: "+err.Error())
		return err
	}
	a.UserCache = cache
	monitor.GlobalEventBus.AddLog(a.Uid, "🔑 登录成功: "+a.User.Account)
	return nil
}

func (a *KetangxActivity) Start() error {
	a.setStopped(false)
	a.setPaused(false)
	a.setRunning(true)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusRunning)
	go a.run()
	return nil
}

func (a *KetangxActivity) Stop() error {
	a.setStopped(true)
	a.setRunning(false)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusStopped)
	return nil
}

func (a *KetangxActivity) Pause() error {
	a.setRunning(false)
	a.setPaused(true)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusPaused)
	return nil
}

func (a *KetangxActivity) run() {
	defer func() {
		a.setRunning(false)
		if r := recover(); r != nil {
			monitor.GlobalEventBus.SetError(a.Uid, fmt.Sprintf("Panic: %v", r))
			return
		}
		if !a.isStopped() && !a.isPaused() {
			monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusCompleted)
			a.sendCompletionEmail()
			monitor.GlobalEventBus.UpdateProgress(a.Uid, 100, "\u6240\u6709\u8bfe\u7a0b\u5b66\u4e60\u5b8c\u6bd5")
		}
	}()
	if a.UserCache == nil {
		monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusLogging)
		monitor.GlobalEventBus.AddLog(a.Uid, "🔑 正在登录...")
		if err := a.Login(); err != nil {
			return
		}
	}
	user := a.User
	cache := a.UserCache.(*ketangxApi.KetangxUserCache)
	courseList := ketangx.PullCourseListAction(cache)
	totalCourses := len(courseList)
	monitor.GlobalEventBus.UpdateCourseProgress(a.Uid, 0, totalCourses)
	monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("获取到 %d 门课程", totalCourses))
	for i, course := range courseList {
		if !a.IsRunning() {
			return
		}
		if len(user.CoursesCustom.ExcludeCourses) != 0 && config.CmpCourse(course.Title, user.CoursesCustom.ExcludeCourses) {
			monitor.GlobalEventBus.AddLog(a.Uid, "跳过(已排除): "+course.Title)
			continue
		}
		if len(user.CoursesCustom.IncludeCourses) != 0 && !config.CmpCourse(course.Title, user.CoursesCustom.IncludeCourses) {
			monitor.GlobalEventBus.AddLog(a.Uid, "跳过(不在包含列表): "+course.Title)
			continue
		}
		monitor.GlobalEventBus.UpdateProgress(a.Uid, float64(i)/float64(totalCourses)*100, "正在学习: "+course.Title)
		monitor.GlobalEventBus.UpdateCourseProgress(a.Uid, i, totalCourses)
		monitor.GlobalEventBus.AddLog(a.Uid, "正在学习: "+course.Title)
		nodeList := ketangx.PullNodeListAction(cache, &course)
		for _, node := range nodeList {
			if !a.IsRunning() {
				return
			}
			if user.CoursesCustom.VideoModel != 0 && !node.IsComplete {
				action, err := ketangx.CompleteVideoAction(cache, &node)
				if err != nil {
					monitor.GlobalEventBus.AddLog(a.Uid, "学习异常: "+node.Title+" - "+err.Error())
				} else {
					status := gojsonq.New().JSONString(action).Find("Success")
					if status != nil && status.(bool) {
						monitor.GlobalEventBus.AddLog(a.Uid, "学习完毕: "+node.Title)
					} else {
						monitor.GlobalEventBus.AddLog(a.Uid, "学习异常: "+node.Title+" - "+action)
					}
				}
			}
		}
		monitor.GlobalEventBus.AddLog(a.Uid, "课程完成: "+course.Title)
	}
	monitor.GlobalEventBus.AddLog(a.Uid, "所有待学习课程学习完毕")
}

// ==================== HQKJ ====================

type HqkjActivity struct{ UserActivityBase }

func (a *HqkjActivity) Login() error {
	cache := &hqkjApi.HqkjUserCache{PreUrl: a.User.URL, Account: a.User.Account, Password: a.User.Password}
	err := haiqikeji.HqkjLoginAction(cache)
	if err != nil {
		monitor.GlobalEventBus.SetError(a.Uid, "❌ 登录失败: "+err.Error())
		return err
	}
	a.UserCache = cache
	monitor.GlobalEventBus.AddLog(a.Uid, "🔑 登录成功: "+a.User.Account)
	return nil
}

func (a *HqkjActivity) Start() error {
	a.setStopped(false)
	a.setPaused(false)
	a.setRunning(true)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusRunning)
	go a.run()
	return nil
}

func (a *HqkjActivity) Stop() error {
	a.setStopped(true)
	a.setRunning(false)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusStopped)
	return nil
}

func (a *HqkjActivity) Pause() error {
	a.setRunning(false)
	a.setPaused(true)
	monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusPaused)
	return nil
}

func (a *HqkjActivity) run() {
	defer func() {
		a.setRunning(false)
		if r := recover(); r != nil {
			monitor.GlobalEventBus.SetError(a.Uid, fmt.Sprintf("Panic: %v", r))
			return
		}
		if !a.isStopped() && !a.isPaused() {
			monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusCompleted)
			a.sendCompletionEmail()
			monitor.GlobalEventBus.UpdateProgress(a.Uid, 100, "\u6240\u6709\u8bfe\u7a0b\u5b66\u4e60\u5b8c\u6bd5")
		}
	}()
	if a.UserCache == nil {
		monitor.GlobalEventBus.UpdateStatus(a.Uid, monitor.StatusLogging)
		monitor.GlobalEventBus.AddLog(a.Uid, "🔑 正在登录...")
		if err := a.Login(); err != nil {
			return
		}
	}
	user := a.User
	cache := a.UserCache.(*hqkjApi.HqkjUserCache)
	courseList, err := haiqikeji.HqkjCourseListAction(cache)
	if err != nil {
		monitor.GlobalEventBus.SetError(a.Uid, "❌ 拉取课程列表失败: "+err.Error())
		return
	}
	totalCourses := len(courseList)
	monitor.GlobalEventBus.UpdateCourseProgress(a.Uid, 0, totalCourses)
	monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("获取到 %d 门课程", totalCourses))
	for i, course := range courseList {
		if !a.IsRunning() {
			return
		}
		if len(user.CoursesCustom.ExcludeCourses) != 0 && config.CmpCourse(course.Name, user.CoursesCustom.ExcludeCourses) {
			monitor.GlobalEventBus.AddLog(a.Uid, "跳过(已排除): "+course.Name)
			continue
		}
		if len(user.CoursesCustom.IncludeCourses) != 0 && !config.CmpCourse(course.Name, user.CoursesCustom.IncludeCourses) {
			monitor.GlobalEventBus.AddLog(a.Uid, "跳过(不在包含列表): "+course.Name)
			continue
		}
		if course.StartDate.After(time.Now()) || course.EndDate.Before(time.Now()) {
			monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("跳过(不在开课时期): %s [%s ~ %s]", course.Name, course.StartDate.Format("2006-01-02"), course.EndDate.Format("2006-01-02")))
			continue
		}
		monitor.GlobalEventBus.UpdateProgress(a.Uid, float64(i)/float64(totalCourses)*100, "正在学习: "+course.Name)
		monitor.GlobalEventBus.UpdateCourseProgress(a.Uid, i, totalCourses)
		monitor.GlobalEventBus.AddLog(a.Uid, "正在学习: "+course.Name)
		nodeList, err2 := haiqikeji.HqkjNodeListAction(cache, course)
		if err2 != nil {
			continue
		}
		switch user.CoursesCustom.VideoModel {
		case 1:
			a.normalModeAction(cache, &course, nodeList)
		case 2:
			a.fastModeAction(cache, &course, nodeList)
		default:
			a.normalModeAction(cache, &course, nodeList)
		}
		monitor.GlobalEventBus.AddLog(a.Uid, "课程完成: "+course.Name)
	}
	monitor.GlobalEventBus.AddLog(a.Uid, "所有待学习课程学习完毕")
}

func (a *HqkjActivity) normalModeAction(cache *hqkjApi.HqkjUserCache, course *haiqikeji.HqkjCourse, nodeList []haiqikeji.HqkjNode) {
	for _, node := range nodeList {
		if !a.IsRunning() {
			return
		}
		if node.TabVideo <= 0 {
			continue
		}
		progress, err := haiqikeji.HqkjGetNodeProgressAction(cache, node)
		if err != nil {
			monitor.GlobalEventBus.AddLog(a.Uid, "拉取进度错误: "+err.Error())
			continue
		}
		if progress >= 100 {
			continue
		}
		sessionId, err := haiqikeji.HqkjStartStudyAction(cache, node)
		if err != nil {
			monitor.GlobalEventBus.AddLog(a.Uid, "获取sessionId失败: "+err.Error())
			return
		}
		nowTime := int(float64(progress) * 0.01 * float64(node.VideoDuration))
		stopTime := 30
		time.Sleep(time.Duration(stopTime) * time.Second)
		for {
			if !a.IsRunning() {
				return
			}
			nowAddV := stopTime
			if nowTime+stopTime > node.VideoDuration {
				nowAddV = node.VideoDuration - nowTime
			}
			nowTime += nowAddV
			submitProgress := int(float64(nowTime) / float64(node.VideoDuration) * 100)
			submitResult, err := haiqikeji.HqkjSubmitStudyTimeAction(cache, node, sessionId, submitProgress)
			if err != nil {
				monitor.GlobalEventBus.AddLog(a.Uid, "提交学时失败: "+err.Error())
			} else {
				msg := gojsonq.New().JSONString(submitResult).Find("msg")
				if msg != nil {
					monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频: %s 进度: %d%% 状态: %s", node.Name, submitProgress, msg.(string)))
				}
			}
			if submitProgress >= 100 {
				endResult, err := haiqikeji.HqkjEndStudyAction(cache, sessionId)
				if err != nil {
					monitor.GlobalEventBus.AddLog(a.Uid, "结束学习失败: "+err.Error())
					break
				}
				monitor.GlobalEventBus.AddLog(a.Uid, "结束学习: "+node.Name+" 返回: "+endResult)
				progress, err = haiqikeji.HqkjGetNodeProgressAction(cache, node)
				if progress < 100 {
					sessionId, err = haiqikeji.HqkjStartStudyAction(cache, node)
					if err != nil {
						break
					}
					time.Sleep(time.Duration(stopTime) * time.Second)
					nowTime = int(float64(node.VideoDuration) * 0.01 * float64(progress))
					continue
				}
				break
			}
			time.Sleep(time.Duration(stopTime) * time.Second)
		}
	}
}

func (a *HqkjActivity) fastModeAction(cache *hqkjApi.HqkjUserCache, course *haiqikeji.HqkjCourse, nodeList []haiqikeji.HqkjNode) {
	var videosLock sync.WaitGroup
	for _, node := range nodeList {
		videosLock.Add(1)
		go func(node haiqikeji.HqkjNode) {
			defer videosLock.Done()
			progress, err := haiqikeji.HqkjGetNodeProgressAction(cache, node)
			if err != nil {
				monitor.GlobalEventBus.AddLog(a.Uid, "拉取进度错误: "+err.Error())
				return
			}
			for {
				if !a.IsRunning() {
					return
				}
				sessionId, err := haiqikeji.HqkjStartStudyAction(cache, node)
				if err != nil {
					monitor.GlobalEventBus.AddLog(a.Uid, "获取sessionId失败: "+err.Error())
					return
				}
				time.Sleep(30 * time.Second)
				submitResult, err := haiqikeji.HqkjSubmitStudyTimeAction(cache, node, sessionId, 100)
				if err != nil {
					monitor.GlobalEventBus.AddLog(a.Uid, "提交学时失败: "+err.Error())
					return
				}
				msg := gojsonq.New().JSONString(submitResult).Find("msg")
				if msg != nil {
					monitor.GlobalEventBus.AddLog(a.Uid, fmt.Sprintf("视频(快速): %s 进度: 100%% 状态: %s", node.Name, msg.(string)))
				}
				endResult, err := haiqikeji.HqkjEndStudyAction(cache, sessionId)
				if err != nil {
					monitor.GlobalEventBus.AddLog(a.Uid, "结束学习失败: "+err.Error())
					break
				}
				monitor.GlobalEventBus.AddLog(a.Uid, "结束学习: "+node.Name+" 返回: "+endResult)
				progress, err = haiqikeji.HqkjGetNodeProgressAction(cache, node)
				if err != nil {
					return
				}
				if progress >= 100 {
					break
				}
				time.Sleep(30 * time.Second)
			}
		}(node)
	}
	videosLock.Wait()
}
