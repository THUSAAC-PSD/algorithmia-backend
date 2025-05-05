package shared

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type ContestEndpointParams struct {
	ContestsGroup *echo.Group
	Validator     *validator.Validate
}
