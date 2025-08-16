package services

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate

func init() {
    Validator = validator.New()
		err := Validator.RegisterValidation("oneof", ValidateStringInSlice)
	if err != nil {
		panic(err)
	}
}


func ValidateStringInSlice(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	allowedValues := strings.Split(fl.Param(), " ")

	for _, v := range allowedValues {
		if v == value {
			return true
		}
	}

	return false
}
