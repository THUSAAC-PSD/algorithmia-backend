package echoweb

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"
	ctxmiddleware "github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb/middleware/context"

	"emperror.dev/errors"
	"github.com/labstack/echo-contrib/session"
)

type SessionAuthProvider struct{}

func NewSessionAuthProvider() *SessionAuthProvider {
	return &SessionAuthProvider{}
}

func (s *SessionAuthProvider) GetUser(ctx context.Context) (*contract.AuthUser, error) {
	eCtx := ctxmiddleware.FromContext(ctx)
	if eCtx == nil {
		return nil, errors.New("echo context not found")
	}

	sess, err := session.Get(SessionName, eCtx)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to get session")
	}

	if user, ok := sess.Values[SessionUserKey].(contract.AuthUser); !ok || user.Username == "" || user.Email == "" {
		return nil, nil
	} else {
		return &user, nil
	}
}

func (s *SessionAuthProvider) MustGetUser(ctx context.Context) (contract.AuthUser, error) {
	user, err := s.GetUser(ctx)
	if err != nil {
		return contract.AuthUser{}, errors.WrapIf(err, "failed to get user")
	}

	if user == nil {
		return contract.AuthUser{}, errors.WithStack(customerror.ErrNotAuthenticated)
	}

	return *user, nil
}
