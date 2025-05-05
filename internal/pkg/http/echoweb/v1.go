package echoweb

import "github.com/labstack/echo/v4"

type V1Group struct {
	*echo.Group
}

func NewV1Group(e *echo.Echo) *V1Group {
	v1 := e.Group("/api/v1")
	return &V1Group{
		Group: v1,
	}
}
