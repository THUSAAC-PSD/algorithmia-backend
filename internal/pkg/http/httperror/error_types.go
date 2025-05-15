package httperror

type ErrorType string

const (
	ErrTypeNotAuthenticated             ErrorType = "not_authenticated"
	ErrTypeInvalidCredentials           ErrorType = "invalid_credentials"
	ErrTypeUserAlreadyExists            ErrorType = "user_already_exists"
	ErrTypeInvalidEmailVerificationCode ErrorType = "invalid_email_verification_code"
	ErrTypeRateLimitExceeded            ErrorType = "rate_limit_exceeded"
	ErrTypeIncompleteProblemDraft       ErrorType = "incomplete_problem_draft"
)

func (e ErrorType) String() string {
	return string(e)
}
