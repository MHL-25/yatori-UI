package aiq

import (
	"encoding/json"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/yatori-dev/yatori-go-core/que-core/qentity"
	"github.com/yatori-dev/yatori-go-core/que-core/qtype"
)

// ParseStringArray 宽松解析 AI 回复，兼容 JSON 数组、JSON 对象、代码块和纯文本
func ParseStringArray(raw string) []string {
	s := normalizeJSONAnswerContent(raw)
	if s == "" {
		return nil
	}
	if value, ok := decodeFirstJSONValue(s); ok {
		if answers := collectAnswerStrings(value, 0); len(answers) > 0 {
			return uniqueNonEmpty(answers)
		}
	}
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = cleanAnswerLine(line)
		if line != "" {
			out = append(out, line)
		}
	}
	if len(out) == 0 {
		out = append(out, strings.TrimSpace(s))
	}
	return uniqueNonEmpty(out)
}

func decodeFirstJSONValue(s string) (interface{}, bool) {
	start := -1
	for i, r := range s {
		if r == '[' || r == '{' || r == '"' {
			start = i
			break
		}
	}
	if start == -1 {
		return nil, false
	}
	decoder := json.NewDecoder(strings.NewReader(s[start:]))
	decoder.UseNumber()
	var value interface{}
	if err := decoder.Decode(&value); err != nil {
		return nil, false
	}
	return value, true
}

func collectAnswerStrings(value interface{}, depth int) []string {
	if depth > 8 {
		return nil
	}
	switch v := value.(type) {
	case string:
		cleaned := cleanAnswerLine(v)
		if cleaned != "" {
			return []string{cleaned}
		}
	case json.Number:
		return []string{v.String()}
	case float64:
		return []string{strconv.FormatFloat(v, 'f', -1, 64)}
	case bool:
		if v {
			return []string{"true"}
		}
		return []string{"false"}
	case []interface{}:
		out := make([]string, 0, len(v))
		for _, item := range v {
			out = append(out, collectAnswerStrings(item, depth+1)...)
		}
		return out
	case map[string]interface{}:
		priorityKeys := []string{
			"answer", "answers", "correct_answer", "correctAnswer",
			"content", "text", "result", "data", "value", "option", "color",
		}
		for _, key := range priorityKeys {
			if child, ok := v[key]; ok {
				if answers := collectAnswerStrings(child, depth+1); len(answers) > 0 {
					return answers
				}
			}
		}
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		out := make([]string, 0, len(v))
		for _, key := range keys {
			out = append(out, collectAnswerStrings(v[key], depth+1)...)
		}
		return out
	}
	return nil
}

func cleanAnswerLine(line string) string {
	line = strings.TrimSpace(line)
	line = strings.Trim(line, "`\"'")
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "- ")
	line = strings.TrimPrefix(line, "* ")
	lower := strings.ToLower(line)
	for _, prefix := range []string{
		"正确答案是：", "正确答案是:", "正确答案是",
		"答案是：", "答案是:", "答案是",
		"答案：", "答案:", "正确答案：", "正确答案:",
		"answer:", "answer：", "选：", "选:",
	} {
		if strings.HasPrefix(lower, strings.ToLower(prefix)) {
			line = strings.TrimSpace(line[len(prefix):])
			lower = strings.ToLower(line)
			break
		}
	}
	return strings.TrimSpace(line)
}

func uniqueNonEmpty(values []string) []string {
	seen := make(map[string]bool)
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	return out
}

// ParseQuestionAnswers 按题型解析 AI 回复：填空题严格按空位数量，其他题型宽松解析
func ParseQuestionAnswers(question qentity.Question, raw string) []string {
	if question.Type == qtype.FillInTheBlank.String() {
		return ParseFillAnswers(raw, len(question.Options))
	}
	return ParseStringArray(raw)
}

var (
	fillMetaPattern   = regexp.MustCompile(`[（(][^）)]*(?:此处|应绘制|应提供|见图|如图|略)[^）)]*[）)]`)
	fillNumberPattern = regexp.MustCompile(`^\s*(?:第\s*\d+\s*空\s*[:：]?|\(?\d+\)?\s*[、，)）:：]|[-*])\s*`)
	fillDotPattern    = regexp.MustCompile(`^\s*\(?\d+\)?\s*\.\s*`)
)

// ParseFillAnswers 解析填空题答案：保留顺序和重复项，并严格校验空位数量。
// 数量与题目空位不一致时返回 nil，避免把开场白或说明文字当成答案提交。
func ParseFillAnswers(raw string, blankCount int) []string {
	if blankCount <= 0 {
		return nil
	}
	s := normalizeJSONAnswerContent(raw)
	if s == "" {
		return nil
	}
	answers := make([]string, 0, blankCount)
	if value, ok := decodeFirstJSONValue(s); ok {
		if arr, ok := value.([]interface{}); ok {
			for _, item := range arr {
				if text := cleanFillAnswer(ContentToString(item)); text != "" {
					answers = append(answers, text)
				}
			}
		} else if text := cleanFillAnswer(ContentToString(value)); text != "" {
			answers = append(answers, text)
		}
	}
	if len(answers) == 0 {
		for _, line := range strings.Split(s, "\n") {
			line = cleanFillAnswer(line)
			if line == "" || isFillNoise(line) {
				continue
			}
			answers = append(answers, line)
		}
	}
	if len(answers) != blankCount {
		return nil
	}
	return answers
}

func cleanFillAnswer(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "**", "")
	s = strings.ReplaceAll(s, "__", "")
	s = fillMetaPattern.ReplaceAllString(s, "")
	s = fillNumberPattern.ReplaceAllString(s, "")
	if loc := fillDotPattern.FindStringIndex(s); loc != nil && len(s) > loc[1] {
		// “1.形状”这类编号去掉；但“1.5”这种数字答案不能误伤
		rest := s[loc[1]:]
		if rest != "" && rest[0] >= '0' && rest[0] <= '9' {
			// 保留原值
		} else {
			s = rest
		}
	}
	s = strings.Trim(strings.TrimSpace(s), "`\"'")
	return strings.TrimSpace(s)
}

func isFillNoise(s string) bool {
	for _, keyword := range []string{
		"以下是", "答案如下", "根据提供", "不需要", "已经给出", "填空的答案", "故答案为", "综上所述",
	} {
		if strings.Contains(s, keyword) {
			return true
		}
	}
	return strings.HasSuffix(s, "：") || strings.HasSuffix(s, ":")
}
