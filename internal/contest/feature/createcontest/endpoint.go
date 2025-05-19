package createcontest

import (
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/contest"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
)

type Endpoint struct {
	*contest.EndpointParams
	handler *CommandHandler
}

func NewEndpoint(params *contest.EndpointParams, handler *CommandHandler) *Endpoint {
	return &Endpoint{
		EndpointParams: params,
		handler:        handler,
	}
}

func (e *Endpoint) MapEndpoint() {
	e.ContestsGroup.POST("", e.handle())
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

		if errors.Is(err, ErrInvalidProblemCount) {
			return httperror.New(http.StatusUnprocessableEntity, "Problem count must be greater than 0")
		} else if errors.Is(err, ErrProblemCountRangeFlipped) {
			return httperror.New(http.StatusUnprocessableEntity, "Problem count range flipped. Max should be greater than min")
		} else if errors.Is(err, ErrInvalidDeadlineDatetime) {
			return httperror.New(http.StatusUnprocessableEntity, "Deadline must be in the future")
		} else if err != nil {
			return httperror.New(http.StatusInternalServerError, err.Error()).WithInternal(err)
		}

		return ctx.JSON(http.StatusCreated, response)
	}
}
