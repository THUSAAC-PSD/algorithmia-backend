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

	user, ok := sess.Values[SessionUserKey].(contract.AuthUser)
	if !ok || user.Username == "" || user.Email == "" {
		return nil, nil
	}

	return &user, nil
}

func (s *SessionAuthProvider) Can(ctx context.Context, permissionName string) (bool, error) {
	user, err := s.MustGetUser(ctx)
	if err != nil {
		return false, errors.WrapIf(err, "failed to get user")
	}

	for _, p := range user.Permissions {
		if p == permissionName {
			return true, nil
		}
	}

	return false, nil
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
