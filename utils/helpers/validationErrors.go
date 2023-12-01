package helpers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	ValidationErrorResponse struct {
		FailedField string      `json:"failed_field"`
		Tag         string      `json:"tag"`
		Value       interface{} `json:"value"`
	}
)

func ErrorValidation(c *fiber.Ctx, err error) error {
	if errValidation, ok := err.(validator.ValidationErrors); ok {
		var errors []ValidationErrorResponse
		for _, fieldError := range errValidation {
			var errorResponse ValidationErrorResponse
			errorResponse.FailedField = fieldError.Field()
			errorResponse.Tag = fieldError.Tag()
			errorResponse.Value = fieldError.Value()
			//errorResponse.Error = map[string]interface{}{
			//	"params":      fieldError.Param(),
			//	"field":       fieldError.Field(),
			//	"tag":         fieldError.Tag(),
			//	"error":       fieldError.Error(),
			//	"actual":      fieldError.ActualTag(),
			//	"kind":        fieldError.Kind().String(),
			//	"namespace":   fieldError.Namespace(),
			//	"structfield": fieldError.StructField(),
			//	"type":        fieldError.Type(),
			//}

			errors = append(errors, errorResponse)
		}
		return c.Status(fiber.StatusBadRequest).JSON(&ErrorResponse{
			Success: false,
			Message: MessageBadRequest,
			Errors:  errors,
		})
	} else {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

}
