package shared

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type ProblemEndpointParams struct {
	ProblemsGroup *echo.Group
	Validator     *validator.Validate
}
