package repository

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

type Account struct {
	Id        uint64 `gorm:"primary_key"`
	Name      string `gorm:"index:idx_name,unique"`
	Password  string
	CreatedAt time.Time
}

type Follow struct {
	UserId    uint64    `gorm:"index:idx_user_id"`
	ToUserId  uint64    `gorm:"index:idx_to_user_id"`
	CreatedAt time.Time `gorm:"index:idx_created_at"`
}

type Video struct {
	Id        uint64 `gorm:"primary_key"`
	UserId    uint64
	ObjectId  string
	Title     string
	CreatedAt time.Time `gorm:"index:idx_created_at"`
}

type Favorite struct {
	UserId    uint64    `gorm:"index:idx_user_id"`
	VideoId   uint64    `gorm:"index:idx_video_id"`
	CreatedAt time.Time `gorm:"index:idx_created_at"`
}

type Comment struct {
	Id        uint64 `gorm:"primary_key"`
	UserId    uint64 `gorm:"index:idx_user_id"`
	VideoId   uint64 `gorm:"index:idx_video_id"`
	Content   string
	CreatedAt time.Time `gorm:"index:idx_created_at"`
}

func InitDB() *gorm.DB {
	dsn := os.Getenv("RDBMS_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v\n", err)
	}
	_ = db.AutoMigrate(&Account{}, &Comment{}, &Follow{},
		&Favorite{}, &Video{})
	return db
}
