package routes

import (
	authController "payment-gateway-lite/controllers/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")

	// user
	v1.Get("/user", authController.GetUser)

	v1.Post("/signup", authController.Signup)
	v1.Post("/signin", authController.Signin)

	//	oauth
	v1.Post("/auth/google", authController.GoogleOauthLogin)

}
