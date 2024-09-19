package utils

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func InitValidator() {
	validate = validator.New()
}

// Validate validates the given struct
func Validate(data interface{}) error {
	return validate.Struct(data)
}