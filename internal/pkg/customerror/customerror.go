package customerror

import "emperror.dev/errors"

var (
	ErrCommandNil       = errors.New("command is nil")
	ErrValidationFailed = errors.New("validation failed")
	ErrNotAuthenticated = errors.New("not authenticated")
	ErrBaseNoPermission = errors.New("no permission")
)

func NewNoPermissionError(permission string) error {
	return errors.Wrap(ErrBaseNoPermission, "missing permission "+permission)
}
