package requestemailverification

import (
	"context"
	"html/template"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/mailing"

	"emperror.dev/errors"
	gomailpkg "github.com/wneessen/go-mail"
)

type GomailEmailSender struct {
	opts             *mailing.Options
	bodyHTMLTemplate *template.Template
}

func NewGomailEmailSender(opts *mailing.Options) (*GomailEmailSender, error) {
	tmpl, err := template.ParseFiles("internal/user/feature/requestemailverification/email_template.gohtml")
	if err != nil {
		return nil, errors.WrapIf(err, "failed to parse email template")
	}

	return &GomailEmailSender{
		opts:             opts,
		bodyHTMLTemplate: tmpl,
	}, nil
}

func (s *GomailEmailSender) SendVerificationEmail(ctx context.Context, email string, code string) error {
	subject := "【清华大学学生算法协会】邮箱验证码"

	message := gomailpkg.NewMsg()
	if err := message.From(s.opts.Sender); err != nil {
		return errors.WrapIf(err, "failed to set From address")
	}

	if err := message.To(email); err != nil {
		return errors.WrapIf(err, "failed to set To address")
	}

	message.Subject(subject)

	message.AddAlternativeString(
		gomailpkg.TypeTextPlain,
		"您好！\n\n感谢您使用清华大学学生算法协会 (THUSAAC) 的服务。您正在进行邮箱验证。\n\n您的邮箱验证码是：\n\n{{.Code}}\n\n请在 10 分钟内 将此验证码输入到验证页面，以完成您的操作。\n\n---\n安全提示：\n为保障您的账户安全，请勿将此验证码分享给任何人。如果您并未请求此验证码，请忽略本邮件，您的账户仍然是安全的。\n---\n\n此致，\n清华大学学生算法协会 (THUSAAC) 团队\n\n这是一封自动发送的邮件，请勿直接回复。",
	)
	if err := message.AddAlternativeHTMLTemplate(s.bodyHTMLTemplate, struct {
		Code string
	}{
		Code: code,
	}); err != nil {
		return errors.WrapIf(err, "failed to set HTML body")
	}

	client, err := gomailpkg.NewClient(s.opts.Host,
		gomailpkg.WithPort(s.opts.Port),
		gomailpkg.WithSMTPAuth(gomailpkg.SMTPAuthPlain), gomailpkg.WithTLSPortPolicy(gomailpkg.TLSMandatory),
		gomailpkg.WithUsername(s.opts.Username), gomailpkg.WithPassword(s.opts.Password),
	)
	if err != nil {
		return errors.WrapIf(err, "failed to create SMTP client")
	}

	if err := client.DialAndSendWithContext(ctx, message); err != nil {
		return errors.WrapIf(err, "failed to send email")
	}

	return nil
}
