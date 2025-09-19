package verifyemail

import (
	"context"
	"fmt"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/feature/login"
	"github.com/google/uuid"
	"emperror.dev/errors"
	"gorm.io/gorm"
)

type CommandHandler struct {
	repository     Repository
	sessionManager login.SessionManager
}

func NewCommandHandler(repository Repository, sessionManager login.SessionManager) *CommandHandler {
	return &CommandHandler{
		repository:     repository,
		sessionManager: sessionManager,
	}
}

func (h *CommandHandler) Handle(ctx context.Context, command *Command) (*Result, error) {
	// Get and verify the email verification code
	verificationCode, err := h.repository.GetEmailVerificationCode(ctx, command.Email, command.Token)
	if err != nil {
		return nil, err
	}

	// Check if user already exists
	var existingUser database.User
	if err := h.repository.(*GormRepository).db.WithContext(ctx).Where("email = ?", command.Email).First(&existingUser).Error; err == nil {
		return nil, ErrUserAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.WrapIf(err, "failed to check existing user")
	}

	// Resolve potential username conflict by generating a unique username
	username := verificationCode.Username
	if username == "" {
		username = command.Email // fallback
	}

	if taken, err := h.repository.UsernameExists(ctx, username); err != nil {
		return nil, errors.WrapIf(err, "failed to check username existence")
	} else if taken {
		base := username
		for i := 1; i <= 1000; i++ {
			candidate := fmt.Sprintf("%s%d", base, i)
			ok, err := h.repository.UsernameExists(ctx, candidate)
			if err != nil {
				return nil, errors.WrapIf(err, "failed to check username existence")
			}
			if !ok {
				username = candidate
				break
			}
		}
	}

	// Create the user with the stored registration data
	newUserID, err := uuid.NewV7()
	if err != nil {
		return nil, errors.WrapIf(err, "failed to generate user id")
	}

	user := &database.User{
		UserID:         newUserID,
		Username:       username,
		HashedPassword: verificationCode.PasswordHash,
		Email:          verificationCode.Email,
		CreatedAt:      time.Now(),
	}

	if err := h.repository.CreateUser(ctx, user); err != nil {
		return nil, errors.WrapIf(err, "failed to create user")
	}

	// Log the user in by setting the session
	if err := h.sessionManager.SetUser(ctx, login.User{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
	}); err != nil {
		return nil, errors.WrapIf(err, "failed to set user session")
	}

	// Clean up the verification code
	if err := h.repository.DeleteEmailVerificationCode(ctx, verificationCode.EmailVerificationCodeID); err != nil {
		// Log error but don't fail the operation
		// Logger would be injected in a real implementation
	}

	return &Result{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
