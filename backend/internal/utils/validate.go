package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

var validate *validator.Validate

// InitValidator initializes the validator instance
func InitValidator() {
	LogInfo("Initializing validator", nil)
	validate = validator.New()
}

// Validate validates the given struct
func Validate(data interface{}) error {
	err := validate.Struct(data)
	if err != nil {
		LogWarn("Validation failed", logrus.Fields{
			"errors": err.Error(),
		})
		return err
	}

	LogInfo("Validation successful", nil)
	return nil
}
