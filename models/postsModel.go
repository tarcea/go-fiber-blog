package models

import (
	"github.com/tarcea/go-fiber-blog/initializers"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title     string `json:"title"`
	Body      string `json:"body"`
	Published bool   `json:"published"`
	UserId    uint   `json:"userId"`
	User      User   `gorm:"foreignKey:ID;references:user_id" json:"user"`
}

type ResponsePost struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Published bool   `json:"published"`
	UserId    uint   `json:"userId"`
	Author    string `json:"author"`
}

func GetAllPublishedPosts(res chan *[]Post) {
	var posts []Post

	initializers.DB.
		Model(&Post{}).
		Order("created_at desc").
		Where("published = ?", "true").
		Preload("User").
		Find(&posts)

	res <- &posts
}

func GetAllPublishedPostsByUser(res chan *[]Post, userId int) {
	var posts []Post
	initializers.DB.
		Where("published = ? AND user_id = ?", "true", userId).
		Order("created_at desc").
		Preload("User").
		Find(&posts)

	res <- &posts
}

func GetPostById(res chan *Post, postId string) {
	type response struct {
		ID     string
		Title  string
		Body   string
		Author string `json:"author"`
	}
	var resp response
	var post *Post

	initializers.DB.
		Preload("User").
		First(&post, postId)

	resp.Body = post.Body

	res <- post

}

func DeletePost(postId string) {
	initializers.DB.Delete(&Post{}, postId)
}
