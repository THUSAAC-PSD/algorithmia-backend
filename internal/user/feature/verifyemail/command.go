package verifyemail

type Command struct {
	Email string `json:"email" validate:"required,email"`
	Token string `json:"token" validate:"required"`
}
