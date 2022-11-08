package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html"
	"github.com/tarcea/go-fiber-blog/controllers"
	"github.com/tarcea/go-fiber-blog/initializers"
	"github.com/tarcea/go-fiber-blog/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
	models.SyncDb()
}

func Setup() *fiber.App {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(cors.New())
	app.Static("/", "./public")

	app.Get("/posts", controllers.PostsIndex)
	app.Get("/posts/:id", controllers.PostsView)
	app.Post("/posts", controllers.PostsAdd)
	app.Delete("/posts/:id", controllers.PostsDelete)
	app.Patch("/posts/:id", controllers.PostsUpdate)

	app.Post("/users/signup", controllers.SignUp)
	app.Post("/users/login", controllers.Login)

	return app
}

func main() {

	app := Setup()
	app.Listen(":" + os.Getenv("PORT"))

}
