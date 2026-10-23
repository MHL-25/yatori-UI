package notify

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SendShowDoc 向多个 ShowDoc 推送地址发送通知。
// 地址格式：https://push.showdoc.com.cn/server/api/push/{token}
func SendShowDoc(addresses []string, title, content string) error {
	if len(addresses) == 0 {
		return nil
	}
	client := &http.Client{Timeout: 20 * time.Second}
	var lastErr error
	sent := 0
	seen := map[string]bool{}
	for _, addr := range addresses {
		addr = strings.TrimSpace(addr)
		if addr == "" || seen[addr] {
			continue
		}
		seen[addr] = true
		form := url.Values{}
		form.Set("title", title)
		form.Set("content", content)
		resp, err := client.PostForm(addr, form)
		if err != nil {
			lastErr = err
			continue
		}
		resp.Body.Close()
		if resp.StatusCode >= 300 {
			lastErr = fmt.Errorf("ShowDoc 返回状态码 %d", resp.StatusCode)
			continue
		}
		sent++
	}
	if sent == 0 && lastErr != nil {
		return lastErr
	}
	return nil
}
