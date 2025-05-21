package login

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"

	"emperror.dev/errors"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var user database.User
	if err := db.WithContext(ctx).
		Preload("Roles").
		Preload("Roles.Permissions").
		Where("username = ?", username).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}

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

	return &User{
		UserID:         user.UserID,
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		Email:          user.Email,
		IsSuperAdmin:   isSuperAdmin,
		Roles:          roles,
		Permissions:    permissions,
	}, nil
}
