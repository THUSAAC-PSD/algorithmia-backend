package httperror

import (
	stderrors "errors"
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/customerror"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
)

type HTTPError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Internal   error  `json:"-"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func (e *HTTPError) Unwrap() error {
	return e.Internal
}

func New(statusCode int, code int, message string) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
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
	var response map[string]interface{}
	statusCode := getStatusCode(err)

	var customErr *HTTPError
	if errors.As(err, &customErr) {
		response = map[string]interface{}{
			"code":    customErr.Code,
			"message": customErr.Message,
		}
	} else {
		var echoErr *echo.HTTPError
		if errors.As(err, &echoErr) {
			statusCode = echoErr.Code
			if msg, ok := echoErr.Message.(string); ok {
				response = map[string]interface{}{"message": msg}
			} else {
				response = map[string]interface{}{"message": echoErr.Message}
			}
		} else if errors.Is(err, customerror.ErrValidationFailed) {
			response = map[string]interface{}{
				"code":    100,
				"message": err.Error(),
			}
		} else if errors.Is(err, customerror.ErrCommandNil) {
			response = map[string]interface{}{
				"code":    200,
				"message": err.Error(),
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
