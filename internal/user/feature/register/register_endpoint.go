package register

import (
	"fmt"
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/user/shared"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type Endpoint struct {
	*shared.UserEndpointParams
}

func NewEndpoint(params *shared.UserEndpointParams) *Endpoint {
	return &Endpoint{
		UserEndpointParams: params,
	}
}

func (e *Endpoint) MapEndpoint() {
	e.AuthGroup.POST("/register", e.handle())
}

func (e *Endpoint) handle() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		command := &Command{}
		if err := ctx.Bind(command); err != nil {
			return httperror.New(http.StatusBadRequest, 100, "invalid request")
		}

		if err := e.Validator.StructCtx(ctx.Request().Context(), command); err != nil {
			return httperror.New(http.StatusBadRequest, 100, fmt.Sprintf("validation error: %s", err.Error()))
		}

		result, err := mediatr.Send[*Command, *Response](
			ctx.Request().Context(),
			command,
		)

		if errors.Is(err, ErrUserAlreadyExists) {
			return httperror.New(http.StatusConflict, 101, "user already exists")
		} else if errors.Is(err, ErrInvalidEmailVerificationCode) {
			return httperror.New(http.StatusBadRequest, 102, "invalid email verification code")
		} else if err != nil {
			return httperror.New(http.StatusInternalServerError, 200, fmt.Sprintf("error in sending command: %s", err.Error()))
		}

		return ctx.JSON(http.StatusCreated, result)
	}
}
