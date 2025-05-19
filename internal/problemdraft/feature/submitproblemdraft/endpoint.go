package submitproblemdraft

import (
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
)

type Endpoint struct {
	*problemdraft.EndpointParams
	handler *CommandHandler
}

func NewEndpoint(params *problemdraft.EndpointParams, handler *CommandHandler) *Endpoint {
	return &Endpoint{
		EndpointParams: params,
		handler:        handler,
	}
}

func (e *Endpoint) MapEndpoint() {
	e.ProblemDraftsGroup.POST("/:problem_draft_id/submit", e.handle())
}

func (e *Endpoint) handle() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		command := &Command{}
		if err := ctx.Bind(command); err != nil {
			return httperror.New(http.StatusBadRequest, "Invalid request format")
		}

		if err := ctx.Validate(command); err != nil {
			return err
		}

		response, err := e.handler.Handle(ctx.Request().Context(), command)
		if errors.Is(err, ErrProblemDraftNotFound) {
			return httperror.New(http.StatusNotFound, "The problem draft you're trying to submit does not exist").
				WithInternal(err)
		} else if errors.Is(err, ErrProblemDraftNotActive) {
			return httperror.New(http.StatusUnprocessableEntity, "The problem draft you're trying to submit is not active")
		} else if errors.Is(err, ErrProblemDoesntNeedRevision) {
			return httperror.New(http.StatusUnprocessableEntity, "The problem draft you're trying to submit does not need revision")
		} else if errors.Is(err, ErrNotCreator) {
			return httperror.New(http.StatusForbidden, "You are not the creator of this problem draft")
		} else if errors.Is(err, ErrMissingProblemDifficulty) {
			return httperror.New(http.StatusUnprocessableEntity, "The problem draft you're trying to submit is missing its problem difficulty").
				WithType(httperror.ErrTypeIncompleteProblemDraft)
		} else if errors.Is(err, ErrContestNotFound) {
			return httperror.New(http.StatusUnprocessableEntity, "The contest you're trying to submit to does not exist")
		} else if err != nil {
			return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
		}

		return ctx.JSON(http.StatusOK, response)
	}
}
