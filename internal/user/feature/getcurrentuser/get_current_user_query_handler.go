package getcurrentuser

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"

	"emperror.dev/errors"
)

type Query struct{}

type QueryHandler struct {
	authProvider contract.AuthProvider
}

func NewQueryHandler(authProvider contract.AuthProvider) *QueryHandler {
	return &QueryHandler{
		authProvider: authProvider,
	}
}

func (h *QueryHandler) Handle(
	ctx context.Context,
	query *Query,
) (*Response, error) {
	if query == nil {
		return nil, errors.WithStack(customerror.ErrCommandNil)
	}

	user, err := h.authProvider.MustGetUser(ctx)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to get user from auth provider")
	}

	return &Response{
		User: ResponseUser{
			UserID:   user.UserID,
			Username: user.Username,
			Email:    user.Email,
			Roles:    user.Roles,
		},
	}, nil
}
