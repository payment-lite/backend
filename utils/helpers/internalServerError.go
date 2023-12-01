package helpers

import "github.com/gofiber/fiber/v2"

func InternalServerError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(&ErrorResponse{
		Success: false,
		Message: MessageUnauthorized,
		Errors:  err.Error(),
	})

}
