package controllers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/tarcea/go-fiber-blog/initializers"
	"github.com/tarcea/go-fiber-blog/models"
)

func PostsIndex(c *fiber.Ctx) error {
	err := errors.New("there are no posts")
	responseChannel := make(chan *[]models.Post)
	go models.GetAllPublishedPosts(responseChannel)
	posts := <-responseChannel

	if len(*posts) == 0 {
		fmt.Println(err)
		return c.Status(400).JSON(map[string]string{"message": err.Error()})
	}

	d := *posts
	var r []models.ResponsePost

	for _, p := range d {
		var resp models.ResponsePost
		resp.ID = p.ID
		resp.Body = p.Body
		resp.Published = p.Published
		resp.Title = p.Title
		resp.UserId = p.UserId
		resp.Author = p.User.Username
		r = append(r, resp)
	}

	return c.JSON(r)
}

func PostsView(c *fiber.Ctx) error {
	postId := c.Params("id")

	err := errors.New("record not found")
	responseChannel := make(chan *models.Post)
	go models.GetPostById(responseChannel, postId)
	post := <-responseChannel

	if post.ID == 0 {
		return c.Status(400).JSON(map[string]string{"message": err.Error()})
	}

	p := *post
	var resp models.ResponsePost

	resp.ID = p.ID
	resp.Body = p.Body
	resp.Published = p.Published
	resp.Title = p.Title
	resp.UserId = p.UserId
	resp.Author = p.User.Username

	return c.JSON(resp)
}

func PostsAdd(c *fiber.Ctx) error {

	post := new(models.Post)

	if err := c.BodyParser(post); err != nil {
		return err
	}

	initializers.DB.Create(&post)

	return c.JSON(post)
}

func PostsDelete(c *fiber.Ctx) error {
	postId := c.Params("id")

	err := errors.New("record not found")
	responseChannel := make(chan *models.Post)
	go models.GetPostById(responseChannel, postId)
	post := <-responseChannel

	if post.ID == 0 {
		return c.Status(400).JSON(map[string]string{"message": err.Error()})
	}

	models.DeletePost(postId)

	return c.JSON(post)
}
