package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/huaouo/t4k/common"
	"github.com/huaouo/t4k/t4k-account-service/model/request"
	"github.com/huaouo/t4k/t4k-account-service/model/response"
	"github.com/huaouo/t4k/t4k-account-service/service"
	"github.com/huaouo/t4k/t4k-rdbms-service/rpc"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

type AccountHandler struct {
	AccountClient rpc.AccountClient
	Signer        service.JwtSigner
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
	h.Sign(c, h.AccountClient.CreateAccount)
}

func (h *AccountHandler) SignIn(c *gin.Context) {
	h.Sign(c, h.AccountClient.Authenticate)
}
