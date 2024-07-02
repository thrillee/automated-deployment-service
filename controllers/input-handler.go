package controllers

import "github.com/go-playground/validator/v10"

var validate = validator.New()

type handlers func(interface{}) ControllerResponse

func HandleRequest(data interface{}, h handlers) ControllerResponse {
	if validationErr := validate.Struct(&data); validationErr != nil {
		return ControllerResponse{
			Success: false,
			Message: validationErr.Error(),
		}
	}

	return h(data)
}
