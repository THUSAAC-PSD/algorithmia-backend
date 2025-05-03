package register

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
	return &GormRepository{
		db: db,
	}
}

func (r *GormRepository) CreateUser(ctx context.Context, user User) error {
	userModel := database.User{
		UserID:         user.UserID,
		Username:       user.Username,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		CreatedAt:      user.CreatedAt,
	}

	if err := r.db.WithContext(ctx).Create(&userModel).Error; err != nil {
		return errors.WrapIf(err, "failed to create user")
	}
	return nil
}

func (r *GormRepository) IsUserUnique(ctx context.Context, username string, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&database.User{}).Where("username = ? OR email = ?", username, email).Count(&count).Error; err != nil {
		return false, errors.WrapIf(err, "failed to check if user is unique")
	}
	return count == 0, nil
}

func (r *GormRepository) CheckAndDeleteEmailVerificationCode(
	ctx context.Context,
	email string,
	code string,
) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&database.EmailVerificationCode{}).Where("email = ? AND code = ?", email, code).Count(&count).Error; err != nil {
		return false, errors.WrapIf(err, "failed to check email verification code")
	}

	if count > 0 {
		if err := r.db.WithContext(ctx).Where("email = ? AND code = ?", email, code).Delete(&database.EmailVerificationCode{}).Error; err != nil {
			return false, errors.WrapIf(err, "failed to delete email verification code")
		}
		return true, nil
	}
	return false, nil
}
