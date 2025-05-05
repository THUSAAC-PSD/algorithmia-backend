package echoweb

import (
	"encoding/gob"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/contract"
)

const (
	SessionName    = "session"
	SessionUserKey = "user"
)

func init() {
	gob.Register(contract.AuthUser{})
}
