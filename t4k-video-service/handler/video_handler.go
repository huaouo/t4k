package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/huaouo/t4k/common"
	mq "github.com/huaouo/t4k/t4k-mq-service/rpc"
	"github.com/huaouo/t4k/t4k-rdbms-service/rpc"
	rdbms "github.com/huaouo/t4k/t4k-rdbms-service/rpc"
	"github.com/huaouo/t4k/t4k-video-service/model/request"
	"github.com/huaouo/t4k/t4k-video-service/model/response"
	"log"
	"net/http"
)

type VideoHandler struct {
	VideoClient           rdbms.VideoClient
	MqClient              mq.MqClient
	HttpClient            http.Client
	ObjectServiceEndpoint string
}

func (h *VideoHandler) Publish(c *gin.Context) {
	var req request.Publish
	resp := response.Publish{}
	err := c.ShouldBind(&req)
	if err != nil {
		h.handleError(c, err, "cannot bind request", "invalid request")
		return
	}

	signInUserId, err := common.ExtractSignInUserId(c)
	if err != nil {
		h.handleError(c, err, "", "")
		return
	}

	rpcResp, err := h.VideoClient.Create(context.TODO(), &rpc.CreateVideoRequest{
		UserId: signInUserId,
		Title:  req.Title,
	})
	if err != nil {
		h.handleError(c, err, "failed to create item in rdbms", "")
		return
	}

	videoFile, err := req.Data.Open()
	if err != nil {
		h.handleError(c, err, "failed to open video file", "")
		return
	}

	url := "http://" + h.ObjectServiceEndpoint +
		common.ObjectServiceVideoPathPrefix + rpcResp.GetObjectId()
	videoReq, err := http.NewRequest(http.MethodPut, url, videoFile)
	if err != nil {
		h.handleError(c, err, "failed to create new video request", "")
		return
	}
	videoResp, err := h.HttpClient.Do(videoReq)
	if err != nil {
		h.handleError(c, err, "failed to upload video", "")
		return
	}
	if videoResp.StatusCode != http.StatusOK {
		log.Printf("failed to upload video: response status is %v", videoResp.Status)
		resp.StatusCode = common.StatusFailure
		resp.StatusMsg = common.ErrInternal.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	pubClient, err := h.MqClient.Publish(context.TODO())
	if err != nil {
		log.Printf("failed to initialize MQ service publish client: %v", err)
	} else {
		err = pubClient.Send(&mq.PubRequest{
			QueueName: common.MqCoverQueueName,
			Content:   []byte(rpcResp.GetObjectId()),
		})
		if err != nil {
			log.Printf("failed to publish cover generating task: %v", err)
		}
	}

	resp.StatusCode = common.StatusSuccess
	c.JSON(http.StatusOK, resp)
}

func (h *VideoHandler) handleError(c *gin.Context, err error, logMsg, returnMsg string) {
	resp := response.Publish{}
	if logMsg != "" {
		log.Printf("%s: %v", logMsg, err)
	}
	resp.StatusCode = common.StatusFailure
	if returnMsg == "" {
		resp.StatusMsg = common.ErrInternal.Error()
	} else {
		resp.StatusMsg = returnMsg
	}
	c.JSON(http.StatusOK, resp)
	return
}
