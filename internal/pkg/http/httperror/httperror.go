package httperror

import (
	stderrors "errors"
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
)

type HTTPError struct {
	Type       ErrorType `json:"type,omitempty"`
	Message    string    `json:"message"`
	StatusCode int       `json:"-"`
	Internal   error     `json:"-"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func (e *HTTPError) Unwrap() error {
	return e.Internal
}

func New(statusCode int, message string) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Message:    message,
	}
}

func (e *HTTPError) WithType(t ErrorType) *HTTPError {
	e.Type = t
	return e
}

func (e *HTTPError) WithInternal(err error) *HTTPError {
	e.Internal = err
	return e
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	var he *HTTPError
	if errors.As(err, &he) {
		return he.StatusCode
	}

	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		return echoErr.Code
	}

	return http.StatusInternalServerError
}

func getRootCause(err error) error {
	for err != nil {
		unwrapped := stderrors.Unwrap(err)
		if unwrapped == nil {
			break
		}
		err = unwrapped
	}

	return err
}

// Handler processes errors for Echo
func Handler(err error, c echo.Context) {
	var response any
	statusCode := getStatusCode(err)

	var httpErr *HTTPError
	if httpErr = mapCustomErrorToHTTPError(err); httpErr != nil {
		response = httpErr
		statusCode = httpErr.StatusCode
	} else if errors.As(err, &httpErr); httpErr != nil {
		response = httpErr
		statusCode = httpErr.StatusCode
	} else {
		var echoErr *echo.HTTPError
		if errors.As(err, &echoErr) {
			statusCode = echoErr.Code

			if msg, ok := echoErr.Message.(string); ok {
				response = map[string]interface{}{"message": msg}
			} else {
				response = map[string]interface{}{"message": echoErr.Message}
			}
		} else {
			rootErr := getRootCause(err)
			if rootErr != nil {
				response = map[string]interface{}{"message": rootErr.Error()}
			}
		}
	}

	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(statusCode)
		} else {
			err = c.JSON(statusCode, response)
		}

		if err != nil {
			c.Echo().Logger.Error(err)
		}
	}
}

func mapCustomErrorToHTTPError(err error) *HTTPError {
	if errors.Is(err, customerror.ErrValidationFailed) {
		return New(http.StatusBadRequest, err.Error()).WithInternal(err)
	} else if errors.Is(err, customerror.ErrCommandNil) {
		return New(http.StatusBadRequest, err.Error()).WithInternal(err)
	} else if errors.Is(err, customerror.ErrNotAuthenticated) {
		return New(http.StatusUnauthorized, err.Error()).WithInternal(err).WithType(ErrTypeNotAuthenticated)
	} else if errors.Is(err, customerror.ErrBaseNoPermission) {
		return New(http.StatusForbidden, err.Error()).WithInternal(err).WithType(ErrTypeNoPermission)
	}

	return nil
}
