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

func (r *GormRepository) GetSubmissionMessages(ctx context.Context, problemID uuid.UUID) ([]ResponseChatMessage, error) {
	db := database.GetDBFromContext(ctx, r.db)

	var problem database.Problem
	if err := db.WithContext(ctx).
		Preload("Creator").
		Where("problem_id = ?", problemID).
		First(&problem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithStack(ErrProblemNotFound)
		}

		return nil, errors.WrapIf(err, "failed to get problem for submission history")
	}

	var versions []database.ProblemVersion
	if err := db.WithContext(ctx).
		Preload("SubmittedByUser").
		Preload("Details").
		Preload("Examples").
		Where("problem_id = ?", problemID).
		Order("created_at ASC").
		Find(&versions).Error; err != nil {
		return nil, errors.WrapIf(err, "failed to get problem versions")
	}

	messages := make([]ResponseChatMessage, 0, len(versions))
	for i, version := range versions {
		user := contract.MessageUser{
			UserID:   version.SubmittedBy,
			Username: version.SubmittedByUser.Username,
		}

		if user.UserID == uuid.Nil {
			user.UserID = problem.CreatorID
			if problem.Creator.UserID != uuid.Nil {
				user.Username = problem.Creator.Username
			}
		}

		messageType := contract.MessageTypeSubmitted
		payload := ResponseChatSubmittedPayload{Submitter: user}
		if i > 0 {
			messageType = contract.MessageTypeEdited
			// Compute changed fields and diffs compared to previous version
			changed, diffs := computeChanged(&versions[i-1], &version)
			if len(changed) > 0 {
				payload.ChangedFields = changed
			}
			if len(diffs) > 0 {
				payload.Diffs = diffs
			}
		}

		messages = append(messages, ResponseChatMessage{
			MessageType: string(messageType),
			Payload:     payload,
			Timestamp:   version.CreatedAt,
		})
	}

	return messages, nil
}

// computeChangedFields compares two consecutive versions and returns a concise list of
// top-level fields that changed. For multi-language details, if any language differs,
// we mark the corresponding field as changed.
func computeChanged(prev *database.ProblemVersion, curr *database.ProblemVersion) ([]string, map[string]FieldDiff) {
	changedSet := make(map[string]struct{})
	diffs := make(map[string]FieldDiff)

	// Difficulty
	if prev.ProblemDifficultyID != curr.ProblemDifficultyID {
		changedSet["difficulty"] = struct{}{}
		diffs["difficulty"] = FieldDiff{Before: prev.ProblemDifficultyID.String(), After: curr.ProblemDifficultyID.String()}
	}

	// Build maps of details by language
	prevDetails := make(map[string]database.ProblemVersionDetail, len(prev.Details))
	for _, d := range prev.Details {
		prevDetails[d.Language] = d
	}
	currDetails := make(map[string]database.ProblemVersionDetail, len(curr.Details))
	for _, d := range curr.Details {
		currDetails[d.Language] = d
	}

	// Union of languages
	langSet := make(map[string]struct{})
	for l := range prevDetails {
		langSet[l] = struct{}{}
	}
	for l := range currDetails {
		langSet[l] = struct{}{}
	}

	// Compare each field across languages; pick en-US when available otherwise the first language
	for l := range langSet {
		pd, pok := prevDetails[l]
		cd, cok := currDetails[l]
		if !pok || !cok {
			// A language was added or removed; mark all textual fields
			changedSet["title"] = struct{}{}
			changedSet["background"] = struct{}{}
			changedSet["statement"] = struct{}{}
			changedSet["input_format"] = struct{}{}
			changedSet["output_format"] = struct{}{}
			changedSet["note"] = struct{}{}
			// Use available values for before/after (empty when missing)
			diffs["title"] = FieldDiff{Before: pd.Title, After: cd.Title}
			diffs["background"] = FieldDiff{Before: pd.Background, After: cd.Background}
			diffs["statement"] = FieldDiff{Before: pd.Statement, After: cd.Statement}
			diffs["input_format"] = FieldDiff{Before: pd.InputFormat, After: cd.InputFormat}
			diffs["output_format"] = FieldDiff{Before: pd.OutputFormat, After: cd.OutputFormat}
			diffs["note"] = FieldDiff{Before: pd.Note, After: cd.Note}
			continue
		}
		if pd.Title != cd.Title {
			changedSet["title"] = struct{}{}
			diffs["title"] = FieldDiff{Before: pd.Title, After: cd.Title}
		}
		if pd.Background != cd.Background {
			changedSet["background"] = struct{}{}
			diffs["background"] = FieldDiff{Before: pd.Background, After: cd.Background}
		}
		if pd.Statement != cd.Statement {
			changedSet["statement"] = struct{}{}
			diffs["statement"] = FieldDiff{Before: pd.Statement, After: cd.Statement}
		}
		if pd.InputFormat != cd.InputFormat {
			changedSet["input_format"] = struct{}{}
			diffs["input_format"] = FieldDiff{Before: pd.InputFormat, After: cd.InputFormat}
		}
		if pd.OutputFormat != cd.OutputFormat {
			changedSet["output_format"] = struct{}{}
			diffs["output_format"] = FieldDiff{Before: pd.OutputFormat, After: cd.OutputFormat}
		}
		if pd.Note != cd.Note {
			changedSet["note"] = struct{}{}
			diffs["note"] = FieldDiff{Before: pd.Note, After: cd.Note}
		}
	}

	// Examples: simple len or pairwise comparison
	if len(prev.Examples) != len(curr.Examples) {
		changedSet["examples"] = struct{}{}
	} else {
		for i := range prev.Examples {
			if prev.Examples[i].Input != curr.Examples[i].Input || prev.Examples[i].Output != curr.Examples[i].Output {
				changedSet["examples"] = struct{}{}
				break
			}
		}
	}
	// Include examples as a simple text blob (line-per-example) when changed
	if _, ok := changedSet["examples"]; ok {
		before := ""
		after := ""
		for i, e := range prev.Examples {
			before += "# Example " + itoa(i+1) + "\n" + "Input:\n" + e.Input + "\nOutput:\n" + e.Output + "\n\n"
		}
		for i, e := range curr.Examples {
			after += "# Example " + itoa(i+1) + "\n" + "Input:\n" + e.Input + "\nOutput:\n" + e.Output + "\n\n"
		}
		diffs["examples"] = FieldDiff{Before: before, After: after}
	}

	// Convert set to slice with a stable order
	order := []string{"title", "background", "statement", "input_format", "output_format", "note", "examples", "difficulty"}
	out := make([]string, 0, len(changedSet))
	for _, key := range order {
		if _, ok := changedSet[key]; ok {
			out = append(out, key)
		}
	}
	return out, diffs
}

// tiny helper without importing strconv everywhere
func itoa(i int) string {
	// simple conversion (fast enough and self-contained)
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	buf := [20]byte{}
	bp := len(buf)
	for i > 0 {
		bp--
		buf[bp] = byte('0' + (i % 10))
		i /= 10
	}
	if neg {
		bp--
		buf[bp] = '-'
	}
	return string(buf[bp:])
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
