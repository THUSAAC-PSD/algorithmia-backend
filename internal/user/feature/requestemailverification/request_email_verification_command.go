package requestemailverification

type Command struct {
	Email string `json:"email" validate:"required"`
}
