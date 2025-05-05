package shared

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type ContestEndpointParams struct {
	Logger        logger.Logger
	ContestsGroup *echo.Group
	Validator     *validator.Validate
}
