package resetpassword

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"

	"emperror.dev/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var user database.User
	if err := db.WithContext(ctx).
		Select("user_id", "hashed_password").
		Where("user_id = ?", userID).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, errors.WrapIf(err, "failed to load user by id")
	}

	return &User{
		UserID:         user.UserID,
		HashedPassword: user.HashedPassword,
	}, nil
}

func (r *GormRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error {
	db := database.GetDBFromContext(ctx, r.db)

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user database.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("user_id", "hashed_password").
			Where("user_id = ?", userID).
			First(&user).Error; err != nil {
			return err
		}

		if err := tx.Model(&user).
			Update("hashed_password", hashedPassword).Error; err != nil {
			return err
		}

		return nil
	})
}
