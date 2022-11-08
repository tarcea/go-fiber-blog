package controllers

import (
	"errors"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
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

func Login(c *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	credentialsError := errors.New("invalid credentials")
	tokenError := errors.New("failed to create token")

	if err := c.BodyParser(&body); err != nil {
		return c.JSON(err)
	}

	var user = new(models.User)
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		return c.Status(400).JSON(map[string]string{"message": credentialsError.Error()})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		return c.Status(400).JSON(map[string]string{"message": credentialsError.Error()})
	}

	// create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return c.Status(400).JSON(map[string]string{"message": tokenError.Error()})
	}

	// set a cookie and send the cookie back
	// cookie := new(fiber.Cookie)
	// cookie.Name = "token"
	// cookie.Value = tokenString
	// cookie.Expires = time.Now().Add(24 * time.Hour)
	// cookie.Secure = true
	// cookie.HTTPOnly = true
	// cookie.SameSite = "lax"
	// c.Cookie(cookie)

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: "lax",
	})

	return c.Status(200).JSON(map[string]string{"message": "Ok"})

}
