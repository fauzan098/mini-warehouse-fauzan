package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Validate(data interface{}) error {
	var errorMessages []string

	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("%s is required", err.Field()))
			case "email":
				errorMessages = append(errorMessages, fmt.Sprintf("%s is not valid email", err.Field()))
			case "min":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be at least %s character long", err.Field(), err.Param()))
			case "max":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be at most %s character long", err.Field(), err.Param()))
			}
		}
		return errors.New("validasi gagal: " + strings.Join(errorMessages, ", "))
	}

	return nil
}

// func joinMessage(errorMessages []string) string {
// 	result := ""
// 	for i, message := range errorMessages {
// 		if i > 0 {
// 			result += ", "
// 		}
// 		result += message
// 	}

// 	return result
// }
