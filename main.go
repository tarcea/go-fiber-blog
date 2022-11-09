package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/tarcea/go-fiber-blog/controllers"
	"github.com/tarcea/go-fiber-blog/initializers"
	"github.com/tarcea/go-fiber-blog/middlewares"
	"github.com/tarcea/go-fiber-blog/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
	models.SyncDb()
}

func Setup() *fiber.App {

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CLIENT_URL"),
		AllowCredentials: true,
	}))

	app.Get("/posts", controllers.PostsIndex)
	app.Get("/posts/:id", controllers.PostsView)
	app.Post("/posts", controllers.PostsAdd)
	app.Delete("/posts/:id", controllers.PostsDelete)
	app.Patch("/posts/:id", controllers.PostsUpdate)

	app.Post("/users/signup", controllers.SignUp)
	app.Post("/users/login", controllers.Login)
	app.Post("/users/logout", controllers.Logout)

	app.Get("/validate", middlewares.CheckToken, controllers.Validate)

	return app
}

func main() {
	app := Setup()
	app.Listen(":" + os.Getenv("PORT"))
}
