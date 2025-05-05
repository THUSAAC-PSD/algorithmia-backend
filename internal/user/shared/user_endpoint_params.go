package shared

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type UserEndpointParams struct {
	UsersGroup *echo.Group
	AuthGroup  *echo.Group
	Validator  *validator.Validate
}
