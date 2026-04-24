package email

import (
	"crypto/tls"
	"fmt"
	"html"
	"strings"

	"gopkg.in/gomail.v2"
)

func SendMail(host string, port int, userName, password string, toMail []string, content string) {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(userName, "Yatori课程助手"))
	m.SetHeader("To", toMail...)
	m.SetHeader("Subject", "Yatori课程助手通知")
	emailHTML := buildEmailHTML("Yatori课程助手", content, false)
	m.SetBody("text/html", emailHTML)

	d := gomail.NewDialer(host, port, userName, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("邮件发送失败: host=%s port=%d user=%s err=%v\n", host, port, userName, err)
	}
}

func buildEmailHTML(title string, contentHTML string, asPlainText bool) string {
	logoURL := "https://avatars.githubusercontent.com/u/185567923?s=1000&v=4"

	if asPlainText {
		contentHTML = html.EscapeString(contentHTML)
		contentHTML = strings.ReplaceAll(contentHTML, "\n", "<br>")
	}

	return fmt.Sprintf(`<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="x-apple-disable-message-reformatting">
  <title>%s</title>
</head>
<body style="margin:0;padding:0;background:#f5f7fb;">
  <table role="presentation" cellpadding="0" cellspacing="0" width="100%%" style="background:#f5f7fb;">
    <tr>
      <td align="center" style="padding:32px 16px;">
        <table role="presentation" cellpadding="0" cellspacing="0" width="600" style="max-width:600px;background:#ffffff;border-radius:16px;box-shadow:0 6px 24px rgba(18,38,63,0.08);">
          <tr>
            <td align="center" style="padding:28px 24px 8px 24px;">
              <img src="%s" width="88" height="88" alt="logo" style="display:block;border-radius:50%%;width:88px;height:88px;border:2px solid #eef2f7;object-fit:cover;">
            </td>
          </tr>
          <tr>
            <td align="center" style="padding:0 24px 8px 24px;">
              <div style="font-family:system-ui,-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica,Arial,sans-serif;font-size:22px;font-weight:700;color:#111827;line-height:1.3;">%s</div>
            </td>
          </tr>
          <tr>
            <td style="padding:8px 24px 0 24px;">
              <div style="height:1px;background:linear-gradient(90deg,#e5e7eb,#f3f4f6,#e5e7eb);"></div>
            </td>
          </tr>
          <tr>
            <td style="padding:18px 24px 8px 24px;">
              <div style="font-family:system-ui,-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica,Arial,sans-serif;font-size:15px;color:#374151;line-height:1.8;">
                %s
              </div>
            </td>
          </tr>
        </table>
        <table role="presentation" cellpadding="0" cellspacing="0" width="600" style="max-width:600px;">
          <tr><td align="center" style="padding:14px 8px 0 8px;color:#6b7280;font-family:system-ui,sans-serif;font-size:12px;line-height:1.6;">
            这是一封系统通知邮件，请勿直接回复。
          </td></tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`,
		html.EscapeString(title),
		logoURL,
		html.EscapeString(title),
		contentHTML,
	)
}
