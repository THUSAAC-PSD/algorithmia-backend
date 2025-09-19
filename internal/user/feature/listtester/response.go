package listtester

import "github.com/google/uuid"

type ResponseTester struct {
	UserID      uuid.UUID `json:"user_id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
}

type Response struct {
	Testers []ResponseTester `json:"testers"`
}
