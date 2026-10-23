package xuexitong

import (
	"encoding/json"
	"fmt"
	"html"
	"regexp"
	"strings"
	"time"

	xuexitongApi "github.com/yatori-dev/yatori-go-core/api/xuexitong"
)

// OuterInfo 课程外层页解析结果
type OuterInfo struct {
	BbsID  string
	Enc    string
	Fid    string
	UID    string
	Stuenc string
	Menus  map[string]string
}

var (
	reMenuItem = regexp.MustCompile(`<li[^>]*dataname="([^"]+)"[^>]*>\s*<a[^>]*data-url="([^"]+)"`)
	reMenuAny  = regexp.MustCompile(`dataname="([^"]+)"[^>]*data-url="([^"]+)"`)
	reStuenc   = regexp.MustCompile(`studentcourse\?[^"']*?stuenc=([0-9a-f]{32})`)
)

func htmlField(h, id string) string {
	re := regexp.MustCompile(`id="` + regexp.QuoteMeta(id) + `"[^>]*value="([^"]*)"`)
	if m := re.FindStringSubmatch(h); len(m) > 1 {
		return m[1]
	}
	return ""
}

// ParseOuterInfo 解析课程外层页菜单与讨论/阅读所需参数
func ParseOuterInfo(page string) OuterInfo {
	info := OuterInfo{Menus: map[string]string{}}
	for _, m := range reMenuItem.FindAllStringSubmatch(page, -1) {
		info.Menus[m[1]] = m[2]
	}
	for _, m := range reMenuAny.FindAllStringSubmatch(page, -1) {
		if _, ok := info.Menus[m[1]]; !ok {
			info.Menus[m[1]] = m[2]
		}
	}
	info.BbsID = htmlField(page, "bbsid")
	info.Enc = htmlField(page, "enc")
	if info.Enc == "" {
		info.Enc = htmlField(page, "showAIEnc")
	}
	info.Fid = htmlField(page, "cfid")
	if info.Fid == "" {
		info.Fid = htmlField(page, "fid")
	}
	info.UID = htmlField(page, "uid")
	if info.UID == "" {
		if m := regexp.MustCompile(`(?:^|;\s*)_uid=([^;]+)`).FindStringSubmatch(page); len(m) > 1 {
			info.UID = m[1]
		}
	}
	if m := reStuenc.FindStringSubmatch(page); len(m) > 1 {
		info.Stuenc = m[1]
	}
	return info
}

// ReadBook 拓展阅读书目
type ReadBook struct {
	ID    string
	Title string
}

// ParseReadCategory 从阅读页提取分类 id
func ParseReadCategory(page string) string {
	if m := regexp.MustCompile(`resource/(\d+)/list`).FindStringSubmatch(page); len(m) > 1 {
		return m[1]
	}
	return "12576"
}

// ParseReadBooks 解析书目列表
func ParseReadBooks(body string) []ReadBook {
	var root map[string]interface{}
	if err := json.Unmarshal([]byte(body), &root); err != nil {
		return nil
	}
	list, _ := root["data"].([]interface{})
	var books []ReadBook
	for _, it := range list {
		m, ok := it.(map[string]interface{})
		if !ok {
			continue
		}
		var ri map[string]interface{}
		switch v := m["resourceInfo"].(type) {
		case string:
			_ = json.Unmarshal([]byte(v), &ri)
		case map[string]interface{}:
			ri = v
		}
		if ri == nil {
			continue
		}
		id := toStr(ri["id"])
		if id == "" {
			continue
		}
		title := ""
		if inf, ok := ri["info"].(map[string]interface{}); ok {
			title = toStr(inf["field1"])
			if title == "" {
				title = toStr(inf["author"])
			}
		}
		if title == "" {
			title = toStr(m["resourceId"])
		}
		if len(title) > 16 {
			title = title[:16]
		}
		books = append(books, ReadBook{ID: id, Title: title})
	}
	return books
}

// RunExtendReadAction 拓展阅读纯协议流程，返回成功上报的书目数
func RunExtendReadAction(cache *xuexitongApi.XueXiTUserCache, courseID, classID, cpi string, maxBooks, pagesPerBook int) (int, string, error) {
	page, err := cache.FetchStuCourseMiddle(courseID, classID, cpi, 3, nil)
	if err != nil {
		return 0, "", err
	}
	info := ParseOuterInfo(page)
	if _, ok := info.Menus["expandread"]; !ok {
		return 0, "该课程无「拓展阅读」菜单", nil
	}
	if info.UID == "" {
		info.UID = strings.TrimSpace(cpi)
	}
	index, err := cache.FetchExtendReadIndex(courseID, classID, cpi, info.Stuenc, info.Fid, info.UID, 3, nil)
	if err != nil {
		return 0, "", err
	}
	if len(index) < 500 {
		return 0, "进入拓展阅读页失败", nil
	}
	cat := ParseReadCategory(index)
	list, err := cache.FetchExtendReadBooks(cat, 3, nil)
	if err != nil {
		return 0, "", err
	}
	books := ParseReadBooks(list)
	if len(books) == 0 {
		return 0, "书目列表为空", nil
	}
	if maxBooks > 0 && len(books) > maxBooks {
		books = books[:maxBooks]
	}
	ok := 0
	for _, b := range books {
		pages := 0
		for pg := 1; pg <= pagesPerBook; pg++ {
			body, err := cache.SubmitExtendReadDot(info.UID, b.ID, pg, 2, nil)
			if err == nil && strings.Contains(strings.ReplaceAll(body, " ", ""), "\"code\":1") {
				pages++
			}
			time.Sleep(800 * time.Millisecond)
		}
		if pages > 0 {
			ok++
		}
	}
	return ok, fmt.Sprintf("共 %d 本书，成功上报 %d 本", len(books), ok), nil
}

// Topic 讨论话题
type Topic struct {
	ID        string
	UUID      string
	CreaterID string
	Content   string
}

// ParseTopics 解析 getTopicList 返回
func ParseTopics(body string) []Topic {
	var root map[string]interface{}
	if err := json.Unmarshal([]byte(body), &root); err != nil {
		return nil
	}
	list, _ := root["datas"].([]interface{})
	var topics []Topic
	for _, it := range list {
		m, ok := it.(map[string]interface{})
		if !ok {
			continue
		}
		t := Topic{
			ID:        toStr(m["id"]),
			UUID:      toStr(m["uuid"]),
			CreaterID: toStr(m["createrId"]),
			Content:   stripHTML(toStr(m["content"])),
		}
		if t.UUID != "" {
			topics = append(topics, t)
		}
	}
	return topics
}

func stripHTML(s string) string {
	s = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(s, " ")
	s = html.UnescapeString(s)
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	if len(s) > 120 {
		s = s[:120]
	}
	return s
}

// FetchAllTopicsAction 拉取课程讨论全部话题（最多 5 页）
func FetchAllTopicsAction(cache *xuexitongApi.XueXiTUserCache, courseID, classID, cpi string) ([]Topic, string, error) {
	page, err := cache.FetchStuCourseMiddle(courseID, classID, cpi, 3, nil)
	if err != nil {
		return nil, "", err
	}
	info := ParseOuterInfo(page)
	if info.BbsID == "" {
		return nil, "该课程无讨论区", nil
	}
	_, _ = cache.FetchTopicListPage(courseID, classID, cpi, info.Stuenc, info.BbsID, 2, nil)
	var all []Topic
	seen := map[string]bool{}
	for p := 1; p <= 5; p++ {
		body, err := cache.FetchTopics(info.BbsID, courseID, classID, p, 2, nil)
		if err != nil {
			break
		}
		list := ParseTopics(body)
		if len(list) == 0 {
			break
		}
		for _, t := range list {
			if !seen[t.UUID] {
				seen[t.UUID] = true
				all = append(all, t)
			}
		}
		if len(list) < 20 {
			break
		}
		time.Sleep(time.Second)
	}
	return all, info.BbsID, nil
}

// ReplyTopicAction 回复单个话题
func ReplyTopicAction(cache *xuexitongApi.XueXiTUserCache, bbsid string, topic Topic, content, courseID, classID, cpi string) (bool, error) {
	body, err := cache.AddTopicReply(bbsid, topic.ID, topic.UUID, content, courseID, classID, cpi, 2, nil)
	if err != nil {
		return false, err
	}
	if !strings.HasPrefix(strings.TrimSpace(body), "{") {
		return false, fmt.Errorf("接口返回非 JSON，可能需要浏览器会话")
	}
	var root map[string]interface{}
	if err := json.Unmarshal([]byte(body), &root); err != nil {
		return false, err
	}
	if toStr(root["status"]) == "1" || toStr(root["msg"]) == "success" {
		return true, nil
	}
	return false, fmt.Errorf("回复未生效: %s", toStr(root["msg"]))
}
