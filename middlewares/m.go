package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func M(c *fiber.Ctx) error {

	fmt.Println("gogo", "123")

	return c.Next()
}
