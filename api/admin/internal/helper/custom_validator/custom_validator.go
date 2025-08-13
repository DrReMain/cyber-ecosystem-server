package custom_validator

import (
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/DrReMain/cyber-ecosystem-server/pkg/errorc"
	"github.com/DrReMain/cyber-ecosystem-server/pkg/msgc"
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
		return errorc.NewHTTPBadRequest(msgc.VALIDATE_FAILURE, err.Error())
	}
	return nil
}
