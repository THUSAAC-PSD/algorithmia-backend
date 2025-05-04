package context

import (
	"context"

	"github.com/labstack/echo/v4"
)

// Middleware is a middleware that sets the echo context in the request context.
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(
				c.Request().
					WithContext(context.
						WithValue(
							c.Request().Context(),
							"echo.context",
							c,
						),
					),
			)

			return next(c)
		}
	}
}
