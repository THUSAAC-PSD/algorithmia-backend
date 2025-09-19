package requestemailverification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/postmark"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/constant"

	"emperror.dev/errors"
)

type PostmarkEmailSender struct {
	opts             *postmark.Options
	bodyHTMLTemplate *template.Template
}

type PostmarkEmailRequest struct {
	From     string `json:"From"`
	To       string `json:"To"`
	Subject  string `json:"Subject"`
	HtmlBody string `json:"HtmlBody"`
	TextBody string `json:"TextBody"`
	Tag      string `json:"Tag,omitempty"`
}

func NewPostmarkEmailSender(opts *postmark.Options) (*PostmarkEmailSender, error) {
	tmplFileName := filepath.Join("resources", "request_email_verification_template.gohtml")

	tmpl, err := template.ParseFiles(tmplFileName)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to parse email template")
	}

	return &PostmarkEmailSender{
		opts:             opts,
		bodyHTMLTemplate: tmpl,
	}, nil
}

func (s *PostmarkEmailSender) SendVerificationEmail(ctx context.Context, email string, code string) error {
	subject := "【清华大学学生算法协会】注册确认邮件"
	
	// URL-encode token and email to ensure special characters (like '+') are preserved correctly
	encodedToken := url.QueryEscape(code)
	encodedEmail := url.QueryEscape(email)

	// Create verification link with the encoded params
	verificationLink := fmt.Sprintf("https://algorithmia.thusaac.com/verify-email?token=%s&email=%s", encodedToken, encodedEmail)
	
	// Generate HTML content
	var htmlBuf bytes.Buffer
	if err := s.bodyHTMLTemplate.Execute(&htmlBuf, struct {
		VerificationLink  string
		ValidDurationMins int
	}{
		VerificationLink:  verificationLink,
		ValidDurationMins: constant.EmailVerificationValidDurationMins,
	}); err != nil {
		return errors.WrapIf(err, "failed to execute HTML template")
	}

	// Plain text fallback
	textBody := fmt.Sprintf(`您好！

感谢您使用清华大学学生算法协会 (THUSAAC) 的服务。您正在进行账户注册。

请点击以下链接完成您的账户注册：

%s

此链接将在 %d 分钟内有效。如果链接无法点击，请复制粘贴到浏览器地址栏中。

---
安全提示：
为保障您的账户安全，请勿将此链接分享给任何人。如果您并未请求注册账户，请忽略本邮件。
---

此致，
清华大学学生算法协会 (THUSAAC) 团队

这是一封自动发送的邮件，请勿直接回复。`, verificationLink, constant.EmailVerificationValidDurationMins)

	// Create Postmark email request
	emailRequest := PostmarkEmailRequest{
		From:     s.opts.FromEmail,
		To:       email,
		Subject:  subject,
		HtmlBody: htmlBuf.String(),
		TextBody: textBody,
		Tag:      "email-verification",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(emailRequest)
	if err != nil {
		return errors.WrapIf(err, "failed to marshal email request")
	}

	// Create HTTP request to Postmark API
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.postmarkapp.com/email", bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.WrapIf(err, "failed to create HTTP request")
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", s.opts.ServerToken)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.WrapIf(err, "failed to send HTTP request to Postmark")
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var responseBody bytes.Buffer
		responseBody.ReadFrom(resp.Body)
		return errors.Errorf("Postmark API returned status %d: %s", resp.StatusCode, responseBody.String())
	}

	return nil
}
