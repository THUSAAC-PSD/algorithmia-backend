package resetpassword

type Command struct {
	CurrentPassword string `json:"current_password" validate:"omitempty"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}
