package yinghua

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	yinghuaApi "github.com/yatori-dev/yatori-go-core/api/yinghua"
)

// NormalizeYingHuaURL 将用户填写的页面地址统一为站点根地址
// 例如 https://cdcas.suwankj.com/user/index -> https://cdcas.suwankj.com
func NormalizeYingHuaURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return raw
	}
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Host == "" {
		return strings.TrimRight(raw, "/")
	}
	return parsed.Scheme + "://" + parsed.Host
}

// IsYingHuaCookieOrToken 判断用户填写的是 Cookie 还是普通密码
func IsYingHuaCookieOrToken(value string) bool {
	value = strings.TrimSpace(value)
	return strings.Contains(value, "token=") || strings.HasPrefix(value, "sid.")
}

// YingHuaCookieLoginAction 使用浏览器 Cookie / token 登录英华
func YingHuaCookieLoginAction(cache *yinghuaApi.YingHuaUserCache) error {
	cache.PreUrl = NormalizeYingHuaURL(cache.PreUrl)
	raw := strings.TrimSpace(cache.Password)
	if raw == "" {
		return errors.New("Cookie 为空")
	}

	var token string
	var cookies []*http.Cookie
	if strings.Contains(raw, "=") {
		req := &http.Request{Header: http.Header{"Cookie": []string{raw}}}
		cookies = req.Cookies()
		for _, c := range cookies {
			if c.Name == "token" {
				token = c.Value
				break
			}
		}
	}
	if token == "" {
		token = strings.TrimSpace(strings.TrimPrefix(raw, "token="))
		if idx := strings.Index(token, ";"); idx >= 0 {
			token = strings.TrimSpace(token[:idx])
		}
		if idx := strings.Index(token, "&"); idx >= 0 {
			token = strings.TrimSpace(token[:idx])
		}
	}
	if token == "" {
		return errors.New("Cookie 中未找到 token")
	}

	cache.SetToken(token)
	cache.SetCookie(raw)
	if len(cookies) > 0 {
		cache.SetCookies(cookies)
	}

	listJson, err := cache.CourseListApi(5, nil)
	if err != nil {
		return fmt.Errorf("英华 Cookie 登录失败: %v", err)
	}
	if !strings.Contains(listJson, "获取数据成功") {
		return fmt.Errorf("英华 Cookie 登录失败，请确认账号与 Cookie 是否匹配：%s", truncate(listJson, 200))
	}
	return nil
}

func truncate(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}
