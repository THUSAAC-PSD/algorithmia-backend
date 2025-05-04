package context

import (
	"context"

	"github.com/labstack/echo/v4"
)

type key struct{}

// Middleware is a middleware that sets the echo context in the request context.
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(
				c.Request().
					WithContext(context.
						WithValue(
							c.Request().Context(),
							key{},
							c,
						),
					),
			)

			return next(c)
		}
	}
}

func FromContext(ctx context.Context) echo.Context {
	if ctx == nil {
		return nil
	}

	c, ok := ctx.Value(key{}).(echo.Context)
	if !ok {
		return nil
	}

	return c
}
