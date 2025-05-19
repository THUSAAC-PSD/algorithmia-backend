package upsertproblemdraft

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
	e.ProblemDraftsGroup.PUT("", e.handle())
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
		if errors.Is(err, ErrInvalidProblemDraftID) {
			return httperror.New(http.StatusUnprocessableEntity, "The problem draft you're trying to update does not exist").
				WithInternal(err)
		} else if errors.Is(err, ErrInvalidProblemDifficultyID) {
			return httperror.New(http.StatusUnprocessableEntity, "The problem difficulty ID you provided doesn't correspond to any valid problem difficulties").WithInternal(err)
		} else if errors.Is(err, ErrNotCreatorOrInactive) {
			return httperror.New(http.StatusForbidden, "You are not the creator of this problem draft or the draft is inactive").WithInternal(err)
		} else if err != nil {
			return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
		}

		return ctx.JSON(http.StatusCreated, response)
	}
}
