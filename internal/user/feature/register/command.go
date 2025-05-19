package register

type Command struct {
	Username              string `json:"username"                validate:"required,min=5"`
	Password              string `json:"password"                validate:"required,min=8"`
	Email                 string `json:"email"                   validate:"required"`
	EmailVerificationCode string `json:"email_verification_code" validate:"required"`
}
