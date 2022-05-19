package main

import (
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/huaouo/t4k/common"
	"github.com/huaouo/t4k/t4k-object-service/handler"
	"github.com/huaouo/t4k/t4k-object-service/util"
	"log"
	"os"
)

func main() {
	binding.Validator = new(common.DefaultValidator)
	router := gin.Default()

	s3Session := util.NewS3Session()
	uploader := s3manager.NewUploader(s3Session)
	downloader := s3manager.NewDownloader(s3Session)
	uploader.Concurrency = 1
	downloader.Concurrency = 1
	objectHandler := handler.ObjectHandler{
		Uploader:   uploader,
		Downloader: downloader,
	}

	coverPath := common.ObjectServiceCoverPathPrefix + ":" + common.ObjectServiceFilenameParam
	videoPath := common.ObjectServiceVideoPathPrefix + ":" + common.ObjectServiceFilenameParam
	router.GET(coverPath, objectHandler.GetCover)
	router.GET(videoPath, objectHandler.GetVideo)
	router.PUT(coverPath, objectHandler.PutCover)
	router.PUT(videoPath, objectHandler.PutVideo)

	port := os.Getenv("OBJECT_SERVICE_LISTEN_PORT")
	err := router.Run(":" + port)
	if err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
