package manageuser

import (
	"context"
	"strings"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type UpdateCommandHandler struct {
	repo         Repository
	validator    *validator.Validate
	authProvider contract.AuthProvider
}

func NewUpdateCommandHandler(
	repo Repository,
	validator *validator.Validate,
	authProvider contract.AuthProvider,
) *UpdateCommandHandler {
	return &UpdateCommandHandler{
		repo:         repo,
		validator:    validator,
		authProvider: authProvider,
	}
}

func (h *UpdateCommandHandler) Handle(
	ctx context.Context,
	command *UpdateCommand,
) (*UpdateResponse, error) {
	if command == nil {
		return nil, errors.WithStack(customerror.ErrCommandNil)
	}

	if err := h.validator.Struct(command); err != nil {
		return nil, errors.WithStack(errors.Append(err, customerror.ErrValidationFailed))
	}

	can, err := h.authProvider.Can(ctx, constant.PermissionUserManageRolesAny)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to check permission for updating user")
	}

	if !can {
		return nil, customerror.NewNoPermissionError(constant.PermissionUserManageRolesAny)
	}

	currentUser, err := h.authProvider.MustGetUser(ctx)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to get current user")
	}

	targetUser, err := h.repo.GetUserWithRoles(ctx, command.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithStack(ErrUserNotFound)
		}

		return nil, errors.WrapIf(err, "failed to load target user")
	}

	normalizedRoles := normalizeRoles(command.Roles)
	if len(normalizedRoles) == 0 {
		return nil, errors.WithStack(ErrRolesRequired)
	}

	targetIsSuperAdmin := userHasRole(targetUser.Roles, "super_admin")
	requestRemovesSuperAdmin := targetIsSuperAdmin && !containsRole(normalizedRoles, "super_admin")

	if requestRemovesSuperAdmin && targetUser.UserID == currentUser.UserID {
		return nil, errors.WithStack(ErrCannotRemoveOwnSuperAdmin)
	}

	if requestRemovesSuperAdmin {
		count, err := h.repo.CountSuperAdmins(ctx)
		if err != nil {
			return nil, errors.WrapIf(err, "failed to count super admins")
		}

		if count <= 1 {
			return nil, errors.WithStack(ErrCannotDeleteLastSuperAdmin)
		}
	}

	if exists, err := h.repo.ExistsEmail(ctx, command.Email, command.UserID); err != nil {
		return nil, errors.WrapIf(err, "failed to check email duplicates")
	} else if exists {
		return nil, errors.WithStack(ErrEmailAlreadyExists)
	}

	if exists, err := h.repo.ExistsUsername(ctx, command.Username, command.UserID); err != nil {
		return nil, errors.WrapIf(err, "failed to check username duplicates")
	} else if exists {
		return nil, errors.WithStack(ErrUsernameAlreadyExists)
	}

	roles, err := h.repo.GetRolesByNames(ctx, normalizedRoles)
	if err != nil {
		return nil, err
	}

	updatedUser, err := h.repo.UpdateUser(
		ctx,
		command.UserID,
		command.Username,
		command.Email,
		roles,
	)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to update user")
	}

	return &UpdateResponse{User: *updatedUser}, nil
}

func normalizeRoles(roles []string) []string {
	unique := make(map[string]struct{}, len(roles))
	result := make([]string, 0, len(roles))

	for _, role := range roles {
		roleName := strings.TrimSpace(role)
		if roleName == "" {
			continue
		}
		roleName = strings.ToLower(roleName)
		if _, exists := unique[roleName]; exists {
			continue
		}
		unique[roleName] = struct{}{}
		result = append(result, roleName)
	}

	return result
}

func userHasRole(roles []database.Role, target string) bool {
	for _, role := range roles {
		if role.Name == target {
			return true
		}
	}

	return false
}

func containsRole(roles []string, target string) bool {
	target = strings.ToLower(target)
	for _, role := range roles {
		if role == target {
			return true
		}
	}

	return false
}
