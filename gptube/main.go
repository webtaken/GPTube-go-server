package main

import (
	"gptube/router"

	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()
	router.SetupRoutes(app)
	app.Listen(":8000")
}
