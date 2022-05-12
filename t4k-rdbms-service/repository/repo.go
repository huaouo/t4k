package repository

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

type Account struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `gorm:"index:idx_name,unique"`
	Password  string
	CreatedAt time.Time
}

type Follow struct {
	UserID    uint      `gorm:"index:idx_user_id"`
	ToUserID  uint      `gorm:"index:idx_to_user_id"`
	CreatedAt time.Time `gorm:"index:idx_created_at"`
}

type Video struct {
	ID        uint `gorm:"primary_key"`
	PlayURL   string
	CoverURL  string
	CreatedAt time.Time `gorm:"index:idx_created_at"`
}

type Favorite struct {
	UserID    uint      `gorm:"index:idx_user_id"`
	VideoID   uint      `gorm:"index:idx_video_id"`
	CreatedAt time.Time `gorm:"index:idx_created_at"`
}

type Comment struct {
	ID        int `gorm:"primary_key"`
	UserID    int `gorm:"index:idx_user_id"`
	VideoID   int `gorm:"index:idx_video_id"`
	Content   string
	CreatedAt time.Time `gorm:"index:idx_created_at"`
}

func InitDB() *gorm.DB {
	dsn := os.Getenv("dsn")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v\n", err)
	}
	_ = db.AutoMigrate(&Account{}, &Comment{}, &Follow{},
		&Favorite{}, &Video{})
	return db
}
