package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "TradingView Data Service",
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "TradingView Data Service",
			"status":  "running",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Listen(":3333")
}
