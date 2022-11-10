package controllers

import (
	"errors"
	"fmt"
	"os"
	"strconv"
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

	// check if the email was already used
	var u = new(models.User)
	initializers.DB.First(&user, "email = ?", body.Email)

	if u.ID == 0 {
		return c.Status(400).JSON(map[string]string{"message": "email already used"})
	}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		fmt.Println("rer", result.Error.Error())
		return c.Status(400).JSON(map[string]string{"message": result.Error.Error()})
	}

	// create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	tokenError := errors.New("failed to create token")

	if err != nil {
		return c.Status(400).JSON(map[string]string{"message": tokenError.Error()})
	}

	c.Set("access-control-expose-headers", "Set-Cookie")

	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: time.Now().Add(24 * time.Hour),
	})

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

	c.Set("access-control-expose-headers", "Set-Cookie")

	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: time.Now().Add(24 * time.Hour),
	})

	return c.Status(200).JSON(map[string]string{"message": "Ok"})

}

func Validate(c *fiber.Ctx) error {
	user := c.Locals("user")
	email := user.(models.User).Email
	i := user.(models.User).ID
	id := strconv.FormatInt(int64(i), 10)
	return c.JSON(map[string]string{"user": email, "uid": id})
}

func Logout(c *fiber.Ctx) error {
	c.Set("access-control-expose-headers", "Set-Cookie")
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Expires: time.Now().Add(-(time.Hour * 2)),
	})

	return c.Status(200).JSON(map[string]string{"message": "Ok, clear"})
}
