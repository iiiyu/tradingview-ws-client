package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func HandleMessage(c *fiber.Ctx) error {
	var msg []byte
	if err := c.BodyParser(&msg); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid message format",
		})
	}

	fmt.Println(string(msg))
	return c.SendStatus(200)
}
