package helper

import "github.com/go-playground/validator/v10"

func Validate(data interface{}) error {
	validate := validator.New()
	err := validate.Struct(data)
	return err
}
