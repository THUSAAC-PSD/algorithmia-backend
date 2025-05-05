package infrastructure

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb"
	ctxmiddleware "github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb/middleware/context"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/login"

	"emperror.dev/errors"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
)

type HTTPSessionManager struct{}

func NewHTTPSessionManager() *HTTPSessionManager {
	return &HTTPSessionManager{}
}

func (m *HTTPSessionManager) SetUser(ctx context.Context, user login.User) error {
	eCtx := ctxmiddleware.FromContext(ctx)
	if eCtx == nil {
		return errors.New("echo context not found")
	}

	sess, err := session.Get(echoweb.SessionName, eCtx)
	if err != nil {
		return errors.WrapIf(err, "failed to get session")
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	sess.Values[echoweb.SessionUserKey] = contract.AuthUser{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
	}

	if err := sess.Save(eCtx.Request(), eCtx.Response()); err != nil {
		return errors.WrapIf(err, "failed to save session")
	}

	return nil
}

func (m *HTTPSessionManager) Delete(ctx context.Context) error {
	eCtx := ctxmiddleware.FromContext(ctx)
	if eCtx == nil {
		return errors.New("echo context not found")
	}

	sess, err := session.Get(echoweb.SessionName, eCtx)
	if err != nil {
		return errors.WrapIf(err, "failed to get session")
	}

	sess.Options.MaxAge = -1
	if err := sess.Save(eCtx.Request(), eCtx.Response()); err != nil {
		return errors.WrapIf(err, "failed to save expired session")
	}

	return nil
}
