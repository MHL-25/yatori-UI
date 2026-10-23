package xuexitong

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	xuexitongApi "github.com/yatori-dev/yatori-go-core/api/xuexitong"
)

// QuizPlanItem 互动测验未完成清单中的一条（以视频为单位）
type QuizPlanItem struct {
	KnowledgeID         int    `json:"knowledgeId"`
	Name                string `json:"name"`
	KnowledgeName       string `json:"knowledgeName"`
	QuestionCount       int    `json:"questionCount"`
	FinishQuestionCount int    `json:"finishQuestionCount"`
	Type                string `json:"type"`
	Rate                string `json:"rate"`
}

// QuizItem 单个弹题
type QuizItem struct {
	ResourceID    string                   `json:"resourceId"`
	Description   string                   `json:"description"`
	Options       []map[string]interface{} `json:"options"`
	QuestionType  string                   `json:"questionType"`
	MemberInfo    string                   `json:"memberinfo"`
	AnswerContent string                   `json:"answerContent"`
	RawOptions    []string                 `json:"-"`
}

// QuizVideo 视频与其弹题
type QuizVideo struct {
	KnowledgeID int
	ObjectID    string
	Mid         string
	Quizzes     []QuizItem
}

func toStr(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		if t == float64(int64(t)) {
			return strconv.FormatInt(int64(t), 10)
		}
		return strconv.FormatFloat(t, 'f', -1, 64)
	case json.Number:
		return t.String()
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", t)
	}
}

// PullQuizPlanAction 拉取互动测验未完成清单
func PullQuizPlanAction(cache *xuexitongApi.XueXiTUserCache, classID, courseID, cpi int) ([]QuizPlanItem, error) {
	var all []QuizPlanItem
	for page := 1; page <= 10; page++ {
		body, err := cache.FetchQuizPlanApi(classID, courseID, cpi, "", page, 16, 3, nil)
		if err != nil {
			if len(all) > 0 {
				break
			}
			return nil, err
		}
		var root map[string]interface{}
		if err := json.Unmarshal([]byte(body), &root); err != nil {
			if len(all) > 0 {
				break
			}
			return nil, err
		}
		data, _ := root["data"].(map[string]interface{})
		if data == nil {
			break
		}
		list, _ := data["data"].([]interface{})
		if len(list) == 0 {
			break
		}
		for _, item := range list {
			m, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			row := QuizPlanItem{
				KnowledgeID:         int(toInt(m["knowledgeId"])),
				Name:                toStr(m["name"]),
				KnowledgeName:       toStr(m["knowledgeName"]),
				QuestionCount:       int(toInt(m["questionCount"])),
				FinishQuestionCount: int(toInt(m["finishQuestionCount"])),
				Type:                toStr(m["type"]),
				Rate:                toStr(m["rate"]),
			}
			if row.KnowledgeID > 0 {
				all = append(all, row)
			}
		}
		pageInfo, _ := data["pageInfo"].(map[string]interface{})
		totalPage := int(toInt(pageInfo["totalPage"]))
		if totalPage > 0 && page >= totalPage {
			break
		}
	}
	return all, nil
}

func toInt(v interface{}) int64 {
	switch t := v.(type) {
	case float64:
		return int64(t)
	case json.Number:
		n, _ := t.Int64()
		return n
	case string:
		n, _ := strconv.ParseInt(t, 10, 64)
		return n
	default:
		return 0
	}
}

// ParseQuizTimeline 解析 initdatawithviewerV2 返回的弹题数据
func ParseQuizTimeline(body string) []QuizItem {
	var v interface{}
	if err := json.Unmarshal([]byte(body), &v); err != nil {
		return nil
	}
	var result []QuizItem
	var consume func(item interface{})
	consume = func(item interface{}) {
		m, ok := item.(map[string]interface{})
		if !ok {
			return
		}
		if datas, ok := m["datas"].([]interface{}); ok {
			for _, d := range datas {
				consume(d)
			}
			return
		}
		if _, ok := m["resourceId"]; !ok {
			return
		}
		q := QuizItem{
			ResourceID:    toStr(m["resourceId"]),
			Description:   toStr(m["description"]),
			QuestionType:  toStr(m["questionType"]),
			MemberInfo:    toStr(m["memberinfo"]),
			AnswerContent: toStr(m["answerContent"]),
		}
		if opts, ok := m["options"].([]interface{}); ok {
			for _, o := range opts {
				if om, ok := o.(map[string]interface{}); ok {
					q.Options = append(q.Options, om)
					q.RawOptions = append(q.RawOptions, toStr(om["content"]))
				} else {
					q.RawOptions = append(q.RawOptions, toStr(o))
				}
			}
		}
		if q.ResourceID != "" && q.Description != "" {
			result = append(result, q)
		}
	}
	switch t := v.(type) {
	case []interface{}:
		for _, item := range t {
			consume(item)
		}
	case map[string]interface{}:
		for _, key := range []string{"datas", "data", "list"} {
			if list, ok := t[key].([]interface{}); ok {
				for _, item := range list {
					consume(item)
				}
			}
		}
	}
	return result
}

// ParseVideoAttachment 从章节卡片 mArg 中提取视频 mid 与 objectId
func ParseVideoAttachment(cardBody string) (mid, objectID string, ok bool) {
	mArg := extractMArgJSON(cardBody)
	if mArg == "" {
		return "", "", false
	}
	var root struct {
		Attachments []struct {
			Type     string `json:"type"`
			Mid      string `json:"mid"`
			ObjectID string `json:"objectId"`
			Property struct {
				Mid string `json:"mid"`
			} `json:"property"`
		} `json:"attachments"`
	}
	if err := json.Unmarshal([]byte(mArg), &root); err != nil {
		return "", "", false
	}
	for _, a := range root.Attachments {
		if a.Type != "video" {
			continue
		}
		mid = a.Mid
		if mid == "" {
			mid = a.Property.Mid
		}
		if mid != "" && a.ObjectID != "" {
			return mid, a.ObjectID, true
		}
	}
	return "", "", false
}

func extractMArgJSON(body string) string {
	key := "mArg = {"
	i := strings.Index(body, key)
	if i < 0 {
		key = "mArg={"
		i = strings.Index(body, key)
		if i < 0 {
			return ""
		}
	}
	start := i + len("mArg = ")
	if key == "mArg={" {
		start = i + len("mArg=")
	}
	depth := 0
	inStr := false
	esc := false
	for k := start; k < len(body); k++ {
		c := body[k]
		if inStr {
			if esc {
				esc = false
			} else if c == '\\' {
				esc = true
			} else if c == '"' {
				inStr = false
			}
			continue
		}
		switch c {
		case '"':
			inStr = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return body[start : k+1]
			}
		}
	}
	return ""
}

// BuildQuizPrompt 构建弹题提示词
func BuildQuizPrompt(q QuizItem) string {
	var sb strings.Builder
	sb.WriteString("请回答下面的题目，只输出选项字母。\n")
	if q.QuestionType != "" {
		sb.WriteString("题型：" + q.QuestionType + "\n")
	}
	sb.WriteString("题目：" + q.Description + "\n")
	if len(q.RawOptions) > 0 {
		sb.WriteString("选项：\n")
		for i, opt := range q.RawOptions {
			sb.WriteString(fmt.Sprintf("%c. %s\n", 'A'+i, opt))
		}
	}
	sb.WriteString("要求：单选题输出一个字母；多选题输出多个字母并用英文逗号分隔。不要输出任何解释。")
	return sb.String()
}

// NormalizeQuizAnswer 从 AI 返回内容中提取答案字母
func NormalizeQuizAnswer(text string) (string, error) {
	text = strings.ToUpper(strings.TrimSpace(text))
	var letters []string
	seen := map[rune]bool{}
	for _, r := range text {
		if r >= 'A' && r <= 'F' && !seen[r] {
			seen[r] = true
			letters = append(letters, string(r))
		}
	}
	if len(letters) == 0 {
		return "", errors.New("无法解析答案")
	}
	return strings.Join(letters, ","), nil
}
