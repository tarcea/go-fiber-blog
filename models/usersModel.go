package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Posts    []*Post `json:"posts"`
}

type ResponseUser struct {
	Username string  `json:"username"`
	Email    User    `json:"email"`
	Posts    []*Post `json:"posts"`
}
