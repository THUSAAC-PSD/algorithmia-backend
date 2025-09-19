package listtester

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"

	"emperror.dev/errors"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db: db,
	}
}

func (r *GormRepository) GetTesters(ctx context.Context) ([]ResponseTester, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var users []database.User
	if err := db.WithContext(ctx).
		Table("users u").
		Select("u.user_id, u.username, COALESCE(u.display_name, u.username) as display_name").
		Joins("LEFT JOIN user_roles ur ON ur.user_user_id = u.user_id").
		Joins("LEFT JOIN roles r ON r.role_id = ur.role_role_id").
		Joins("LEFT JOIN role_permissions rp ON rp.role_role_id = ur.role_role_id").
		Joins("LEFT JOIN permissions p ON p.permission_id = rp.permission_permission_id").
		Where("p.name IN ?", []string{
			constant.PermissionProblemTestAssigned,
			constant.PermissionProblemTestOverride,
		}).Or("r.is_super_admin = ?", true).
		Find(&users).Error; err != nil {
		return nil, errors.WrapIf(err, "failed to get users")
	}

	testers := make([]ResponseTester, 0, len(users))
	for _, user := range users {
		testers = append(testers, ResponseTester{
			UserID:      user.UserID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
		})
	}

	return testers, nil
}
