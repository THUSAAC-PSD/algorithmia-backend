package shared

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type ProblemDraftEndpointParams struct {
	ProblemDraftsGroup *echo.Group
	Validator          *validator.Validate
}
