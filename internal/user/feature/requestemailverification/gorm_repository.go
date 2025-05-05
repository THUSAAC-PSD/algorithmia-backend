package requestemailverification

import (
	"context"
	"fmt"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"

	"emperror.dev/errors"
	"github.com/google/uuid"
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

func (r *GormRepository) CreateEmailVerificationCode(ctx context.Context, email string, code string) error {
	db := database.GetDBFromContext(ctx, r.db)

	id, err := uuid.NewV7()
	if err != nil {
		return errors.WrapIf(err, "failed to generate UUID")
	}

	model := &database.EmailVerificationCode{
		EmailVerificationCodeID: id,
		Email:                   email,
		Code:                    code,
	}

	if err := db.WithContext(ctx).Create(model).Error; err != nil {
		return errors.WrapIf(err, "failed to create email verification code")
	}

	return nil
}

func (r *GormRepository) IsNotTimedOut(ctx context.Context, email string) (bool, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var count int64
	if err := db.WithContext(ctx).Model(&database.EmailVerificationCode{}).
		Where(fmt.Sprintf("email = ? AND created_at >= NOW() - INTERVAL '%d' MINUTE", timeoutDurationMins), email).
		Count(&count).Error; err != nil {
		return false, errors.WrapIf(err, "failed to check if email is not timed out")
	}

	return count == 0, nil
}

func (r *GormRepository) IsNotAssociatedWithUser(ctx context.Context, email string) (bool, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var count int64
	if err := db.WithContext(ctx).Model(&database.User{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		return false, errors.WrapIf(err, "failed to check if email is not associated with user")
	}

	return count == 0, nil
}
