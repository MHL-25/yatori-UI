package xuexitong

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/yatori-dev/yatori-go-core/api/xuexitong"
	"github.com/yatori-dev/yatori-go-core/models/ctype"
	"github.com/yatori-dev/yatori-go-core/que-core/aiq"
	"github.com/yatori-dev/yatori-go-core/que-core/external"
	"github.com/yatori-dev/yatori-go-core/que-core/qentity"
	"github.com/yatori-dev/yatori-go-core/que-core/qtype"
	log2 "github.com/yatori-dev/yatori-go-core/utils/log"
	"github.com/yatori-dev/yatori-go-core/utils"
)

// 学习通考试结构体
type XXTWork struct {
	Name              string `json:"name"`
	Status            string `json:"status"`
	RemainTime        string `json:"remain_time"`
	RawURL            string `json:"raw_url"`
	Params            map[string]string
	CourseId          string `json:"course_id"`
	UserId            string `json:"user_id"`
	ClazzId           string `json:"clazz_id"`
	Type              string `json:"type"`
	EncTask           string `json:"enc_task"`
	TaskRefId         string `json:"taskrefId"`
	MsgId             string `json:"msgId"`
	CaptchaCaptchaId  string
	ExamRelationId    string
	AnswerId          string `json:"answerId"`
	Cpi               string
	Validate          string //过验证码用的
	QuestionTotal     int    //题目数量
	Enc               string
	EncRemainTime     string
	EncLastUpdateTime string
}

type XXTWorkQuestion struct {
	xuexitong.XXTWorkQuestionSubmitEntity
}

// 拉取考试列表
func PullWorkListAction(cache *xuexitong.XueXiTUserCache, course XueXiTCourse) ([]XXTWork, error) {
	examList := []XXTWork{}
	examListHtml, err := cache.PullWorkListHtmlApi(course.CourseID, course.Key, fmt.Sprintf("%d", course.Cpi), 3, nil)
	if err != nil {
		return examList, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(examListHtml))
	if err != nil {
		return examList, fmt.Errorf("作业列表页面解析失败: %v", err)
	}

	// 遍历所有 <li>
	doc.Find("ul.nav li").Each(func(i int, li *goquery.Selection) {
		rawURL, _ := li.Attr("data")

		// 解析 URL 参数
		parsed, _ := url.Parse(rawURL)
		params := map[string]string{}
		for k, v := range parsed.Query() {
			if len(v) > 0 {
				params[k] = v[0]
			}
		}

		div := li.Find("div")

		name := strings.TrimSpace(div.Find("p").Text())

		spanList := div.Find("span")
		status := strings.TrimSpace(spanList.Eq(0).Text())

		remain := ""
		if spanList.Length() > 1 {
			remain = strings.TrimSpace(spanList.Eq(1).Text())
		}

		exam := XXTWork{
			Name:       name,
			Status:     status,
			RemainTime: remain,
			RawURL:     rawURL,
			Params:     params,
			TaskRefId:  params["taskrefId"],
			CourseId:   params["courseId"],
			UserId:     params["userId"],
			ClazzId:    params["clazzId"],
			Type:       params["type"],
			EncTask:    params["enc_task"],
			MsgId:      params["msgId"],
		}
		examList = append(examList, exam)
	})

	return examList, nil
}

// EnterWorkAction 进入作业
func EnterWorkAction(cache *xuexitong.XueXiTUserCache, exam *XXTWork) error {
	//这一步拉取必要的参数，比如滑块验证码参数等,注意这里的refererUrl会在后面的滑块验证码中用到
	enterHtml, refererUrl, err := cache.PullWorkEnterInformHtmlApi(exam.TaskRefId, exam.MsgId, exam.CourseId, exam.UserId, exam.ClazzId, exam.Type, exam.EncTask, 3, nil)
	if err != nil {
		//fmt.Println(refererUrl)
		return err
	}

	re := regexp.MustCompile(`共包含\s*(\d+)\s*道题目`)
	match := re.FindStringSubmatch(enterHtml)

	if len(match) > 1 {
		count, _ := strconv.Atoi(match[1])
		exam.QuestionTotal = count
	}
	//如果是待重做
	if strings.Contains(enterHtml, "待重做") {
		re = regexp.MustCompile(`共\s*(\d+)\s*题`)
		match = re.FindStringSubmatch(enterHtml)

		if len(match) > 1 {
			count, _ := strconv.Atoi(match[1])
			exam.QuestionTotal = count
		}
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(enterHtml))
	if err != nil {
		return fmt.Errorf("作业进入页面解析失败: %v", err)
	}
	type HiddenField struct {
		ID    string
		Name  string
		Value string
	}

	var fields []HiddenField
	// 选择所有隐藏字段
	doc.Find("input[type='hidden']").Each(func(i int, sel *goquery.Selection) {
		id, _ := sel.Attr("id")
		name, _ := sel.Attr("name")
		value, _ := sel.Attr("value")

		fields = append(fields, HiddenField{
			ID:    id,
			Name:  name,
			Value: value,
		})
	})

	// 输出结果
	for _, f := range fields {
		if f.ID == "captchaCaptchaId" {
			exam.CaptchaCaptchaId = f.Value
		}
		if f.ID == "testPaperId" {
			exam.ExamRelationId = f.Value
		}
		if f.ID == "testUserRelationId" {
			exam.AnswerId = f.Value
		}
		if f.ID == "cpi" {
			exam.Cpi = f.Value
		}
		//fmt.Printf("ID: %-35s Name: %-20s Value: %s\n", f.ID, f.Name, f.Value)
	}
	if exam.CaptchaCaptchaId != "" {
		slider := XueXiTSlider{
			CaptchaId: exam.CaptchaCaptchaId,
			Referer:   refererUrl,
		}
		for {
			validate, passErr := slider.Pass(cache)
			if passErr != nil {
				if strings.Contains(passErr.Error(), `"result":false`) {
					continue
				}
			}
			exam.Validate = validate
			break //如果成功了那么直接退出循环
		}
	}
	cpi, answerId, enc := extractParams(enterHtml)
	exam.Cpi = cpi
	exam.AnswerId = answerId
	exam.Enc = enc
	//pullPaperHtml, err := cache.PullExamPaperHtmlApi(exam.CourseId, exam.ClazzId, exam.ExamRelationId, "0", exam.AnswerId, exam.Cpi, "1", xuexitong.IMEI, exam.Validate, "0", 3, nil)
	pullPaperHtml, err := cache.PullWorkPaperHtmlApi(exam.CourseId, exam.ClazzId, exam.TaskRefId, "0", exam.MsgId, exam.Cpi, exam.AnswerId, exam.Enc, "1", 3, nil)
	if err != nil {
		return err
	}
	if strings.Contains(pullPaperHtml, "已过时效，不能操作!") {
		return errors.New("已过时效，不能操作!")
	}
	qsEntity, err1 := HtmlWorkQuestionTurnEntity(pullPaperHtml)
	if err1 != nil {
		return err1
	}
	exam.Enc = qsEntity.Enc
	exam.EncRemainTime = qsEntity.EncRemainTime
	exam.EncLastUpdateTime = qsEntity.EncLastUpdateTime

	//fmt.Println(pullPaperHtml)
	return nil
}

func extractParams(js string) (cpi, workAnswerId, enc string) {
	//cpiRe := regexp.MustCompile(`cpi=(\d+).*?workAnswerId=(\d+).*?enc=([a-fA-F0-9]+)`)
	cpiRe := regexp.MustCompile(`cpi=(\d+).*?`)
	cpiMatch := cpiRe.FindStringSubmatch(js)
	answerIdRe := regexp.MustCompile(`workAnswerId=(\d+).*?`)
	answerIdMatch := answerIdRe.FindStringSubmatch(js)
	encRe := regexp.MustCompile(`enc=([a-fA-F0-9]+)`)
	encMatch := encRe.FindStringSubmatch(js)
	cpiC := ""
	workAnswerIdC := ""
	encC := ""
	if len(cpiMatch) > 1 {
		cpiC = cpiMatch[1]
	}
	if len(answerIdMatch) > 1 {
		workAnswerIdC = answerIdMatch[1]
	}
	if len(encMatch) > 1 {
		encC = encMatch[1]
	}
	return cpiC, workAnswerIdC, encC
}

// 拉取题目
func (exam *XXTWork) PullWorkQuestionAction(cache *xuexitong.XueXiTUserCache, index int /*第几道题*/) (XXTWorkQuestion, error) {
	pullQuestion, err1 := cache.PullWorkQuestionApi(exam.CourseId, exam.ClazzId, exam.TaskRefId, "0", exam.MsgId, exam.Cpi, exam.AnswerId, exam.Enc, "1", index, 3, nil)
	if err1 != nil {
		return XXTWorkQuestion{}, err1
	}
	//isLastQuestion := strings.Contains(pullQuestion, `<div class="lastQuestion"> 已经是最后一题了</div>`)
	qsEntity, err := HtmlWorkQuestionTurnEntity(pullQuestion)
	//fmt.Println(pullPaperHtml)
	if err != nil {
		return XXTWorkQuestion{}, err
	}
	qsEntity.ExamRelationId = exam.ExamRelationId
	qsEntity.AnswerId = exam.AnswerId
	qsEntity.RemainTimeParam = exam.RemainTime
	qsEntity.Tid = exam.TaskRefId
	return qsEntity, nil
}

// AI写题
func (question *XXTWorkQuestion) WriteQuestionForAIAction(cache *xuexitong.XueXiTUserCache, aiUrl, model string, aiType ctype.AiType, apiKey string, visionAiUrl, visionModel string, visionAiType ctype.AiType, visionApiKey string) error {
	var aiAnswer string
	var err error

	xuexitong.PrepareQuestionImages(cache, &question.Question)
	if question.Question.HasImages() && visionApiKey != "" {
		imageDataMap := make(map[string]qentity.QuestionImage)
		cookies := cache.GetCookies()
		for _, img := range question.Question.ContentImages {
			if img.Url != "" && img.Base64 == "" {
				downloaded, dlErr := utils.DownloadImageAsQuestionImage(img.Url, cookies)
				if dlErr == nil {
					imageDataMap[img.Url] = downloaded
				}
			} else {
				imageDataMap[img.Url] = img
			}
		}
		for _, imgs := range question.Question.OptionImages {
			for _, img := range imgs {
				if img.Url != "" && img.Base64 == "" {
					downloaded, dlErr := utils.DownloadImageAsQuestionImage(img.Url, cookies)
					if dlErr == nil {
						imageDataMap[img.Url] = downloaded
					}
				} else {
					imageDataMap[img.Url] = img
				}
			}
		}

		visualMessages := aiq.BuildVisualAiQuestionMessage(question.Question, imageDataMap)
		aiAnswer, err = aiq.AggregationVisualAIApi(visionAiUrl, visionModel, visionAiType, visualMessages, visionApiKey)
		if err != nil {
			log2.Print(log2.INFO, `[`, cache.Name, `] `, log2.BoldRed, "视觉AI异常，尝试本地OCR兜底：", err.Error())
			ocrFallbackQuestion(question)
			aiChatMessages := aiq.BuildAiQuestionMessage(question.Question)
			aiAnswer, err = aiq.AggregationAIApi(aiUrl, model, aiType, aiChatMessages, apiKey)
		}
	} else {
		if question.Question.HasImages() {
			ocrFallbackQuestion(question)
		}
		aiChatMessages := aiq.BuildAiQuestionMessage(question.Question)
		aiAnswer, err = aiq.AggregationAIApi(aiUrl, model, aiType, aiChatMessages, apiKey)
	}

	if err != nil {
		log2.Print(log2.INFO, `[`, cache.Name, `] `, log2.BoldRed, "Ai异常，返回信息：", err.Error())
	}

	answers := aiq.ParseQuestionAnswers(question.Question, aiAnswer)
	if len(answers) == 0 {
		if question.Question.Type == qtype.FillInTheBlank.String() {
			answers = make([]string, len(question.Question.Options))
		} else {
			answers = []string{"A"}
		}
		log2.Print(log2.INFO, "AI回复解析为空，已采用保底答案:", fmt.Sprintf("题目：%v \nAI回复： %v", question.Question.Content, aiAnswer))
	}
	question.Question.Answers = answers
	return nil
}

// WriteQuestionForVisionOrXXTAIAction 有图优先识图，无图或识图失败继续使用学习通内置AI
func (question *XXTWorkQuestion) WriteQuestionForVisionOrXXTAIAction(cache *xuexitong.XueXiTUserCache,
	visionAiUrl, visionModel string, visionAiType ctype.AiType, visionApiKey string,
	classId, courseId, cpi string) error {
	if !question.Question.HasImages() {
		return question.WriteQuestionForXXTAIAction(cache, classId, courseId, cpi)
	}
	xuexitong.PrepareQuestionImages(cache, &question.Question)
	if visionApiKey != "" && visionAiUrl != "" {
		setter := questionAnswerSetter{answers: &question.Question.Answers}
		err := xuexitong.AnswerQuestionWithVision(cache, question.Question, setter,
			visionAiUrl, visionModel, visionAiType, visionApiKey)
		if err == nil {
			return nil
		}
		log2.Print(log2.INFO, `[`, cache.Name, `] `, log2.BoldRed, "识图AI异常，尝试OCR后使用学习通内置AI：", err.Error())
	}
	if ocrText := utils.OCRQuestionImages(question.Question); ocrText != "" {
		question.Question.Content += "\n[图片OCR识别内容]\n" + ocrText
	}
	return question.WriteQuestionForXXTAIAction(cache, classId, courseId, cpi)
}

type questionAnswerSetter struct {
	answers *qentity.StringArray
}

func (s questionAnswerSetter) SetAnswers(answers []string) {
	*s.answers = qentity.StringArray(answers)
}

// ocrFallbackQuestion 图片题无法使用识图模型时，尝试本地 OCR 并把文字追加到题干
func ocrFallbackQuestion(question *XXTWorkQuestion) {
	if !question.Question.HasImages() {
		return
	}
	if ocrText := utils.OCRQuestionImages(question.Question); ocrText != "" {
		question.Question.Content += "\n[图片OCR识别内容]\n" + ocrText
	}
}

// 外置题库写题
func (question *XXTWorkQuestion) WriteQuestionForExternalAction(exUrl string) {
	//赋值选项
	request, err := external.ApiQueRequest(question.Question, exUrl, 5, nil)
	if err != nil {
		log2.Print(log2.INFO, "外置题库请求失败，已使用默认答案A:", err.Error())
		question.Question.Answers = []string{"A"}
		return
	}
	//赋值答案
	var answers []string
	if request.Code == 404 {
		answers = []string{"A"}
		log2.Print(log2.INFO, "外置题库未找到答案，已使用默认答案A:", fmt.Sprintf("题目：%v \n外置题库回复： %v", question.Question, request))
		request.Question.Answers = answers
		return
	}
	if request.Question.Answers == nil {
		answers = []string{"A"}
		fmt.Println("回复解析错误:", err, fmt.Sprintf("题目：%v \n外置题库回复： %v", question.Question, request))
		request.Question.Answers = answers
	}
	question.Question.Answers = request.Question.Answers
}

// 学习通内置题库写题
func (question *XXTWorkQuestion) WriteQuestionForXXTAIAction(cache *xuexitong.XueXiTUserCache, classId, courseId, cpi string) error {
	aiChatMessages := aiq.BuildAiQuestionMessage(question.Question)

	aiAnswer, err := cache.XueXiTongAIAggregation(classId, courseId, cpi, aiChatMessages)

	if err != nil {
		log2.Print(log2.INFO, `[`, cache.Name, `] `, log2.BoldRed, "学习通内置AI作答失败：", err.Error())
		return err
	}
	answers := aiq.ParseStringArray(aiAnswer)
	if len(answers) == 0 {
		answers = []string{"A"}
		log2.Print(log2.INFO, "AI回复解析为空，已采用随机答案:", fmt.Sprintf("题目：%v \nAI回复： %v", aiChatMessages, aiAnswer))
	}
	question.Question.Answers = answers
	return nil
}

// 提交学习通考试答案
func (question *XXTWorkQuestion) SubmitWorkAnswerAction(cache *xuexitong.XueXiTUserCache, isSubmit bool /*是否提交，true为提交，false为暂存*/) (string, error) {
	//api, err := cache.SubmitExamAnswerApi(exam.ClazzId, exam.CourseId, exam.Paper.TestPaperId, exam.Paper.TestUserRelationId, exam.Cpi, exam.Paper.RemainTime, exam.Paper.EncRemainTime, exam.Paper.EncLastUpdateTime, exam.ExamRelationId, exam.AnswerId, exam.RemainTime, !isSubmit, exam.Paper.Enc, exam.Paper.EnterPageTime, exam.Paper.XXTExamQuestion.QuestionId, exam.Paper.Type, exam.Paper.XXTExamQuestion.TypeName, &exam.Paper)
	api, err := cache.SubmitWorkAnswerApi(&question.XXTWorkQuestionSubmitEntity, !isSubmit, 3, nil)
	if err != nil {
		return "", err
	}
	//fmt.Println(api)
	return api, nil
}

// html转Exam实体
func HtmlWorkQuestionTurnEntity(paperHtml string) (XXTWorkQuestion, error) {
	//xxtExamPaper := XXTExamPaper{}
	question := XXTWorkQuestion{}
	// 使用 goquery 解析 HTML
	paperDoc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(paperHtml)))
	if err != nil {
		return question, fmt.Errorf("作业题目页面解析失败: %v", err)
	}
	questionId, exists := paperDoc.Find("#questionId").Attr("value") //题目id
	if exists {
		//fmt.Println("question:", questionId)
		log2.Print(log2.DEBUG, questionId)
	}
	questionTypeCode, exists := paperDoc.Find(`input[name="` + `type` + questionId + `"]`).Attr("value")
	if exists {
		question.QuestionTypeCode = questionTypeCode
		log2.Print(log2.DEBUG, questionTypeCode)
		//fmt.Println("questionType:", questionTypeCode)
	}
	questionTypeStr, _ := paperDoc.Find(`span.focusSpan`).Attr("aria-label")
	// 去除前置题序，如 "1. "、"12. "
	re := regexp.MustCompile(`^\s*\d+\.\s*`)
	questionTypeStr = re.ReplaceAllString(questionTypeStr, "")

	get := func(id string) string {
		v, _ := paperDoc.Find("#" + id).Attr("value")
		return v
	}
	getForName := func(id string) string {
		v, _ := paperDoc.Find(`input[name="` + id + `"]`).Attr("value")
		return v
	}
	switch questionTypeCode {
	case "0": //单选题
		turn, err1 := singleWorkTurn(paperDoc)
		if err1 != nil {
			fmt.Println(err1)
		}
		question = turn
		//fmt.Println(turn)
	case "1": //多选题
		turn, err1 := multipleWorkTurn(paperDoc)
		if err1 != nil {
			fmt.Println(err1)
		}
		question = turn
		//fmt.Println(turn)
	case "2": //填空题
		turn, err1 := fillWorkTurn(paperDoc)
		if err1 != nil {
			fmt.Println(err1)
		}
		question = turn
		//fmt.Println(turn)
	case "3": //判断题
		turn, err1 := trueOrFalseWorkTurn(paperDoc)
		if err1 != nil {
			fmt.Println(err1)
		}
		question = turn
		//fmt.Println(turn)
	case "4": //简答题
		turn, err1 := shortAnswerWorkTurn(paperDoc)
		if err1 != nil {
			fmt.Println(err1)
		}
		question = turn
		//fmt.Println(turn)
	case "6": //论述题
		turn, err1 := essayWorkTurn(paperDoc)
		if err1 != nil {
			fmt.Println(err1)
		}
		question = turn
		//fmt.Println(turn)
	case "11": //连线题
		turn, err1 := matchingWorkTurn(paperDoc)
		if err1 != nil {
			fmt.Println(err1)
		}
		question = turn
	default:
		fmt.Printf("未知题型%s   %s", questionTypeCode, questionTypeStr)
	}
	question.CourseId = get("courseId")
	question.TestUserRelationId = get("testUserRelationId")
	question.ClassId = get("classId")
	question.Type = get("type")
	question.IsPhone = get("isphone")
	question.Imei = get("imei")
	question.SubCount = get("subCount")
	question.RemainTime = get("remainTime")
	question.TempSave = get("tempSave")
	question.TimeOver = get("timeOver")
	question.EncRemainTime = get("encRemainTime")
	question.EncLastUpdateTime = get("encLastUpdateTime")
	question.Cpi = get("cpi")
	question.Enc = get("enc")
	question.Source = get("source")
	question.Score = getForName("score" + questionId)
	question.UserId = get("userId")
	question.EnterPageTime = get("enterPageTime")
	question.AnsweredView = get("answeredView")
	question.AnsweredView = get("answeredView")
	question.PaperGroupId = get("paperGroupId")
	question.WordId = get("workId")
	question.CurrentTime = get("currentTime")
	question.CurrentCpi = get("currentCpi")
	question.CurrentUploadEnc = get("currentUploadEnc")
	question.MatchEnc = get("matchEnc")
	question.Cfid = get("cfid")
	question.AddTimes = get("addTimes")
	question.LimitWorkSubmitTimes = get("limitWorkSubmitTimes")
	question.QuestionTypeCode = questionTypeCode
	question.QuestionTypeStr = questionTypeStr
	question.EncWork = get("encWork")
	question.Index = get("index")
	return question, nil
}

// 提取作业题干（避免把选项 .workWrap 误当题干）
func extractWorkTitle(sel *goquery.Selection) string {
	titleSel := sel.Find(`div.ans-cc.timuStyle.fontLabel.workWrap`).First()
	if titleSel.Length() == 0 {
		titleSel = sel.Find(`div.ans-cc.workWrap`).First()
	}
	if titleSel.Length() == 0 {
		titleSel = sel.Find(`.workWrap`).First()
	}
	if titleSel.Length() == 0 {
		return ""
	}
	return extractQuestion(titleSel)
}

// 提取选项文本（兼容 div.centerSpan / p 文本）
func extractWorkOptionText(sel *goquery.Selection) string {
	text := strings.TrimSpace(sel.Text())
	if text == "" {
		text = strings.TrimSpace(sel.Find("p").Text())
	}
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	return strings.TrimSpace(text)
}

func extractWorkImages(sel *goquery.Selection) []qentity.QuestionImage {
	var images []qentity.QuestionImage
	seen := make(map[string]bool)
	sel.Find("img").Each(func(i int, img *goquery.Selection) {
		src, exists := img.Attr("src")
		if !exists || src == "" {
			return
		}
		if strings.HasPrefix(src, "//") {
			src = "https:" + src
		} else if strings.HasPrefix(src, "/") {
			src = "https://mooc1.chaoxing.com" + src
		}
		if seen[src] {
			return
		}
		seen[src] = true
		alt, _ := img.Attr("alt")
		images = append(images, qentity.QuestionImage{Url: src, AltText: alt})
	})
	return images
}

func singleWorkTurn(paperDoc *goquery.Document) (XXTWorkQuestion, error) {
	question := XXTWorkQuestion{}
	question.QType = qtype.SingleChoice
	question.Question.Type = question.QType.String()
	paperDoc.Find("div.singleQuesId").Each(func(i int, sel *goquery.Selection) {
		questionId, exists := sel.Attr("data")
		if exists {
			question.QuestionId = questionId
		}

		titleSel := sel.Find(`div.ans-cc.timuStyle.fontLabel.workWrap`).First()
		if titleSel.Length() == 0 {
			titleSel = sel.Find(`div.ans-cc.workWrap`).First()
		}
		if titleSel.Length() == 0 {
			titleSel = sel.Find(`.workWrap`).First()
		}
		if titleSel.Length() > 0 {
			question.Question.Content = extractQuestion(titleSel)
			question.Question.ContentImages = extractWorkImages(titleSel)
		}

		resOptions := make(map[string]string)
		question.Question.OptionImages = make(map[string][]qentity.QuestionImage)
		sel.Find(`div.centerSpan`).Each(func(i int, optSel *goquery.Selection) {
				letter, _ := optSel.Attr("id")
				text := extractWorkOptionText(optSel)
				imgs := extractWorkImages(optSel)
				if letter == "" || (text == "" && len(imgs) == 0) {
					return
				}
				resOptions[letter] = text
				if len(imgs) > 0 {
					question.Question.OptionImages[letter] = imgs
				}
			})
			resSelect := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N"}
			for _, slt := range resSelect {
				v, ok := resOptions[slt]
				if !ok {
					break
				}
				question.Question.Options = append(question.Question.Options, v)
			}
	})
	return question, nil
}

// matchingWorkTurn 连线题转换
func matchingWorkTurn(paperDoc *goquery.Document) (XXTWorkQuestion, error) {
	question := XXTWorkQuestion{}
	question.QType = qtype.Matching
	question.Question.Type = question.QType.String()
	paperDoc.Find("div.singleQuesId").Each(func(i int, sel *goquery.Selection) {
		questionId, exists := sel.Attr("data")
		if exists {
			question.QuestionId = questionId
		}
		titleSel := sel.Find(`div.ans-cc.timuStyle.fontLabel.workWrap`).First()
		if titleSel.Length() == 0 {
			titleSel = sel.Find(`div.ans-cc.workWrap`).First()
		}
		if titleSel.Length() == 0 {
			titleSel = sel.Find(`.workWrap`).First()
		}
		if titleSel.Length() > 0 {
			question.Question.Content = extractQuestion(titleSel)
			question.Question.ContentImages = extractWorkImages(titleSel)
		}
		options := []string{}
		selects := []string{}
		sel.Find("ul.answerList-line").Each(func(i int, s *goquery.Selection) {
			if i%2 == 0 {
				s.Find("li").Each(func(_ int, li *goquery.Selection) {
					options = append(options, strings.TrimSpace(li.Text()))
				})
			} else {
				s.Find("li").Each(func(_ int, li *goquery.Selection) {
					selects = append(selects, strings.TrimSpace(li.Text()))
				})
			}
		})
		question.Question.Options = buildMatchingStandardOptions(options, selects)
	})
	return question, nil
}

func multipleWorkTurn(paperDoc *goquery.Document) (XXTWorkQuestion, error) {
	question := XXTWorkQuestion{}
	question.QType = qtype.MultipleChoice
	question.Question.Type = question.QType.String()
	paperDoc.Find("div.singleQuesId").Each(func(i int, sel *goquery.Selection) {
		questionId, exists := sel.Attr("data")
		if exists {
			question.QuestionId = questionId
		}

		titleSel := sel.Find(`div.ans-cc.timuStyle.fontLabel.workWrap`).First()
		if titleSel.Length() == 0 {
			titleSel = sel.Find(`div.ans-cc.workWrap`).First()
		}
		if titleSel.Length() == 0 {
			titleSel = sel.Find(`.workWrap`).First()
		}
		if titleSel.Length() > 0 {
			question.Question.Content = extractQuestion(titleSel)
			question.Question.ContentImages = extractWorkImages(titleSel)
		}

		resOptions := make(map[string]string)
		question.Question.OptionImages = make(map[string][]qentity.QuestionImage)
		sel.Find(`div.centerSpan`).Each(func(i int, optSel *goquery.Selection) {
				letter, _ := optSel.Attr("id")
				text := extractWorkOptionText(optSel)
				imgs := extractWorkImages(optSel)
				if letter == "" || (text == "" && len(imgs) == 0) {
					return
				}
				resOptions[letter] = text
				if len(imgs) > 0 {
					question.Question.OptionImages[letter] = imgs
				}
			})
			resSelect := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N"}
			for _, slt := range resSelect {
				v, ok := resOptions[slt]
				if !ok {
					break
				}
				question.Question.Options = append(question.Question.Options, v)
			}
	})

	return question, nil
}

func fillWorkTurn(paperDoc *goquery.Document) (XXTWorkQuestion, error) {
	question := XXTWorkQuestion{}
	question.QType = qtype.FillInTheBlank
	question.Question.Type = question.QType.String()
	paperDoc.Find("div.singleQuesId").Each(func(i int, sel *goquery.Selection) {
		questionId, exists := sel.Attr("data")
		if exists {
			question.QuestionId = questionId
		}
		typeName, exists := paperDoc.Find(`input[name="` + `typeName` + questionId + `"]`).Attr("value")
		if exists {
			question.TypeName = typeName
		}

		title := ""
		sel.Find(`.workWrap`).Each(func(i int, sel *goquery.Selection) {
			title = extractQuestion(sel)
			question.Question.ContentImages = extractWorkImages(sel)
		})
		question.Question.Content = title
		sel.Find(`.blankItemName`).Each(func(i int, sel *goquery.Selection) {
			letter, _ := sel.Attr(`aria-label`)
			question.Question.Options = append(question.Question.Options, letter)
		})

	})
	return question, nil
}

func trueOrFalseWorkTurn(paperDoc *goquery.Document) (XXTWorkQuestion, error) {
	question := XXTWorkQuestion{}
	question.QType = qtype.TrueOrFalse
	question.Question.Type = question.QType.String()
	paperDoc.Find("div.singleQuesId").Each(func(i int, sel *goquery.Selection) {
		questionId, exists := sel.Attr("data")
		if exists {
			question.QuestionId = questionId
		}
		typeName, exists := paperDoc.Find(`input[name="` + `typeName` + questionId + `"]`).Attr("value")
		if exists {
			question.TypeName = typeName
		}

		title := ""
		sel.Find(`.workWrap`).Each(func(i int, sel *goquery.Selection) {
			title = extractQuestion(sel)
			question.Question.ContentImages = extractWorkImages(sel)
		})
		question.Question.Content = title
		sel.Find(`.Answer`).Each(func(i int, sel *goquery.Selection) {
			letter := sel.Find("span.check").Text()
			text := sel.Find("p.fl").Text()
			question.Question.Options = append(question.Question.Options, letter+text)
		})

	})
	return question, nil
}

func shortAnswerWorkTurn(paperDoc *goquery.Document) (XXTWorkQuestion, error) {
	question := XXTWorkQuestion{}
	question.QType = qtype.ShortAnswer
	question.Question.Type = question.QType.String()
	paperDoc.Find("div.singleQuesId").Each(func(i int, sel *goquery.Selection) {
		questionId, exists := sel.Attr("data")
		if exists {
			question.QuestionId = questionId
		}
		typeName, exists := paperDoc.Find(`input[name="` + `typeName` + questionId + `"]`).Attr("value")
		if exists {
			question.TypeName = typeName
		}

		title := ""
		sel.Find(`.workWrap`).Each(func(i int, sel *goquery.Selection) {
			title = extractQuestion(sel)
			question.Question.ContentImages = extractWorkImages(sel)
		})
		question.Question.Content = title

	})
	return question, nil
}

func essayWorkTurn(paperDoc *goquery.Document) (XXTWorkQuestion, error) {
	question := XXTWorkQuestion{}
	question.QType = qtype.Essay
	question.Question.Type = question.QType.String()
	paperDoc.Find("div.singleQuesId").Each(func(i int, sel *goquery.Selection) {
		questionId, exists := sel.Attr("data")
		if exists {
			question.QuestionId = questionId
		}
		typeName, exists := paperDoc.Find(`input[name="` + `typeName` + questionId + `"]`).Attr("value")
		if exists {
			question.TypeName = typeName
		}

		title := ""
		sel.Find(`.workWrap`).Each(func(i int, sel *goquery.Selection) {
			title = extractQuestion(sel)
			question.Question.ContentImages = extractWorkImages(sel)
		})
		question.Question.Content = title

	})
	return question, nil
}
