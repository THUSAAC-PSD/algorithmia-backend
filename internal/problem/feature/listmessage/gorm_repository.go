package listmessage

import (
	"context"
	"database/sql"

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
	return &GormRepository{
		db: db,
	}
}

func (r *GormRepository) IsUserPartOfRoom(ctx context.Context, problemID uuid.UUID, userID uuid.UUID) (bool, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var count int64
	if err := db.WithContext(ctx).
		Model(&database.Problem{}).
		Joins("INNER JOIN problem_testers ON problem_testers.problem_problem_id = problems.problem_id AND user_user_id = ?", userID).
		Where("problem_id = ? AND (creator_id = ? OR reviewer_id = ?)", problemID, userID, userID).
		Count(&count).Error; err != nil {
		return false, errors.WrapIf(err, "failed to check if user is part of room")
	}

	return count > 0, nil
}

func (r *GormRepository) GetUserChatMessages(ctx context.Context, problemID uuid.UUID) ([]ResponseChatMessage, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var messages []database.ProblemChatMessage
	if err := db.WithContext(ctx).
		Model(&database.ProblemChatMessage{}).
		Preload("Attachments").
		Where("problem_id = ?", problemID).
		Find(&messages).Error; err != nil {
		return nil, errors.WrapIf(err, "failed to get user chat messages")
	}

	senderIDs := make([]uuid.UUID, 0, len(messages))
	mediaCount := 0

	for _, message := range messages {
		senderIDs = append(senderIDs, message.SenderID)
		mediaCount += len(message.Attachments)
	}

	senderMap, err := r.fetchUsers(ctx, db, senderIDs)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to fetch users")
	}

	mediaIDs := make([]uuid.UUID, 0, mediaCount)
	for _, message := range messages {
		for _, attachment := range message.Attachments {
			mediaIDs = append(mediaIDs, attachment.MediaID)
		}
	}

	var media []database.Media
	if err := db.WithContext(ctx).
		Model(&database.Media{}).
		Select("media_id, url").
		Where("media_id IN ?", mediaIDs).
		Find(&media).Error; err != nil {
		return nil, errors.WrapIf(err, "failed to get medias")
	}

	mediaMap := make(map[uuid.UUID]contract.MessageAttachment, len(media))
	for _, m := range media {
		mediaMap[m.MediaID] = contract.MessageAttachment{
			URL:      m.URL,
			FileName: m.FileName,
			MIMEType: m.MIMEType,
			Size:     m.FileSize,
		}
	}

	dtos := make([]ResponseChatMessage, 0, len(messages))
	for _, message := range messages {
		payload := ResponseChatUserPayload{
			MessageID:   message.MessageID,
			Sender:      senderMap[message.SenderID],
			Content:     message.Content,
			Attachments: make([]contract.MessageAttachment, len(message.Attachments)),
		}

		for i, attachment := range message.Attachments {
			payload.Attachments[i] = mediaMap[attachment.MediaID]
		}

		dtos = append(dtos, ResponseChatMessage{
			MessageType: string(contract.MessageTypeUser),
			Payload:     payload,
			Timestamp:   message.CreatedAt,
		})
	}

	return dtos, nil
}

func (r *GormRepository) GetReviewedMessages(ctx context.Context, problemID uuid.UUID) ([]ResponseChatMessage, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var reviews []database.ProblemReview
	if err := db.WithContext(ctx).
		Table("problem_reviews pr").
		Select("pr.*").
		Joins("LEFT JOIN problem_versions ON problem_versions.problem_version_id = pr.version_id").
		Joins("LEFT JOIN problems ON problems.problem_id = problem_versions.problem_id").
		Where("problems.problem_id = ?", problemID).
		Find(&reviews).Error; err != nil {
		return nil, errors.WrapIf(err, "failed to get reviews ")
	}

	reviewerIDs := make([]uuid.UUID, 0, len(reviews))
	for _, review := range reviews {
		reviewerIDs = append(reviewerIDs, review.ReviewerID)
	}

	reviewerMap, err := r.fetchUsers(ctx, db, reviewerIDs)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to fetch users")
	}

	dtos := make([]ResponseChatMessage, 0, len(reviews))
	for _, review := range reviews {
		dtos = append(dtos, ResponseChatMessage{
			MessageType: string(contract.MessageTypeReviewed),
			Payload: ResponseChatReviewedPayload{
				Reviewer: reviewerMap[review.ReviewerID],
				Decision: review.Decision,
			},
			Timestamp: review.CreatedAt,
		})
	}

	return dtos, nil
}

func (r *GormRepository) GetTestedMessages(ctx context.Context, problemID uuid.UUID) ([]ResponseChatMessage, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var tests []database.ProblemTestResult
	if err := db.WithContext(ctx).
		Table("problem_test_results ptr").
		Select("ptr.*").
		Joins("LEFT JOIN problem_versions ON problem_versions.problem_version_id = ptr.version_id").
		Joins("LEFT JOIN problems ON problems.problem_id = problem_versions.problem_id").
		Where("problems.problem_id = ?", problemID).
		Find(&tests).Error; err != nil {
		return nil, errors.WrapIf(err, "failed to get tests")
	}

	testerIDs := make([]uuid.UUID, 0, len(tests))
	for _, test := range tests {
		testerIDs = append(testerIDs, test.TesterID)
	}

	testerMap, err := r.fetchUsers(ctx, db, testerIDs)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to fetch users")
	}

	dtos := make([]ResponseChatMessage, 0, len(tests))
	for _, test := range tests {
		dtos = append(dtos, ResponseChatMessage{
			MessageType: string(contract.MessageTypeTested),
			Payload: ResponseChatTestedPayload{
				Tester: testerMap[test.TesterID],
				Status: test.Status,
			},
			Timestamp: test.CreatedAt,
		})
	}

	return dtos, nil
}

func (r *GormRepository) GetCompletedMessage(ctx context.Context, problemID uuid.UUID) (*ResponseChatMessage, error) {
	db := database.GetDBFromContext(ctx, r.db)

	type dbModel struct {
		CompletedAt         sql.NullTime
		CompletedByID       uuid.NullUUID
		CompletedByUsername string
	}

	var problem dbModel
	if err := db.WithContext(ctx).
		Table("problems p").
		Select("p.completed_at, p.completed_by AS completed_by_id, users.username AS completed_by_username").
		Joins("LEFT JOIN users ON users.user_id = p.completed_by").
		Where("problem_id = ?", problemID).
		First(&problem).Error; err != nil {
		return nil, errors.WrapIf(err, "failed to get problem")
	}

	if !problem.CompletedAt.Valid || !problem.CompletedByID.Valid {
		return nil, nil
	}

	return &ResponseChatMessage{
		MessageType: string(contract.MessageTypeCompleted),
		Payload: ResponseChatCompletedPayload{
			Completer: contract.MessageUser{
				UserID:   problem.CompletedByID.UUID,
				Username: problem.CompletedByUsername,
			},
		},
		Timestamp: problem.CompletedAt.Time,
	}, nil
}

func (r *GormRepository) fetchUsers(
	ctx context.Context,
	db *gorm.DB,
	userIDs []uuid.UUID,
) (map[uuid.UUID]contract.MessageUser, error) {
	var users []database.User
	if err := db.WithContext(ctx).
		Model(&database.User{}).
		Select("user_id, username").
		Where("user_id IN ?", userIDs).
		Find(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithStack(ErrProblemNotFound)
		}

		return nil, errors.WrapIf(err, "failed to get users")
	}

	userMap := make(map[uuid.UUID]contract.MessageUser, len(users))
	for _, user := range users {
		userMap[user.UserID] = contract.MessageUser{
			UserID:   user.UserID,
			Username: user.Username,
		}
	}

	return userMap, nil
}
