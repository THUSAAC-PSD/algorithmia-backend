package contest

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb"

	"github.com/labstack/echo/v4"
)

type EndpointParams struct {
	ContestsGroup *echo.Group
}

func NewEndpointParams(
	v1Group *echoweb.V1Group,
) *EndpointParams {
	contests := v1Group.Group.Group("/contests")
	return &EndpointParams{
		ContestsGroup: contests,
	}
}
