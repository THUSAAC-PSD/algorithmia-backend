package getcurrentuser

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"

	"emperror.dev/errors"
)

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
) (*Response, error) {
	user, err := h.authProvider.MustGetUser(ctx)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to get user from auth provider")
	}

	details, err := h.authProvider.MustGetUserDetails(ctx, user.UserID)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to get user details")
	}

	return &Response{
		User: ResponseUser{
			UserID:       user.UserID,
			Username:     details.Username,
			Email:        user.Email,
			IsSuperAdmin: details.IsSuperAdmin,
			Roles:        details.Roles,
			Permissions:  details.Permissions,
		},
	}, nil
}
