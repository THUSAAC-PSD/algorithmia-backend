package shared

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/logger"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type UserEndpointParams struct {
	Logger     logger.Logger
	UsersGroup *echo.Group
	AuthGroup  *echo.Group
	Validator  *validator.Validate
}
