package main

import (
	"bytes"
	"fmt"
	"github.com/huaouo/t4k/common"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"log"
	"net/http"
	"os"
)

const VideoPath = "/tmp/video"

type CoverGenerator struct {
	VideoUrlPrefix string
	CoverUrlPrefix string
	HttpClient     http.Client
}

func (g *CoverGenerator) GenerateCover(objectId string) error {
	err := g.fetchVideo(objectId)
	if err != nil {
		return err
	}
	r, err := g.takeFrame()
	if err != nil {
		return err
	}
	err = g.uploadCover(objectId, r)
	if err != nil {
		return err
	}
	return nil
}

func (g *CoverGenerator) fetchVideo(objectId string) error {
	url := g.VideoUrlPrefix + objectId
	videoReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("failed to create request to get the video: %v", err)
		return err
	}
	videoResp, err := g.HttpClient.Do(videoReq)
	if err != nil {
		log.Printf("failed to request the video: %v", err)
		return err
	}
	defer videoResp.Body.Close()
	if videoResp.StatusCode != http.StatusOK {
		log.Printf("failed to request the video: http code is %v", videoResp.StatusCode)
		return common.ErrInternal
	}
	videoFile, err := os.Create(VideoPath)
	if err != nil {
		log.Printf("failed to create video file: %v", err)
		return err
	}
	defer videoFile.Close()
	_, err = io.Copy(videoFile, videoResp.Body)
	if err != nil {
		log.Printf("failed to write video file: %v", err)
		return err
	}
	if err != nil {
		log.Printf("failed to write video file: %v", err)
		return err
	}
	return nil
}

func (g *CoverGenerator) takeFrame() (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(VideoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 5)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "png"}).
		WithOutput(buf).
		Run()
	if err != nil {
		log.Printf("failed to generate cover: %v", err)
		return nil, err
	}
	return buf, nil
}

func (g *CoverGenerator) uploadCover(objectId string, r io.Reader) error {
	url := g.CoverUrlPrefix + objectId
	coverReq, err := http.NewRequest(http.MethodPut, url, r)
	if err != nil {
		log.Printf("failed to create request to upload the video: %v", err)
		return err
	}
	coverResp, err := g.HttpClient.Do(coverReq)
	if err != nil {
		log.Printf("failed to upload the cover: %v", err)
		return err
	}
	if coverResp.StatusCode != http.StatusOK {
		log.Printf("failed to upload the cover: http code is %v", coverResp.StatusCode)
		return common.ErrInternal
	}
	return nil
}
