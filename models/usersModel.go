package models

import (
	"github.com/tarcea/go-fiber-blog/initializers"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string  `json:"username"`
	Email    string  `gorm:"unique" json:"email"`
	Password string  `json:"password"`
	Posts    []*Post `json:"posts"`
}

type ResponseUser struct {
	Username string  `json:"username"`
	Email    User    `json:"email"`
	Posts    []*Post `json:"posts"`
}

func GetUserById(userId int) *User {
	var user *User
	initializers.DB.First(&user, userId)
	return user
}
