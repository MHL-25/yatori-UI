package yinghuahelper

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/yatori-dev/yatori-go-core/aggregation/yinghua"
	yinghuaApi "github.com/yatori-dev/yatori-go-core/api/yinghua"
	"yatori-UI/internal/config"
)

// Target 英华单个平台登录目标
type Target struct {
	URL      string
	Password string
}

// Targets 解析账号配置中的英华平台列表（含旧版单 URL 回退）
func Targets(user config.User) []Target {
	items := user.YingHuaLoginTargets()
	targets := make([]Target, 0, len(items))
	for _, item := range items {
		targets = append(targets, Target{URL: item.URL, Password: item.Password})
	}
	return targets
}

// LoginAccount 登录单个账号在各平台的所有入口，返回成功会话与失败说明
func LoginAccount(account string, user config.User) (caches []*yinghuaApi.YingHuaUserCache, failures []string) {
	targets := Targets(user)
	if len(targets) == 0 {
		return nil, []string{"未配置英华平台URL"}
	}
	proxy := user.IsProxy == 1
	for _, target := range targets {
		cache, err := loginWithAccount(target, account, proxy)
		if err != nil {
			failures = append(failures, fmt.Sprintf("%s 登录失败: %v", Host(target.URL), err))
			continue
		}
		caches = append(caches, cache)
	}
	return caches, failures
}

// LoginTarget 登录单个平台入口，返回已缓存的 Cookie/Token 会话
func LoginTarget(target Target, account string, useProxy bool) (*yinghuaApi.YingHuaUserCache, error) {
	return loginWithAccount(target, account, useProxy)
}

func loginWithAccount(target Target, account string, useProxy bool) (*yinghuaApi.YingHuaUserCache, error) {
	cache := &yinghuaApi.YingHuaUserCache{
		PreUrl:   target.URL,
		Account:  account,
		Password: target.Password,
	}
	if useProxy {
		cache.IpProxySW = true
	}
	var err error
	if yinghua.IsYingHuaCookieOrToken(target.Password) {
		err = yinghua.YingHuaCookieLoginAction(cache)
	} else {
		err = yinghua.YingHuaLoginAction(cache)
	}
	if err != nil {
		return nil, err
	}
	return cache, nil
}

// Host 返回 URL 的站点主机名，用于日志区分平台
func Host(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Host == "" {
		return raw
	}
	return parsed.Host
}
