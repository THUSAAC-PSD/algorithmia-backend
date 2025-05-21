package sendmessage

import (
	"context"
	"time"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"

	"emperror.dev/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) IsUserPartOfRoom(ctx context.Context, problemID uuid.UUID, userID uuid.UUID) (bool, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var p database.Problem
	if err := db.WithContext(ctx).
		Model(&database.Problem{}).
		Preload("Testers").
		Where("problem_id = ?", problemID).
		First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, errors.WrapIf(err, "failed to check if user is part of room")
	}

	if p.CreatorID == userID || p.ReviewerID.UUID == userID {
		return true, nil
	}

	if len(p.Testers) > 0 {
		for _, tester := range p.Testers {
			if tester.UserID == userID {
				return true, nil
			}
		}

		return false, nil
	}

	return true, nil
}

func (r *GormRepository) CreateChatMessage(
	ctx context.Context,
	command *Command,
	senderID uuid.UUID,
	createdAt time.Time,
) (uuid.UUID, error) {
	db := database.GetDBFromContext(ctx, r.db)

	messageID, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, errors.WrapIf(err, "failed to generate message ID")
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		var problem database.Problem
		if err := tx.WithContext(ctx).
			Select("problem_id").
			Where("problem_id = ?", command.ProblemID).
			First(&problem).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.WrapIf(err, "problem not found")
			}

			return errors.WrapIf(err, "failed to find problem")
		}

		var mediaIDs []uuid.UUID
		if err := tx.WithContext(ctx).
			Model(&database.Media{}).
			Where("media_id IN ?", command.AttachmentMediaIDs).
			Pluck("media_id", &mediaIDs).Error; err != nil {
			return errors.WrapIf(err, "failed to check media existence")
		}

		if len(mediaIDs) != len(command.AttachmentMediaIDs) {
			return errors.WithStack(ErrMediaNotFound)
		}

		chatMessage := &database.ProblemChatMessage{
			MessageID:   messageID,
			ProblemID:   command.ProblemID,
			SenderID:    senderID,
			Content:     command.Content,
			CreatedAt:   createdAt,
			Attachments: make([]database.ProblemChatMessageAttachment, 0),
		}

		for _, mediaID := range command.AttachmentMediaIDs {
			attachmentID, err := uuid.NewV7()
			if err != nil {
				return errors.WrapIf(err, "failed to generate attachment ID")
			}

			chatMessage.Attachments = append(chatMessage.Attachments, database.ProblemChatMessageAttachment{
				AttachmentID: attachmentID,
				MessageID:    messageID,
				MediaID:      mediaID,
			})
		}

		if err := tx.WithContext(ctx).Create(chatMessage).Error; err != nil {
			return errors.WrapIf(err, "failed to create chat message")
		}

		return nil
	}); err != nil {
		return uuid.Nil, errors.WrapIf(err, "failed to run transaction")
	}

	return messageID, nil
}

func (r *GormRepository) GetAttachmentByMediaIDs(
	ctx context.Context,
	mediaIDs []uuid.UUID,
) ([]contract.MessageAttachment, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var attachments []contract.MessageAttachment
	if err := db.WithContext(ctx).
		Model(&database.Media{}).
		Where("media_id IN ?", mediaIDs).
		Select("url, file_name, mime_type, file_size").
		Scan(&attachments).Error; err != nil {
		return nil, errors.WrapIf(err, "failed to get attachments by media IDs")
	}

	return attachments, nil
}
