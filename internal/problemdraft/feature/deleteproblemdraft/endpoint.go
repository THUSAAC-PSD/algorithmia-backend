package deleteproblemdraft

import (
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft"

	"emperror.dev/errors"
	"github.com/google/uuid"
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
	e.ProblemDraftsGroup.DELETE("/:problem_draft_id", e.handle())
}

func (e *Endpoint) handle() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		command := &Command{}
		if err := ctx.Bind(command); err != nil {
			return httperror.New(http.StatusBadRequest, "Invalid request format")
		}

		// Manually parse UUID path param because Echo binder does not always populate uuid.UUID fields
		if command.ProblemDraftID == uuid.Nil {
			if raw := ctx.Param("problem_draft_id"); raw != "" {
				id, err := uuid.Parse(raw)
				if err != nil {
					return httperror.New(http.StatusBadRequest, "Invalid problem_draft_id path parameter")
				}
				command.ProblemDraftID = id
			} else {
				return httperror.New(http.StatusBadRequest, "Missing problem_draft_id path parameter")
			}
		}

		if err := ctx.Validate(command); err != nil {
			return err
		}

		err := e.handler.Handle(ctx.Request().Context(), command)
		if errors.Is(err, ErrProblemDraftNotFound) {
			return httperror.New(http.StatusNotFound, "The problem draft does not exist").
				WithInternal(err)
		} else if err != nil {
			return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
		}

		return ctx.NoContent(http.StatusNoContent)
	}
}
