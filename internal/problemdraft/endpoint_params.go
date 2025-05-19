package problemdraft

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb"

	"github.com/labstack/echo/v4"
)

type EndpointParams struct {
	ProblemDraftsGroup *echo.Group
}

func NewEndpointParams(
	v1Group *echoweb.V1Group,
) *EndpointParams {
	problemDrafts := v1Group.Group.Group("/problem-drafts")
	return &EndpointParams{
		ProblemDraftsGroup: problemDrafts,
	}
}
