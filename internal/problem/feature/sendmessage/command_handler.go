package sendmessage

import (
	"context"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract/uowhelper"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

var (
	ErrUserNotPartOfRoom = errors.New("user is not part of room")
	ErrMediaNotFound     = errors.New("media not found")
)

type Repository interface {
	IsUserPartOfRoom(ctx context.Context, problemID uuid.UUID, userID uuid.UUID) (bool, error)
	CreateChatMessage(ctx context.Context, command *Command, senderID uuid.UUID, createdAt time.Time) (uuid.UUID, error)
	GetAttachmentByMediaIDs(ctx context.Context, mediaIDs []uuid.UUID) ([]contract.MessageAttachment, error)
}

type CommandHandler struct {
	repo         Repository
	validator    *validator.Validate
	authProvider contract.AuthProvider
	uowFactory   contract.UnitOfWorkFactory
	broadcaster  contract.MessageBroadcaster
	l            logger.Logger
}

func NewCommandHandler(
	repo Repository,
	validator *validator.Validate,
	authProvider contract.AuthProvider,
	uowFactory contract.UnitOfWorkFactory,
	broadcaster contract.MessageBroadcaster,
	l logger.Logger,
) *CommandHandler {
	return &CommandHandler{
		repo:         repo,
		validator:    validator,
		authProvider: authProvider,
		uowFactory:   uowFactory,
		broadcaster:  broadcaster,
		l:            l,
	}
}

func (h *CommandHandler) Handle(ctx context.Context, command *Command) (mediatr.Unit, error) {
	if command == nil {
		return mediatr.Unit{}, errors.WithStack(customerror.ErrCommandNil)
	}

	if err := h.validator.Struct(command); err != nil {
		return mediatr.Unit{}, errors.WithStack(errors.Append(err, customerror.ErrValidationFailed))
	}

	user, err := h.authProvider.MustGetUser(ctx)
	if err != nil {
		return mediatr.Unit{}, errors.WrapIf(err, "failed to get user from auth provider")
	}

	uow := h.uowFactory.New()
	return uowhelper.DoWithResult(ctx, uow, h.l, func(ctx context.Context) (mediatr.Unit, error) {
		if ok, err := h.repo.IsUserPartOfRoom(ctx, command.ProblemID, user.UserID); err != nil {
			return mediatr.Unit{}, errors.WrapIf(err, "failed to check if user is part of room")
		} else if !ok {
			return mediatr.Unit{}, errors.WithStack(ErrUserNotPartOfRoom)
		}

		timestamp := time.Now()

		messageID, err := h.repo.CreateChatMessage(ctx, command, user.UserID, timestamp)
		if err != nil {
			return mediatr.Unit{}, errors.WrapIf(err, "failed to create chat message")
		}

		attachments, err := h.repo.GetAttachmentByMediaIDs(ctx, command.AttachmentMediaIDs)
		if err != nil {
			return mediatr.Unit{}, errors.WrapIf(err, "failed to get attachments")
		}

		if err := h.broadcaster.BroadcastUserMessage(command.ProblemID, messageID, command.Content, contract.MessageUser{
			UserID:   user.UserID,
			Username: user.Username,
		}, attachments, timestamp); err != nil {
			return mediatr.Unit{}, errors.WrapIf(err, "failed to broadcast message")
		}

		return mediatr.Unit{}, nil
	})
}
