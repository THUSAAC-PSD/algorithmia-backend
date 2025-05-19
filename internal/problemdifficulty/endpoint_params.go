package problemdifficulty

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb"

	"github.com/labstack/echo/v4"
)

type EndpointParams struct {
	ProblemDifficultiesGroup *echo.Group
}

func NewEndpointParams(
	v1Group *echoweb.V1Group,
) *EndpointParams {
	problemDifficulties := v1Group.Group.Group("/problem-difficulties")
	return &EndpointParams{
		ProblemDifficultiesGroup: problemDifficulties,
	}
}
