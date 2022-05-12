package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/huaouo/t4k/common"
	"github.com/huaouo/t4k/t4k-account-service/model/request"
	"github.com/huaouo/t4k/t4k-account-service/model/response"
	"github.com/huaouo/t4k/t4k-rdbms-service/rpc"
	"log"
	"net/http"
)

type AccountHandler struct {
	AccountClient rpc.AccountClient
}

func (h *AccountHandler) SignUp(c *gin.Context) {
	var req request.SignUp
	resp := response.SignUp{}
	log.Printf("%s", c.Query("username"))
	log.Printf("%s", c.Query("password"))
	err := c.ShouldBind(&req)
	if err != nil {
		log.Printf("%v", err)
		resp.StatusCode = common.StatusFailure
		resp.StatusMsg = "invalid username or password"
		c.JSON(http.StatusOK, resp)
		return
	}

	_, err = h.AccountClient.CreateAccount(context.TODO(), &rpc.CreateAccountRequest{
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		log.Printf("%v", err)
		resp.StatusCode = common.StatusFailure
		resp.StatusMsg = err.Error()
		c.JSON(http.StatusOK, resp)
	}

	resp.StatusCode = common.StatusSuccess
	c.JSON(http.StatusOK, resp)
}

func (h *AccountHandler) SignIn(c *gin.Context) {

}
