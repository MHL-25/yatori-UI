package xuexitong

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	//"github.com/yatori-dev/yatori-go-core/api/xuexitong"
	"github.com/yatori-dev/yatori-go-core/models/ctype"
	"github.com/yatori-dev/yatori-go-core/que-core/aiq"
	"github.com/yatori-dev/yatori-go-core/que-core/external"
	"github.com/yatori-dev/yatori-go-core/que-core/qentity"
	"github.com/yatori-dev/yatori-go-core/que-core/qtype"
	"github.com/yatori-dev/yatori-go-core/utils/log"
)

// XueXiTCourse 课程所有信息
type XueXiTCourseJson struct {
	Result           int           `json:"result"`
	Msg              string        `json:"msg"`
	ChannelList      []ChannelItem `json:"channelList"`
	Mcode            string        `json:"mcode"`
	Createcourse     int           `json:"createcoursed"`
	TeacherEndCourse int           `json:"teacherEndCourse"`
	ShowEndCourse    int           `json:"showEndCourse"`
	HasMore          bool          `json:"hasMore"`
	StuEndCourse     int           `json:"stuEndCourse"`
}

// ChannelItem 课程列表
type ChannelItem struct {
	Cfid     int    `json:"cfid"`
	Norder   int    `json:"norder"`
	CataName string `json:"cataName"`
	Cataid   string `json:"cataid"`
	Id       int    `json:"id"`
	Cpi      int    `json:"cpi"`
	Key      any    `json:"key"`
	Content  struct {
		Studentcount int    `json:"studentcount"`
		Chatid       string `json:"chatid"`
		IsFiled      int    `json:"isFiled"`
		Isthirdaq    int    `json:"isthirdaq"`
		Isstart      bool   `json:"isstart"`
		Isretire     int    `json:"isretire"`
		Name         string `json:"name"`
		Course       struct {
			Data []struct {
				BelongSchoolId     string `json:"belongSchoolId"`
				Coursestate        int    `json:"coursestate"`
				Teacherfactor      string `json:"teacherfactor"`
				IsCourseSquare     int    `json:"isCourseSquare"`
				Schools            string `json:"schools"`
				CourseSquareUrl    string `json:"courseSquareUrl"`
				Imageurl           string `json:"imageurl"`
				AppInfo            string `json:"appInfo"`
				Name               string `json:"name"`
				DefaultShowCatalog int    `json:"defaultShowCatalog"`
				Id                 int    `json:"id"`
				AppData            int    `json:"appData"`
			} `json:"data"`
		} `json:"course"`
		Roletype int    `json:"roletype"`
		Id       int    `json:"id"`
		State    int    `json:"state"`
		Cpi      int    `json:"cpi"`
		Bbsid    string `json:"bbsid"`
		IsSquare int    `json:"isSquare"`
	} `json:"content"`
	Topsign int `json:"topsign"`
}

// XueXiTCourse 关键信息过滤截取最终的实体
type XueXiTCourse struct {
	CourseName string //课程名称
	ClassId    string //classId
	CourseId   string //课程Id
	Cpi        string //不知道是啥玩意，反正需要
	PersonId   string //个人Id
	UserId     string //UserId
}

// ExamTopics holds a map of ExamTopic indexed by answerId
type YingHuaExamTopics struct {
	YingHuaExamTopics map[string]YingHuaExamTopic
}

// ExamTopic represents a single exam question
type YingHuaExamTopic struct {
	AnswerId string `json:"answerId"`
	Index    string `json:"index"`
	Source   string `json:"source"`
	qentity.Question
	//Content  string        `json:"content"`
	//Type     string        `json:"type"`
	//Selects  []TopicSelect `json:"selects"`
	//Answers  string        `json:"answers"`
}

// TopicSelect represents a possible answer choice
//type TopicSelect struct {
//	Value string `json:"value"`
//	Num   string `json:"num"`
//	Text  string `json:"text"`
//}

// QuestionMedia 题目图片（题干图 / 选项图）
type QuestionMedia struct {
	ContentImages []qentity.QuestionImage            `json:"contentImages,omitempty"`
	OptionImages  map[string][]qentity.QuestionImage `json:"optionImages,omitempty"`
}

// HasQuestionImages 是否包含图片
func (m *QuestionMedia) HasQuestionImages() bool {
	if len(m.ContentImages) > 0 {
		return true
	}
	for _, imgs := range m.OptionImages {
		if len(imgs) > 0 {
			return true
		}
	}
	return false
}

// ChoiceQue 选择类型
type ChoiceQue struct {
	QuestionMedia
	Type    qtype.QType       `json:"type"`
	Qid     string            `json:"qid"` //题目ID
	Text    string            `json:"text"`
	Options map[string]string `json:"options"`
	Answers []string          `json:"answers"` // 答案
}

// JudgeQue 判断类型
type JudgeQue struct {
	QuestionMedia
	Type    qtype.QType       `json:"type"`
	Qid     string            `json:"qid"` //题目ID
	Text    string            `json:"text"`
	Options map[string]string `json:"options"`
	Answers []string          `json:"answers"` // 答案
}

// FillQue 填空类型
type FillQue struct {
	QuestionMedia
	Type qtype.QType `json:"type"`
	Qid  string      `json:"qid"`
	Text string      `json:"text"`
	//OpFromAnswer map[string][]string `json:"opFromAnswer"` // 位置与答案
	OpFromAnswer []string `json:"opFromAnswer"` // 位置与答案
	BlankCount   int      `json:"blankCount"`   // 原始空位数量，答题时保持不漂移
}

// 简答类型
type ShortQue struct {
	QuestionMedia
	Type         qtype.QType         `json:"type"`
	Qid          string              `json:"qid"`
	Text         string              `json:"text"`
	OpFromAnswer map[string][]string `json:"opFromAnswer"`
}

// 名词解释类型
type TermExplanationQue struct {
	QuestionMedia
	Type         qtype.QType         `json:"type"`
	Qid          string              `json:"qid"`
	Text         string              `json:"text"`
	OpFromAnswer map[string][]string `json:"opFromAnswer"`
}

// 论述题
type EssayQue struct {
	QuestionMedia
	Type         qtype.QType         `json:"type"`
	Qid          string              `json:"qid"`
	Text         string              `json:"text"`
	OpFromAnswer map[string][]string `json:"opFromAnswer"`
}

// 连线题
type MatchingQue struct {
	QuestionMedia
	Type    qtype.QType `json:"type"`
	Qid     string      `json:"qid"` //题目ID
	Text    string      `json:"text"`
	Options []string    `json:"options1"`
	Selects []string    `json:"options2"`
	Answers []string    `json:"answers"` // 答案
}

// “其他”题
type OtherQue struct {
	QuestionMedia
	Type         qtype.QType         `json:"type"`
	Qid          string              `json:"qid"`
	Text         string              `json:"text"`
	OpFromAnswer map[string][]string `json:"opFromAnswer"`
}

// Question TODO 这里考虑是否在其中直接将答案做出 直接上报提交 或 保存提交
type Question struct {
	Title            string               `json:"title"` //试卷标题
	Cpi              string               `json:"cpi"`
	JobId            string               `json:"jobId"`
	WorkId           string               `json:"workId"`
	ClassId          string               `json:"classId"`
	CourseId         string               `json:"courseId"`
	Ua               string               `json:"ua"`
	FormType         string               `json:"formType"`
	SaveStatus       string               `json:"saveStatus"`
	Version          string               `json:"version"`
	Tempsave         string               `json:"tempsave"`
	PyFlag           string               `json:"pyFlag"`
	UserId           string               `json:"userId"`
	Knowledgeid      string               `json:"knowledgeId"`
	OldWorkId        string               `json:"oldWorkId"`   //最原始作业id
	FullScore        string               `json:"fullScore"`   //满分是多少
	OldSchoolId      string               `json:"oldSchoolId"` //原始作业单位id
	Api              string               `json:"api"`         //api值
	WorkRelationId   string               `json:"workRelationId"`
	Enc_work         string               `json:"enc_Work"`
	Isphone          string               `json:"isphone"`
	RandomOptions    string               `json:"randomOptions"`
	WorkAnswerId     string               `json:"workAnswerId"`
	AnswerId         string               `json:"answerId"`
	TotalQuestionNum string               `json:"totalQuestionNum"`
	Choice           []ChoiceQue          //选择类型
	Judge            []JudgeQue           //判断类型
	Fill             []FillQue            //填空类型
	Short            []ShortQue           //简答类型
	TermExplanation  []TermExplanationQue //名词解释类型
	Essay            []EssayQue           //论述题类型
	Matching         []MatchingQue        //连线题
	Other            []OtherQue           //其他题
}

// 序列化输出答题信息
func (f Question) String() string {
	marshal, err := json.Marshal(f)
	if err != nil {
		return "题目序列化失败: " + err.Error()
	}
	return string(marshal)
}

type ExamTurn struct {
	XueXChoiceQue          ChoiceQue
	XueXJudgeQue           JudgeQue
	XueXFillQue            FillQue
	XueXShortQue           ShortQue
	XueXTermExplanationQue TermExplanationQue
	XueXEssayQue           EssayQue
	XueXMatchingQue        MatchingQue
	XueXOtherQue           OtherQue
	YingHuaExamTopic
}

type AnswerSetter interface {
	SetAnswers([]string)
}

func (q *ChoiceQue) SetAnswers(answers []string) {
	q.Answers = answers
}

func (q *JudgeQue) SetAnswers(answers []string) {
	q.Answers = answers
}

func (q *FillQue) SetAnswers(answers []string) {
	if len(answers) == 0 {
		return
	}
	// 保持原始空位数量，避免AI多给或少给答案时空位整体错位
	blankCount := q.BlankCount
	if blankCount <= 0 {
		blankCount = len(q.OpFromAnswer)
	}
	if blankCount > 0 {
		normalized := make([]string, blankCount)
		copy(normalized, answers)
		q.OpFromAnswer = normalized
		return
	}
	q.OpFromAnswer = answers
}

func (q *ShortQue) SetAnswers(answers []string) {
	q.OpFromAnswer["简答"] = answers
}

func (q *TermExplanationQue) SetAnswers(answers []string) {
	q.OpFromAnswer["名词解释"] = answers
}
func (q *EssayQue) SetAnswers(answers []string) {
	q.OpFromAnswer["论述"] = answers
}
func (q *MatchingQue) SetAnswers(answers []string) {
	q.Answers = answers
}

func (q *OtherQue) SetAnswers(answers []string) {
	q.OpFromAnswer["回复"] = answers
}

// 从键中提取序号（例如："0第3空" → 2，注意索引从0开始）
func extractIndexFromKey(key string) int {
	// 简单实现，实际可能需要更复杂的字符串处理
	// 这里假设key格式为"0第X空"，提取X并转换为整数
	// 示例实现，需要根据实际格式调整
	regex := regexp.MustCompile(`第(\d+)空`)
	matches := regex.FindStringSubmatch(key)
	if len(matches) >= 2 {
		if idx, err := strconv.Atoi(matches[1]); err == nil {
			return idx - 1 // 转换为0-based索引
		}
	}
	return -1 // 无效索引
}

func GetAIAnswer(as AnswerSetter, userID string, url, model string, aiType ctype.AiType, aiChatMessages aiq.AIChatMessages, apiKey string) {
	aiAnswer, err := aiq.AggregationAIApi(url, model, aiType, aiChatMessages, apiKey)
	if err != nil {
		log.Print(log.INFO, `[`, userID, `] `, log.BoldRed, "Ai异常，返回信息：", err.Error())
		//os.Exit(0)
	}

	answers := aiq.ParseStringArray(aiAnswer)
	if len(answers) == 0 {
		answers = []string{"A"}
		log.Print(log.INFO, "AI回复解析为空，已采用随机答案:", fmt.Sprintf("题目：%v \nAI回复： %v", aiChatMessages, aiAnswer))
	}

	as.SetAnswers(answers)
}

// 学习通AI答题回复
func GetXXTAIAnswer(as AnswerSetter, cache *XueXiTUserCache, classId, courseId, cpi string, aiChatMessages aiq.AIChatMessages) {

	aiAnswer, err := cache.XueXiTongAIAggregation(classId, courseId, cpi, aiChatMessages)
	if err != nil {
		log.Print(log.INFO, `[`, cache.Name, `] `, log.BoldRed, "学习通内置AI异常，返回信息：", err.Error())
		as.SetAnswers([]string{"A"})
		return
	}
	answers := aiq.ParseStringArray(aiAnswer)
	if len(answers) == 0 {
		answers = []string{"A"}
		log.Print(log.INFO, "AI回复解析为空，已采用随机答案:", fmt.Sprintf("题目：%v \nAI回复： %v", aiChatMessages, aiAnswer))
	}

	as.SetAnswers(answers)
}

// AnswerAIGet ChoiceQue的AI回答获取方法
func (q *ChoiceQue) AnswerAIGet(userID,
	url, model string, aiType ctype.AiType, aiChatMessages aiq.AIChatMessages, apiKey string) {
	GetAIAnswer(q, userID, url, model, aiType, aiChatMessages, apiKey)
}

// 调用学习通内部AI答题
func (q *ChoiceQue) AnswerXXTAIGet(cache *XueXiTUserCache, classId, courseId, cpi string, aiChatMessages aiq.AIChatMessages) {
	GetXXTAIAnswer(q, cache, classId, courseId, cpi, aiChatMessages)
}

// AnswerExternalGet ChoiceQue的外挂题库回答获取方法
func (q *ChoiceQue) AnswerExternalGet(exUrl string) {
	question := qentity.Question{
		Type:    q.Type.String(),
		Content: q.Text,
	}
	//赋值选项
	for _, k := range q.Options {
		question.Options = append(question.Options, k)
	}
	request, err := external.ApiQueRequest(question, exUrl, 5, nil)
	if err != nil {
		log.Print(log.INFO, "外置题库请求失败，已跳过本题:", err.Error())
		return
	}
	//赋值答案
	var answers []string
	if request.Code == 404 {
		answers = []string{"A"}
		log.Print(log.INFO, "外置题库未找到答案，已使用默认答案A:", fmt.Sprintf("题目：%v \n外置题库回复： %v", q, request))
		request.Question.Answers = answers
		return
	}
	if request.Question.Answers == nil {
		answers = []string{"A"}
		fmt.Println("回复解析错误:", err, fmt.Sprintf("题目：%v \n外置题库回复： %v", q, request))
		request.Question.Answers = answers
	}
	q.Answers = request.Question.Answers
}

// AnswerAIGet JudgeQue的AI回答获取方法
func (q *JudgeQue) AnswerAIGet(userID,
	url, model string, aiType ctype.AiType, aiChatMessages aiq.AIChatMessages, apiKey string) {
	GetAIAnswer(q, userID, url, model, aiType, aiChatMessages, apiKey)

}

// 调用学习通内部AI答题
func (q *JudgeQue) AnswerXXTAIGet(cache *XueXiTUserCache, classId, courseId, cpi string, aiChatMessages aiq.AIChatMessages) {
	GetXXTAIAnswer(q, cache, classId, courseId, cpi, aiChatMessages)
}
func (q *JudgeQue) AnswerExternalGet(exUrl string) {
	question := qentity.Question{
		Type:    q.Type.String(),
		Content: q.Text,
	}
	//赋值选项
	for _, k := range q.Options {
		question.Options = append(question.Options, k)
	}
	request, err := external.ApiQueRequest(question, exUrl, 5, nil)
	if err != nil {
		log.Print(log.INFO, "外置题库请求失败，已跳过本题:", err.Error())
		return
	}
	if request.Code == 404 {
		log.Print(log.INFO, "外置题库未找到答案，已使用默认答案：正确:", fmt.Sprintf("题目：%v \n外置题库回复： %v", q, request))
		request.Question.Answers = []string{"对"}
		return
	}

	var answers []string
	if request.Question.Answers == nil {
		answers = []string{"A"}
		fmt.Println("回复解析错误:", err, fmt.Sprintf("题目：%v \n外置题库回复： %v", q, request))
		request.Question.Answers = answers
	}
	//赋值答案
	q.Answers = request.Question.Answers
}

// AnswerAIGet FillQue的AI回答获取方法
func (q *FillQue) AnswerAIGet(userID,
	url, model string, aiType ctype.AiType, aiChatMessages aiq.AIChatMessages, apiKey string) {
	GetAIAnswer(q, userID, url, model, aiType, aiChatMessages, apiKey)
}

// 调用学习通内部AI答题
func (q *FillQue) AnswerXXTAIGet(cache *XueXiTUserCache, classId, courseId, cpi string, aiChatMessages aiq.AIChatMessages) {
	GetXXTAIAnswer(q, cache, classId, courseId, cpi, aiChatMessages)
}
func (q *FillQue) AnswerExternalGet(exUrl string) {
	question := qentity.Question{
		Type:    q.Type.String(),
		Content: q.Text,
	}
	//赋值选项

	request, err := external.ApiQueRequest(question, exUrl, 5, nil)
	if err != nil {
		log.Print(log.INFO, "外置题库请求失败，已跳过本题:", err.Error())
		return
	}
	if request.Code == 404 {
		log.Print(log.INFO, "外置题库未找到答案，已自动留空", fmt.Sprintf("题目：%v \n外置题库回复： %v", q, request))
		request.Question.Answers = []string{}
		for i := 0; i < len(q.OpFromAnswer); i++ {
			request.Question.Answers = append(request.Question.Answers, "")
		}
		return
	}
	q.OpFromAnswer = request.Question.Answers

}

// AnswerAIGet ShortQue的AI回答获取方法
func (q *ShortQue) AnswerAIGet(userID,
	url, model string, aiType ctype.AiType, aiChatMessages aiq.AIChatMessages, apiKey string) {
	GetAIAnswer(q, userID, url, model, aiType, aiChatMessages, apiKey)
}

// 调用学习通内部AI答题
func (q *ShortQue) AnswerXXTAIGet(cache *XueXiTUserCache, classId, courseId, cpi string, aiChatMessages aiq.AIChatMessages) {
	GetXXTAIAnswer(q, cache, classId, courseId, cpi, aiChatMessages)
}
func (q *ShortQue) AnswerExternalGet(exUrl string) {
	question := qentity.Question{
		Type:    q.Type.String(),
		Content: q.Text,
	}
	//赋值选项
	request, err := external.ApiQueRequest(question, exUrl, 5, nil)
	if err != nil {
		log.Print(log.INFO, "外置题库请求失败，已跳过本题:", err.Error())
		return
	}
	var answers []string
	if request.Question.Answers == nil {
		answers = []string{"A"}
		fmt.Println("回复解析错误:", err, fmt.Sprintf("题目：%v \n外置题库回复： %v", q, request))
		request.Question.Answers = answers
	}
	//赋值答案
	q.SetAnswers(request.Question.Answers)
}

// AnswerAIGet TermExplanation的AI回答获取方法
func (q *TermExplanationQue) AnswerAIGet(userID,
	url, model string, aiType ctype.AiType, aiChatMessages aiq.AIChatMessages, apiKey string) {
	GetAIAnswer(q, userID, url, model, aiType, aiChatMessages, apiKey)
}

// 调用学习通内部AI答题
func (q *TermExplanationQue) AnswerXXTAIGet(cache *XueXiTUserCache, classId, courseId, cpi string, aiChatMessages aiq.AIChatMessages) {
	GetXXTAIAnswer(q, cache, classId, courseId, cpi, aiChatMessages)
}

func (q *TermExplanationQue) AnswerExternalGet(exUrl string) {
	question := qentity.Question{
		Type:    q.Type.String(),
		Content: q.Text,
	}
	//赋值选项
	request, err := external.ApiQueRequest(question, exUrl, 5, nil)
	if err != nil {
		log.Print(log.INFO, "外置题库请求失败，已跳过本题:", err.Error())
		return
	}
	//处理回复为空的情况
	if len(request.Question.Answers) == 0 {
		log.Print(log.INFO, "回复数据为空，已自动回复默认值A:", fmt.Sprintf("%x", request))
		request.Question.Answers = []string{"A"}
	}
	//赋值答案
	q.SetAnswers(request.Question.Answers)
}

// AnswerAIGet Essay的AI回答获取方法
func (q *EssayQue) AnswerAIGet(userID,
	url, model string, aiType ctype.AiType, aiChatMessages aiq.AIChatMessages, apiKey string) {
	GetAIAnswer(q, userID, url, model, aiType, aiChatMessages, apiKey)
}

// 调用学习通内部AI答题
func (q *EssayQue) AnswerXXTAIGet(cache *XueXiTUserCache, classId, courseId, cpi string, aiChatMessages aiq.AIChatMessages) {
	GetXXTAIAnswer(q, cache, classId, courseId, cpi, aiChatMessages)
}

func (q *EssayQue) AnswerExternalGet(exUrl string) {
	question := qentity.Question{
		Type:    q.Type.String(),
		Content: q.Text,
	}
	//赋值选项
	request, err := external.ApiQueRequest(question, exUrl, 5, nil)
	if err != nil {
		log.Print(log.INFO, "外置题库请求失败，已跳过本题:", err.Error())
		return
	}

	//处理回复为空的情况
	if len(request.Question.Answers) == 0 {
		log.Print(log.INFO, "回复数据为空，已自动回复默认值A:", fmt.Sprintf("%x", request))
		request.Question.Answers = []string{"A"}
	}
	//赋值答案
	q.SetAnswers(request.Question.Answers)
}

func (q *MatchingQue) AnswerAIGet(userID,
	url, model string, aiType ctype.AiType, aiChatMessages aiq.AIChatMessages, apiKey string) {
	aiAnswer, err := aiq.AggregationAIApi(url, model, aiType, aiChatMessages, apiKey)
	if err != nil {
		log.Print(log.INFO, `[`, userID, `] `, log.BoldRed, "连线题Ai异常，返回信息：", err.Error())
	}
	q.Answers = ParseMatchingAnswers(aiAnswer)
}

// 调用学习通内部AI答题
func (q *MatchingQue) AnswerXXTAIGet(cache *XueXiTUserCache, classId, courseId, cpi string, aiChatMessages aiq.AIChatMessages) {
	aiAnswer, err := cache.XueXiTongAIAggregation(classId, courseId, cpi, aiChatMessages)
	if err != nil {
		log.Print(log.INFO, `[`, cache.Name, `] `, log.BoldRed, "连线题学习通AI异常，返回信息：", err.Error())
		return
	}
	q.Answers = ParseMatchingAnswers(aiAnswer)
}

// ParseMatchingAnswers 解析连线题AI返回，兼容 JSON 数组、对象数组、带箭头文本和代码块包裹
func ParseMatchingAnswers(aiAnswer string) []string {
	raw := strings.TrimSpace(aiAnswer)
	if raw == "" {
		return nil
	}
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```JSON")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	if start := strings.Index(raw, "["); start >= 0 {
		if end := strings.LastIndex(raw, "]"); end > start {
			candidate := raw[start : end+1]
			var answers []string
			if err := json.Unmarshal([]byte(candidate), &answers); err == nil && len(answers) > 0 {
				return answers
			}
			var objects []struct {
				Name    any    `json:"name"`
				Content string `json:"content"`
			}
			if err := json.Unmarshal([]byte(candidate), &objects); err == nil && len(objects) > 0 {
				out := make([]string, 0, len(objects))
				for _, item := range objects {
					if strings.TrimSpace(item.Content) != "" {
						out = append(out, fmt.Sprintf("%v->%s", item.Name, item.Content))
					}
				}
				if len(out) > 0 {
					return out
				}
			}
		}
	}

	out := make([]string, 0)
	for _, line := range strings.Split(raw, "\n") {
		line = strings.Trim(strings.TrimSpace(line), `"',`)
		if line == "" {
			continue
		}
		if strings.Contains(line, "->") || strings.Contains(line, "→") || strings.Contains(line, "=>") {
			out = append(out, line)
		}
	}
	return out
}

func (q *MatchingQue) AnswerExternalGet(exUrl string) {
	question := qentity.Question{
		Type:    q.Type.String(),
		Content: q.Text,
	}
	//赋值选项
	request, err := external.ApiQueRequest(question, exUrl, 5, nil)
	if err != nil {
		log.Print(log.INFO, "外置题库请求失败，已跳过本题:", err.Error())
		return
	}
	//赋值答案
	q.SetAnswers(request.Question.Answers)
}

// AnswerAIGet Other的AI回答获取方法
func (q *OtherQue) AnswerAIGet(userID,
	url, model string, aiType ctype.AiType, aiChatMessages aiq.AIChatMessages, apiKey string) {
	GetAIAnswer(q, userID, url, model, aiType, aiChatMessages, apiKey)
}

// 调用学习通内部AI答题
func (q *OtherQue) AnswerXXTAIGet(cache *XueXiTUserCache, classId, courseId, cpi string, aiChatMessages aiq.AIChatMessages) {
	GetXXTAIAnswer(q, cache, classId, courseId, cpi, aiChatMessages)
}

func (q *OtherQue) AnswerExternalGet(exUrl string) {
	question := qentity.Question{
		Type:    q.Type.String(), //直接采用论述题方案
		Content: q.Text,
	}
	//赋值选项
	request, err := external.ApiQueRequest(question, exUrl, 5, nil)
	if err != nil {
		log.Print(log.INFO, "外置题库请求失败，已跳过本题:", err.Error())
		return
	}
	//赋值答案
	q.SetAnswers(request.Question.Answers)
}
