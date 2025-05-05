package customerror

import "emperror.dev/errors"

var (
	ErrCommandNil       = errors.New("command is nil")
	ErrValidationFailed = errors.New("validation failed")
	ErrNotAuthenticated = errors.New("not authenticated")
)
