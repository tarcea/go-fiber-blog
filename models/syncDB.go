package models

import "github.com/tarcea/go-fiber-blog/initializers"

func SyncDb() {
	initializers.DB.AutoMigrate(&Post{})
	initializers.DB.AutoMigrate(&User{})
}
