package middlewares

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tarcea/go-fiber-blog/initializers"
	"github.com/tarcea/go-fiber-blog/models"
)

func CheckToken(c *fiber.Ctx) error {

	// get the cookie of RequireAuth
	tokenString := c.Cookies("token")
	fmt.Println(tokenString)
	// Decode/validate init
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check the expiration database
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.Status(fiber.ErrUnauthorized.Code)
		}
		// find the user with token sub
		var user models.User
		initializers.DB.First(&user, claims["sub"])
		if user.ID == 0 {
			c.Status(fiber.ErrUnauthorized.Code)
		}
		// attach to locals
		c.Locals("user", user)
		// continue
		return c.Next()
	} else {
		return c.Status(fiber.ErrUnauthorized.Code).JSON(map[string]string{"message": "not authorized"})
	}

}
