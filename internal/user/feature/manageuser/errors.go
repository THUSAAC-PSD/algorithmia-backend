package manageuser

import "emperror.dev/errors"

var (
	ErrUserNotFound               = errors.New("user not found")
	ErrRoleNotFound               = errors.New("role not found")
	ErrEmailAlreadyExists         = errors.New("email already exists")
	ErrUsernameAlreadyExists      = errors.New("username already exists")
	ErrRolesRequired              = errors.New("at least one role is required")
	ErrCannotRemoveOwnSuperAdmin  = errors.New("cannot remove super admin role from your own account")
	ErrCannotDeleteSelf           = errors.New("cannot delete your own account")
	ErrCannotDeleteLastSuperAdmin = errors.New("cannot remove the last super admin")
)
