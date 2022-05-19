package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/huaouo/t4k/common"
	"github.com/huaouo/t4k/t4k-account-service/handler"
	"github.com/huaouo/t4k/t4k-account-service/util"
	"github.com/huaouo/t4k/t4k-rdbms-service/rpc"
	"log"
	"os"
)

func main() {
	binding.Validator = new(common.DefaultValidator)
	router := gin.Default()

	accountHandler := handler.AccountHandler{
		AccountClient: rpc.NewAccountClient(rpc.NewRdbmsClient()),
		Signer:        util.NewJwtSigner(),
	}
	router.POST("/douyin/user/register/", accountHandler.SignUp)
	router.POST("/douyin/user/login/", accountHandler.SignIn)
	router.GET("/douyin/user/", accountHandler.Info)

	port := os.Getenv("ACCOUNT_SERVICE_LISTEN_PORT")
	err := router.Run(":" + port)
	if err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
