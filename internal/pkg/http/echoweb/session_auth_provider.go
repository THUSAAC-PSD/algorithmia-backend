package echoweb

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"
	ctxmiddleware "github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb/middleware/context"

	"emperror.dev/errors"
	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"gorm.io/gorm"
)

type SessionAuthProvider struct {
	db *gorm.DB
}

func NewSessionAuthProvider(db *gorm.DB) *SessionAuthProvider {
	return &SessionAuthProvider{db: db}
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
	if !ok || user.Email == "" {
		return nil, nil
	}

	return &user, nil
}

func (s *SessionAuthProvider) MustGetUserDetails(
	ctx context.Context,
	userID uuid.UUID,
) (*contract.AuthUserDetails, error) {
	db := database.GetDBFromContext(ctx, s.db)

	var user database.User
	if err := db.WithContext(ctx).
		Preload("Roles").
		Preload("Roles.Permissions").
		Where("user_id", userID).
		First(&user).Error; err != nil {
		return nil, err
	}

	roles := make([]string, 0, len(user.Roles))
	permissionCount := 0
	isSuperAdmin := false

	for _, role := range user.Roles {
		roles = append(roles, role.Name)
		if role.Permissions != nil {
			permissionCount += len(*role.Permissions)
		}

		if role.IsSuperAdmin {
			isSuperAdmin = true
		}
	}

	uniquePermissions := make(map[string]struct{})
	for _, role := range user.Roles {
		if role.Permissions != nil {
			for _, permission := range *role.Permissions {
				uniquePermissions[permission.Name] = struct{}{}
			}
		}
	}

	permissions := make([]string, 0, len(uniquePermissions))
	for permission := range uniquePermissions {
		permissions = append(permissions, permission)
	}

	return &contract.AuthUserDetails{
		Username:     user.Username,
		IsSuperAdmin: isSuperAdmin,
		Roles:        roles,
		Permissions:  permissions,
	}, nil
}

func (s *SessionAuthProvider) Can(ctx context.Context, permissionName string) (bool, error) {
	user, err := s.MustGetUser(ctx)
	if err != nil {
		return false, errors.WrapIf(err, "failed to get user")
	}

	db := database.GetDBFromContext(ctx, s.db)

	// TODO: Cache
	var permissionNames []string
	if err := db.WithContext(ctx).
		Model(&database.User{
			UserID: user.UserID,
		}).
		Select("permissions.name").
		Preload("Permissions").
		Scan(&permissionNames).Error; err != nil {
		return false, errors.WrapIf(err, "failed to get permissions")
	}

	for _, p := range permissionNames {
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
