package xuexitong

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/yatori-dev/yatori-go-core/utils"
)

// ==================== 学习通互动测验（视频弹题）API ====================

func (cache *XueXiTUserCache) quizGet(rawURL string, referer string, retry int, lastErr error) (string, error) {
	if retry < 0 {
		return "", lastErr
	}
	tr := &http.Transport{}
	if cache.IpProxySW {
		tr.Proxy = func(req *http.Request) (*url.URL, error) {
			return url.Parse("http://" + cache.ProxyIP)
		}
	}
	client := &http.Client{Transport: tr, Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("User-Agent", GetUA("web"))
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Referer", referer)
	for _, cookie := range cache.cookies {
		req.AddCookie(cookie)
	}
	res, err := client.Do(req)
	if err != nil {
		time.Sleep(time.Duration(3-retry) * time.Second)
		return cache.quizGet(rawURL, referer, retry-1, err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	utils.CookiesAddNoRepetition(&cache.cookies, res.Cookies())
	return string(body), nil
}

// FetchQuizTimelineApi 拉取视频弹题时间轴
func (cache *XueXiTUserCache) FetchQuizTimelineApi(mid string, cpi, classID, courseID int, retry int, lastErr error) (string, error) {
	params := url.Values{}
	params.Add("mid", mid)
	params.Add("cpi", strconv.Itoa(cpi))
	params.Add("classid", strconv.Itoa(classID))
	params.Add("courseid", strconv.Itoa(courseID))
	u := "https://mooc1.chaoxing.com/mooc-ans/richvideo/initdatawithviewerV2?" + params.Encode()
	return cache.quizGet(u, "https://mooc1.chaoxing.com/", retry, lastErr)
}

// FetchQuizPlanApi 拉取互动测验权威清单（status 为空表示未完成）
func (cache *XueXiTUserCache) FetchQuizPlanApi(classID, courseID, cpi int, status string, pageNo, pageSize int, retry int, lastErr error) (string, error) {
	params := url.Values{}
	params.Add("clazzid", strconv.Itoa(classID))
	params.Add("courseid", strconv.Itoa(courseID))
	params.Add("cpi", strconv.Itoa(cpi))
	params.Add("ut", "s")
	params.Add("pageNo", strconv.Itoa(pageNo))
	params.Add("pageSize", strconv.Itoa(pageSize))
	params.Add("status", status)
	u := "https://stat2-ans.chaoxing.com/study-videoTest/finish?" + params.Encode()
	return cache.quizGet(u, "https://stat2-ans.chaoxing.com/study-videoTest/index", retry, lastErr)
}

// SubmitQuizAnswerApi 提交单个弹题答案。
// answer 格式："A"、"B" 或 "A,B,C"
func (cache *XueXiTUserCache) SubmitQuizAnswerApi(classID, cpi int, objectID, eventID, memberInfo, answer string, retry int, lastErr error) (string, error) {
	params := url.Values{}
	params.Add("classid", strconv.Itoa(classID))
	params.Add("cpi", strconv.Itoa(cpi))
	params.Add("objectid", objectID)
	params.Add("eventid", eventID)
	params.Add("memberinfo", memberInfo)
	params.Add("answerContent", answer)
	u := "https://mooc1.chaoxing.com/question/quiz-validation?" + params.Encode()
	body, err := cache.quizGet(u, "https://mooc1.chaoxing.com/ananas/modules/video/index.html", retry, lastErr)
	if err != nil {
		return "", err
	}
	if len(body) > 2 && body[:2] == "<!" {
		return body, errors.New("会话降级，返回HTML")
	}
	return body, nil
}

// FetchChapterCords2Raw 直接返回章节卡片页原始内容（用于解析视频 mid/objectId）
func (cache *XueXiTUserCache) FetchChapterCords2Raw(clazzid, courseid, knowledgeid, cpi string, retry int, lastErr error) (string, error) {
	u := fmt.Sprintf("https://mooc1.chaoxing.com/mooc-ans/knowledge/cards?clazzid=%s&courseid=%s&knowledgeid=%s&num=0&ut=s&cpi=%s&v=2025-0424-1038-3&mooc2=1&isMicroCourse=false&editorPreview=0",
		clazzid, courseid, knowledgeid, cpi)
	return cache.quizGet(u, "https://mooc1.chaoxing.com/mooc-ans/mycourse/studentstudy", retry, lastErr)
}
