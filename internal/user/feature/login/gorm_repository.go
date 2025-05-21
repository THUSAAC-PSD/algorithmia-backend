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
		Where("username = ?", username).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}

		return nil, err
	}

	return &User{
		UserID:         user.UserID,
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		Email:          user.Email,
	}, nil
}
