package custom_validator

import (
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/DrReMain/cyber-ecosystem-server/api/admin/internal/helper/usual_err"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) Validate(_ *http.Request, data any) error {
	if err := v.validate.Struct(data); err != nil {
		return usual_err.HTTPBadRequest(err.Error())
	}
	return nil
}
