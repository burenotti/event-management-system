package main

import (
	_ "github.com/burenotti/rtu-it-lab-recruit/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func main() {
	app := fiber.New()
	app.Get("/docs/*", swagger.HandlerDefault)
	_ = app.Listen(":8001")
}
