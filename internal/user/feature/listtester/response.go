package listtester

import "github.com/google/uuid"

type ResponseTester struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
}

type Response struct {
	Testers []ResponseTester `json:"testers"`
}
