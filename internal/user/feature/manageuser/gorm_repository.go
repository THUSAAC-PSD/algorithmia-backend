package manageuser

import (
	"context"
	"sort"
	"strings"

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

func (r *GormRepository) ListUsers(ctx context.Context) ([]ResponseUser, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var users []database.User
	if err := db.WithContext(ctx).
		Preload("Roles").
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, errors.WrapIf(err, "failed to list users")
	}

	responses := make([]ResponseUser, len(users))
	for i, user := range users {
		responses[i] = toResponseUser(user)
	}

	return responses, nil
}

func (r *GormRepository) ListRoles(ctx context.Context) ([]ResponseRole, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var roles []database.Role
	if err := db.WithContext(ctx).
		Order("LOWER(name) ASC").
		Find(&roles).Error; err != nil {
		return nil, errors.WrapIf(err, "failed to list roles")
	}

	responses := make([]ResponseRole, len(roles))
	for i, role := range roles {
		responses[i] = toResponseRole(role)
	}

	return responses, nil
}

func (r *GormRepository) GetUserWithRoles(ctx context.Context, userID uuid.UUID) (*database.User, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var user database.User
	if err := db.WithContext(ctx).
		Preload("Roles").
		Where("user_id = ?", userID).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *GormRepository) ExistsEmail(
	ctx context.Context,
	email string,
	excludeUserID uuid.UUID,
) (bool, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var count int64
	if err := db.WithContext(ctx).
		Model(&database.User{}).
		Where("LOWER(email) = LOWER(?) AND user_id <> ?", email, excludeUserID).
		Count(&count).Error; err != nil {
		return false, errors.WrapIf(err, "failed to check duplicate email")
	}

	return count > 0, nil
}

func (r *GormRepository) ExistsUsername(
	ctx context.Context,
	username string,
	excludeUserID uuid.UUID,
) (bool, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var count int64
	if err := db.WithContext(ctx).
		Model(&database.User{}).
		Where("LOWER(username) = LOWER(?) AND user_id <> ?", username, excludeUserID).
		Count(&count).Error; err != nil {
		return false, errors.WrapIf(err, "failed to check duplicate username")
	}

	return count > 0, nil
}

func (r *GormRepository) GetRolesByNames(
	ctx context.Context,
	names []string,
) ([]database.Role, error) {
	db := database.GetDBFromContext(ctx, r.db)

	if len(names) == 0 {
		return nil, errors.WithStack(ErrRolesRequired)
	}

	var roles []database.Role
	if err := db.WithContext(ctx).
		Where("name IN ?", names).
		Find(&roles).Error; err != nil {
		return nil, errors.WrapIf(err, "failed to fetch roles")
	}

	found := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		found[role.Name] = struct{}{}
	}

	missing := make([]string, 0)
	for _, name := range names {
		if _, ok := found[name]; !ok {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		return nil, errors.Wrapf(ErrRoleNotFound, "missing roles: %s", strings.Join(missing, ", "))
	}

	return roles, nil
}

func (r *GormRepository) UpdateUser(
	ctx context.Context,
	userID uuid.UUID,
	username string,
	email string,
	roles []database.Role,
) (*ResponseUser, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var response ResponseUser
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user database.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Roles").
			Where("user_id = ?", userID).
			First(&user).Error; err != nil {
			return err
		}

		user.Username = username
		user.Email = email

		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		if err := tx.Model(&user).Association("Roles").Replace(roles); err != nil {
			return err
		}

		if err := tx.Preload("Roles").
			Where("user_id = ?", userID).
			First(&user).Error; err != nil {
			return err
		}

		response = toResponseUser(user)
		return nil
	})
	if err != nil {
		return nil, errors.WrapIf(err, "failed to update user")
	}

	return &response, nil
}

func (r *GormRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	db := database.GetDBFromContext(ctx, r.db)

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user database.User
		if err := tx.
			Where("user_id = ?", userID).
			First(&user).Error; err != nil {
			return err
		}

		if err := tx.Model(&user).Association("Roles").Clear(); err != nil {
			return err
		}

		if err := tx.Delete(&user).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return errors.WrapIf(err, "failed to delete user")
	}

	return nil
}

func (r *GormRepository) CountSuperAdmins(ctx context.Context) (int64, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var count int64
	if err := db.WithContext(ctx).
		Model(&database.User{}).
		Joins("JOIN user_roles ur ON ur.user_user_id = users.user_id").
		Joins("JOIN roles r ON r.role_id = ur.role_role_id").
		Where("r.name = ?", "super_admin").
		Count(&count).Error; err != nil {
		return 0, errors.WrapIf(err, "failed to count super admins")
	}

	return count, nil
}

func (r *GormRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error {
	db := database.GetDBFromContext(ctx, r.db)

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user database.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("user_id").
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

func toResponseUser(user database.User) ResponseUser {
	roles := make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}

	sort.Strings(roles)

	return ResponseUser{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		Roles:     roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func toResponseRole(role database.Role) ResponseRole {
	return ResponseRole{
		Name:         role.Name,
		Description:  role.Description,
		IsSuperAdmin: role.IsSuperAdmin,
	}
}
