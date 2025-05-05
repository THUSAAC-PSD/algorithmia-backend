package upsertproblemdraft

import (
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/problemdraft/shared"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type Endpoint struct {
	*shared.ProblemDraftEndpointParams
}

func NewEndpoint(params *shared.ProblemDraftEndpointParams) *Endpoint {
	return &Endpoint{
		ProblemDraftEndpointParams: params,
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

		if err := e.Validator.StructCtx(ctx.Request().Context(), command); err != nil {
			return httperror.New(http.StatusBadRequest, err.Error()).WithInternal(err)
		}

		response, err := mediatr.Send[*Command, *Response](
			ctx.Request().Context(),
			command,
		)

		if errors.Is(err, ErrInvalidProblemDraftID) {
			return httperror.New(http.StatusUnprocessableEntity, "The problem draft you're trying to update does not exist").
				WithInternal(err)
		} else if errors.Is(err, ErrInvalidProblemDifficultyID) {
			return httperror.New(http.StatusUnprocessableEntity, "The problem difficulty ID you provided doesn't correspond to any valid problem difficulties").WithInternal(err)
		} else if err != nil {
			return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
		}

		return ctx.JSON(http.StatusCreated, response)
	}
}
