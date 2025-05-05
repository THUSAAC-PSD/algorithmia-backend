package infrastructure

import (
	"context"
	"encoding/gob"

	ctxmiddleware "github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb/middleware/context"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/login"

	"emperror.dev/errors"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
)

const (
	sessionName = "session"
)

type sessionUser struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

func init() {
	gob.Register(sessionUser{})
}

type HTTPSessionManager struct{}

func NewHTTPSessionManager() *HTTPSessionManager {
	return &HTTPSessionManager{}
}

func (m *HTTPSessionManager) SetUser(ctx context.Context, user login.User) error {
	eCtx := ctxmiddleware.FromContext(ctx)
	if eCtx == nil {
		return errors.New("echo context not found")
	}

	sess, err := session.Get(sessionName, eCtx)
	if err != nil {
		return errors.WrapIf(err, "failed to get session")
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	sess.Values["user"] = sessionUser{
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

	sess, err := session.Get(sessionName, eCtx)
	if err != nil {
		return errors.WrapIf(err, "failed to get session")
	}

	sess.Options.MaxAge = -1
	if err := sess.Save(eCtx.Request(), eCtx.Response()); err != nil {
		return errors.WrapIf(err, "failed to save expired session")
	}

	return nil
}
