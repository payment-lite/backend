package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/viper"
	"payment-gateway-lite/database"
	"payment-gateway-lite/routes"
	"payment-gateway-lite/utils/helpers"
)

func main() {
	helpers.LoadENV()
	database.InitDB()

	app := fiber.New(fiber.Config{
		AppName: viper.GetString("APP.NAME"),
		Prefork: true,
	})
	app.Use(logger.New())

	routes.SetupRoutes(app)
	app.Use(helpers.FiberErrorHandler)

	err := app.Listen(fmt.Sprintf("%s:%s", viper.GetString("APP.HOST"), viper.GetString("APP.PORT")))
	if err != nil {
		panic(spew.Sprintf("failed to start server: %v", err))
	}
}
