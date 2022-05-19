package handler

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/huaouo/t4k/common"
	"io"
	"log"
	"net/http"
)

// FakeWriterAt Ref: https://dev.to/flowup/using-io-reader-io-writer-in-go-to-stream-data-3i7b
type FakeWriterAt struct {
	w io.Writer
}

func (fw FakeWriterAt) WriteAt(p []byte, _ int64) (n int, err error) {
	// ignore 'offset' because we forced sequential downloads
	return fw.w.Write(p)
}

type ObjectHandler struct {
	Uploader   *s3manager.Uploader
	Downloader *s3manager.Downloader
}

func (h *ObjectHandler) Get(c *gin.Context, bucket string) {
	filename := c.Param(common.ObjectServiceFilenameParam)
	headOutput, err := h.Downloader.S3.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		log.Printf("failed to get object attributes: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	r, w := io.Pipe()
	go func() {
		_, err := h.Downloader.Download(
			FakeWriterAt{w},
			&s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(filename),
			})
		if err != nil {
			log.Printf("failed to get object: %v", err)
		}
	}()
	c.DataFromReader(http.StatusOK, *headOutput.ContentLength, "application/octet-stream", r, nil)
}

func (h *ObjectHandler) Put(c *gin.Context, bucket string) {
	filename := c.Param(common.ObjectServiceFilenameParam)
	_, err := h.Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   c.Request.Body,
	})
	if err != nil {
		log.Printf("failed to put object: %v", err)
		c.Status(http.StatusInternalServerError)
	}
	c.Status(http.StatusOK)
}

func (h *ObjectHandler) GetCover(c *gin.Context) {
	h.Get(c, common.S3CoverBucketName)
}

func (h *ObjectHandler) PutCover(c *gin.Context) {
	h.Put(c, common.S3CoverBucketName)
}

func (h *ObjectHandler) GetVideo(c *gin.Context) {
	h.Get(c, common.S3VideoBucketName)
}

func (h *ObjectHandler) PutVideo(c *gin.Context) {
	h.Put(c, common.S3VideoBucketName)
}
