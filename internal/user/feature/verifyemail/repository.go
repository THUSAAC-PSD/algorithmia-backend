package verifyemail

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"emperror.dev/errors"
)

type Repository interface {
	GetEmailVerificationCode(ctx context.Context, email string, code string) (*database.EmailVerificationCode, error)
	DeleteEmailVerificationCode(ctx context.Context, id uuid.UUID) error
	CreateUser(ctx context.Context, user *database.User) error
	UsernameExists(ctx context.Context, username string) (bool, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) GetEmailVerificationCode(ctx context.Context, email string, code string) (*database.EmailVerificationCode, error) {
	var verificationCode database.EmailVerificationCode
	// Compare expiration against DB server time to avoid timezone discrepancies
	err := r.db.WithContext(ctx).Where("email = ? AND code = ? AND expires_at > NOW()",
		email, code).First(&verificationCode).Error
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInvalidOrExpiredToken
	}
	
	if err != nil {
		return nil, errors.WrapIf(err, "failed to get email verification code")
	}
	
	return &verificationCode, nil
}

func (r *GormRepository) DeleteEmailVerificationCode(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&database.EmailVerificationCode{}, "email_verification_code_id = ?", id).Error
}

func (r *GormRepository) CreateUser(ctx context.Context, user *database.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *GormRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&database.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, errors.WrapIf(err, "failed to check username existence")
	}
	return count > 0, nil
}
