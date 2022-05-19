package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/huaouo/t4k/common"
	"github.com/huaouo/t4k/t4k-account-service/model/request"
	"github.com/huaouo/t4k/t4k-account-service/model/response"
	"github.com/huaouo/t4k/t4k-account-service/util"
	"github.com/huaouo/t4k/t4k-rdbms-service/rpc"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

type AccountHandler struct {
	AccountClient rpc.AccountClient
	Signer        util.JwtSigner
}

func (h *AccountHandler) Sign(c *gin.Context,
	f func(context.Context, *rpc.AuthNRequest, ...grpc.CallOption) (*rpc.AuthNResponse, error)) {
	var req request.Sign
	resp := response.Sign{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Printf("cannot bind request: %v", err)
		resp.StatusCode = common.StatusFailure
		resp.StatusMsg = "invalid username or password"
		c.JSON(http.StatusOK, resp)
		return
	}

	authNResp, err := f(context.TODO(), &rpc.AuthNRequest{
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		log.Printf("%v", err)
		resp.StatusCode = common.StatusFailure
		resp.StatusMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Token, err = h.Signer.Sign(authNResp.UserId)
	if err != nil {
		log.Printf("%v", err)
		resp.StatusCode = common.StatusFailure
		resp.StatusMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.StatusCode = common.StatusSuccess
	resp.UserId = authNResp.UserId
	c.JSON(http.StatusOK, resp)
}

func (h *AccountHandler) SignUp(c *gin.Context) {
	h.Sign(c, h.AccountClient.Create)
}

func (h *AccountHandler) SignIn(c *gin.Context) {
	h.Sign(c, h.AccountClient.Authenticate)
}

func (h *AccountHandler) Info(c *gin.Context) {
	var req request.Info
	resp := response.Info{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.Printf("cannot bind request: %v", err)
		resp.StatusCode = common.StatusFailure
		resp.StatusMsg = "invalid request"
		c.JSON(http.StatusOK, resp)
		return
	}

	signInUserId, err := common.ExtractSignInUserId(c)
	if err != nil {
		resp.StatusCode = common.StatusFailure
		resp.StatusMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	infoResp, err := h.AccountClient.GetUserInfo(context.TODO(), &rpc.InfoRequest{
		SignInUserId: signInUserId,
		UserId:       req.UserId,
	})
	if err != nil {
		log.Printf("%v", err)
		resp.StatusCode = common.StatusFailure
		resp.StatusMsg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.User = &response.User{
		Id:            req.UserId,
		Name:          infoResp.GetName(),
		FollowCount:   infoResp.GetFollowCount(),
		FollowerCount: infoResp.GetFollowerCount(),
		IsFollow:      infoResp.GetIsFollow(),
	}
	c.JSON(http.StatusOK, resp)
	return
}
