package main

import (
	"flag"
	_ "github.com/burenotti/rtu-it-lab-recruit/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "0.0.0.0:80", "Documentation address")
	flag.Parse()
	app := fiber.New()
	app.Get("/*", swagger.HandlerDefault)
	_ = app.Listen(addr)
}
