package manageuser

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"

	"emperror.dev/errors"
)

type ListQueryHandler struct {
	repo         Repository
	authProvider contract.AuthProvider
}

func NewListQueryHandler(repo Repository, authProvider contract.AuthProvider) *ListQueryHandler {
	return &ListQueryHandler{
		repo:         repo,
		authProvider: authProvider,
	}
}

func (h *ListQueryHandler) Handle(ctx context.Context) (*ListResponse, error) {
	can, err := h.authProvider.Can(ctx, constant.PermissionUserListAll)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to check permission for listing users")
	}

	if !can {
		return nil, customerror.NewNoPermissionError(constant.PermissionUserListAll)
	}

	users, err := h.repo.ListUsers(ctx)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to fetch users")
	}

	roles, err := h.repo.ListRoles(ctx)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to fetch roles")
	}

	return &ListResponse{
		Users: users,
		Roles: roles,
	}, nil
}
