package main

import (
	"gptube/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	app.Get("/", routes.HomeHandler)
	app.Post("/YT/pre-analysis", routes.YoutubePreAnalysisHandler)

	app.Listen(":8000")
}
