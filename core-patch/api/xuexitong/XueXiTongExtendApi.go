package xuexitong

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/yatori-dev/yatori-go-core/utils"
)

func (cache *XueXiTUserCache) extPost(rawURL, referer string, form url.Values, retry int, lastErr error) (string, error) {
	if retry < 0 {
		return "", lastErr
	}
	tr := &http.Transport{}
	if cache.IpProxySW {
		tr.Proxy = func(req *http.Request) (*url.URL, error) {
			return url.Parse("http://" + cache.ProxyIP)
		}
	}
	client := &http.Client{Transport: tr, Timeout: 25 * time.Second}
	req, err := http.NewRequest("POST", rawURL, nil)
	if err != nil {
		return "", err
	}
	req.PostForm = form
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", GetUA("web"))
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	if referer != "" {
		req.Header.Set("Referer", referer)
	}
	for _, c := range cache.cookies {
		req.AddCookie(c)
	}
	res, err := client.Do(req)
	if err != nil {
		return cache.extPost(rawURL, referer, form, retry-1, err)
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	utils.CookiesAddNoRepetition(&cache.cookies, res.Cookies())
	return string(raw), nil
}

// FetchStuCourseMiddle 拉取课程外层页（含菜单、bbsid、stuenc）
func (cache *XueXiTUserCache) FetchStuCourseMiddle(courseID, classID, cpi string, retry int, lastErr error) (string, error) {
	u := fmt.Sprintf("https://mooc1-1.chaoxing.com/mooc-ans/visit/stucoursemiddle?courseid=%s&clazzid=%s&vc=1&cpi=%s&ismooc2=1&v=2",
		courseID, classID, cpi)
	return cache.quizGet(u, "https://mooc1.chaoxing.com/", retry, lastErr)
}

// FetchExtendReadIndex 进入拓展阅读页
func (cache *XueXiTUserCache) FetchExtendReadIndex(courseID, classID, cpi, stuenc, fid, uid string, retry int, lastErr error) (string, error) {
	u := fmt.Sprintf("https://appcd.chaoxing.com/read/index?courseid=%s&clazzid=%s&cpi=%s&ut=s&t=%d&stuenc=%s&fid=%s&userId=%s",
		courseID, classID, cpi, time.Now().UnixMilli(), stuenc, fid, uid)
	return cache.quizGet(u, "https://mooc1.chaoxing.com/", retry, lastErr)
}

// FetchExtendReadBooks 拉取拓展阅读书目列表
func (cache *XueXiTUserCache) FetchExtendReadBooks(categoryID string, retry int, lastErr error) (string, error) {
	u := fmt.Sprintf("https://appcd.chaoxing.com/read/resource/%s/list?page=1&pageSize=20", categoryID)
	return cache.quizGet(u, "https://appcd.chaoxing.com/read/index", retry, lastErr)
}

// SubmitExtendReadDot 上报阅读点
func (cache *XueXiTUserCache) SubmitExtendReadDot(uid, resourceID string, page int, retry int, lastErr error) (string, error) {
	u := fmt.Sprintf("https://appcd.chaoxing.com/read/api/dot/%s/%s_%d", uid, resourceID, page)
	return cache.extPost(u, "https://appcd.chaoxing.com/read/index", url.Values{}, retry, lastErr)
}

// FetchTopicListPage 打开发起讨论会话
func (cache *XueXiTUserCache) FetchTopicListPage(courseID, classID, cpi, stuenc, bbsid string, retry int, lastErr error) (string, error) {
	u := fmt.Sprintf("https://groupweb.chaoxing.com/course/topic/topicList?courseid=%s&clazzid=%s&cpi=%s&ut=s&t=%d&stuenc=%s&bbsid=%s&learnSilverStartTime=&learnSilverEndTime=",
		courseID, classID, cpi, time.Now().UnixMilli(), stuenc, bbsid)
	return cache.quizGet(u, "https://mooc1.chaoxing.com/", retry, lastErr)
}

// FetchTopics 拉取课程讨论话题列表
func (cache *XueXiTUserCache) FetchTopics(bbsid, courseID, classID string, page int, retry int, lastErr error) (string, error) {
	u := fmt.Sprintf("https://groupweb.chaoxing.com/course/topic/%s/getTopicList?folder_uuid=&page=%d&pageSize=20&kw=&courseId=%s&isSetTop=&selectedStartTime=&selectedEndTime=&selectedType=0&learnSilverStartTime=&learnSilverEndTime=",
		bbsid, page, courseID)
	referer := fmt.Sprintf("https://groupweb.chaoxing.com/course/topic/topicList?courseid=%s&clazzid=%s", courseID, classID)
	form := url.Values{}
	form.Set("tags", "classId0000001,classId"+classID+",courseId"+courseID)
	return cache.extPost(u, referer, form, retry, lastErr)
}

// AddTopicReply 回复讨论话题
func (cache *XueXiTUserCache) AddTopicReply(bbsid, topicID, uuid, content, courseID, classID, cpi string, retry int, lastErr error) (string, error) {
	form := url.Values{}
	form.Set("bbsid", bbsid)
	form.Set("topicId", topicID)
	form.Set("uuid", uuid)
	form.Set("content", content)
	form.Set("courseId", courseID)
	form.Set("classId", classID)
	form.Set("cpi", cpi)
	referer := fmt.Sprintf("https://groupweb.chaoxing.com/course/topic/topicList?courseid=%s&clazzid=%s&cpi=%s", courseID, classID, cpi)
	u := fmt.Sprintf("https://groupweb.chaoxing.com/course/topic/%s/addReply", bbsid)
	return cache.extPost(u, referer, form, retry, lastErr)
}
