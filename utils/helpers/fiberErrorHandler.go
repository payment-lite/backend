package helpers

import "github.com/gofiber/fiber/v2"

func FiberErrorHandler(c *fiber.Ctx) error {
	err := c.Next()
	if err != nil {
		if fiberError, ok := err.(*fiber.Error); ok {
			return c.Status(fiberError.Code).JSON(&ErrorResponse{
				Success: false,
				Message: fiberError.Message,
			})
		}
	}
	return c.Status(fiber.StatusInternalServerError).JSON(&ErrorResponse{
		Success: false,
		Message: "Something went wrong",
		Errors:  err.Error(),
	})
}
