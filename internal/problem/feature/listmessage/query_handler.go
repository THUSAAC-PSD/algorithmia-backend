package listmessage

import (
	"context"
	"slices"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"

	"emperror.dev/errors"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

var (
	ErrUserNotPartOfRoom = errors.New("user is not part of room")
	ErrProblemNotFound   = errors.New("problem not found")
)

type Query struct {
	ProblemID uuid.UUID `param:"problem_id" validate:"required"`
}

type Repository interface {
	IsUserPartOfRoom(ctx context.Context, problemID uuid.UUID, userID uuid.UUID) (bool, error)
	GetSubmissionMessages(ctx context.Context, problemID uuid.UUID) ([]ResponseChatMessage, error)
	GetUserChatMessages(ctx context.Context, problemID uuid.UUID) ([]ResponseChatMessage, error)
	GetReviewedMessages(ctx context.Context, problemID uuid.UUID) ([]ResponseChatMessage, error)
	GetTestedMessages(ctx context.Context, problemID uuid.UUID) ([]ResponseChatMessage, error)
	GetCompletedMessage(ctx context.Context, problemID uuid.UUID) (*ResponseChatMessage, error)
}

type QueryHandler struct {
	repo         Repository
	authProvider contract.AuthProvider
}

func NewQueryHandler(repo Repository, authProvider contract.AuthProvider) *QueryHandler {
	return &QueryHandler{
		repo:         repo,
		authProvider: authProvider,
	}
}

func (q *QueryHandler) Handle(ctx context.Context, query *Query) (*Response, error) {
	if query == nil {
		return nil, errors.WithStack(customerror.ErrCommandNil)
	}

	user, err := q.authProvider.MustGetUser(ctx)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to get user ID from auth provider")
	}

	if ok, err := q.repo.IsUserPartOfRoom(ctx, query.ProblemID, user.UserID); err != nil {
		return nil, errors.WrapIf(err, "failed to check if user is part of room")
	} else if !ok {
		return nil, errors.WithStack(ErrUserNotPartOfRoom)
	}

	var submissionMessages []ResponseChatMessage
	var userMessages []ResponseChatMessage
	var reviewedMessages []ResponseChatMessage
	var testedMessages []ResponseChatMessage
	var completedMessage *ResponseChatMessage

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		submissionMessages, err = q.repo.GetSubmissionMessages(ctx, query.ProblemID)
		return errors.WrapIf(err, "failed to get submission messages")
	})

	g.Go(func() error {
		var err error
		userMessages, err = q.repo.GetUserChatMessages(ctx, query.ProblemID)
		return errors.WrapIf(err, "failed to get user chat messages")
	})

	g.Go(func() error {
		var err error
		reviewedMessages, err = q.repo.GetReviewedMessages(ctx, query.ProblemID)
		return errors.WrapIf(err, "failed to get reviewed messages")
	})

	g.Go(func() error {
		var err error
		testedMessages, err = q.repo.GetTestedMessages(ctx, query.ProblemID)
		return errors.WrapIf(err, "failed to get tested messages")
	})

	g.Go(func() error {
		var err error
		completedMessage, err = q.repo.GetCompletedMessage(ctx, query.ProblemID)
		return errors.WrapIf(err, "failed to get completed message")
	})

	if err := g.Wait(); err != nil {
		return nil, errors.WrapIf(err, "failed to get messages")
	}

	messages := make([]ResponseChatMessage, 0, len(submissionMessages)+len(userMessages)+len(reviewedMessages)+len(testedMessages)+1)
	messages = append(messages, submissionMessages...)
	messages = append(messages, userMessages...)
	messages = append(messages, reviewedMessages...)
	messages = append(messages, testedMessages...)
	if completedMessage != nil {
		messages = append(messages, *completedMessage)
	}

	slices.SortFunc(messages, func(a, b ResponseChatMessage) int {
		if a.Timestamp.After(b.Timestamp) {
			return -1
		} else if a.Timestamp.Before(b.Timestamp) {
			return 1
		} else {
			return 0
		}
	})

	return &Response{
		Messages: messages,
	}, nil
}
