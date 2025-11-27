package manageuser

import (
	"context"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"

	"github.com/google/uuid"
)

type Repository interface {
	ListUsers(ctx context.Context) ([]ResponseUser, error)
	ListRoles(ctx context.Context) ([]ResponseRole, error)
	GetUserWithRoles(ctx context.Context, userID uuid.UUID) (*database.User, error)
	ExistsEmail(ctx context.Context, email string, excludeUserID uuid.UUID) (bool, error)
	ExistsUsername(ctx context.Context, username string, excludeUserID uuid.UUID) (bool, error)
	GetRolesByNames(ctx context.Context, names []string) ([]database.Role, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, username, email string, roles []database.Role) (*ResponseUser, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	CountSuperAdmins(ctx context.Context) (int64, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error
}
