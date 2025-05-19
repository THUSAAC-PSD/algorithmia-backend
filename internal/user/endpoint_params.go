package user

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb"

	"github.com/labstack/echo/v4"
)

type EndpointParams struct {
	UsersGroup *echo.Group
	AuthGroup  *echo.Group
}

func NewEndpointParams(
	v1Group *echoweb.V1Group,
) *EndpointParams {
	users := v1Group.Group.Group("/users")
	auth := v1Group.Group.Group("/auth")

	return &EndpointParams{
		UsersGroup: users,
		AuthGroup:  auth,
	}
}
