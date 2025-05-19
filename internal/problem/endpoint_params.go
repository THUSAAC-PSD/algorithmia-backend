package problem

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb"

	"github.com/labstack/echo/v4"
)

type EndpointParams struct {
	ProblemsGroup *echo.Group
}

func NewEndpointParams(
	v1Group *echoweb.V1Group,
) *EndpointParams {
	problems := v1Group.Group.Group("/problems")
	return &EndpointParams{
		ProblemsGroup: problems,
	}
}
