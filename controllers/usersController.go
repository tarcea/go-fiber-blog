package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tarcea/go-fiber-blog/initializers"
	"github.com/tarcea/go-fiber-blog/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.JSON(err)
	}
	// hash the Password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		return c.Status(400).JSON(map[string]string{"message": err.Error()})
	}

	user := models.User{Email: body.Email, Password: string(hash), Username: body.Username}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		return c.Status(400).JSON(map[string]string{"message": result.Error.Error()})
	}
	return c.JSON(user)
}
