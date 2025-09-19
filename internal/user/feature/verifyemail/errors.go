package verifyemail

import "emperror.dev/errors"

var (
	ErrInvalidOrExpiredToken = errors.New("invalid or expired verification token")
	ErrUserAlreadyExists     = errors.New("user with this email already exists")
)
