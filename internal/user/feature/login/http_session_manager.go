package login

import (
	"context"
	"encoding/gob"

	"emperror.dev/errors"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
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

func (m *HTTPSessionManager) SetUser(ctx context.Context, user User) error {
	eCtx, ok := ctx.Value("echo.context").(echo.Context)
	if !ok {
		return errors.New("echo context not found")
	}

	sess, err := session.Get("session", eCtx)
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
