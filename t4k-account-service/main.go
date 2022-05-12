package main

import (
	"github.com/gin-gonic/gin"
	"github.com/huaouo/t4k/t4k-account-service/handler"
	"github.com/huaouo/t4k/t4k-account-service/service"
	"log"
)

func main() {
	//binding.Validator = new(common.DefaultValidator)
	router := gin.Default()

	accountHandler := handler.AccountHandler{
		AccountClient: service.NewRDBMSAccountClient(),
	}
	router.POST("/douyin/user/register/", accountHandler.SignUp)
	router.POST("/douyin/user/login/", accountHandler.SignIn)

	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
