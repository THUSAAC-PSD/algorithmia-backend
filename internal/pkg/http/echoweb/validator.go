package echoweb

import (
	"net/http"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/httperror"

	"github.com/go-playground/validator"
)

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return httperror.New(http.StatusBadRequest, err.Error())
	}
	return nil
}
