package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/huaouo/t4k/common"
	"github.com/huaouo/t4k/t4k-rdbms-service/rpc"
	"github.com/huaouo/t4k/t4k-video-service/handler"
	"log"
	"os"
)

func main() {
	binding.Validator = new(common.DefaultValidator)
	router := gin.Default()

	objectServiceEndpoint := os.Getenv("OBJECT_SERVICE_ADDR") + ":" + os.Getenv("OBJECT_SERVICE_LISTEN_PORT")
	videoHandler := handler.VideoHandler{
		VideoClient:           rpc.NewVideoClient(rpc.NewRdbmsClient()),
		ObjectServiceEndpoint: objectServiceEndpoint,
	}
	router.POST("/douyin/publish/action/", videoHandler.Publish)

	port := os.Getenv("VIDEO_SERVICE_LISTEN_PORT")
	err := router.Run(":" + port)
	if err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
