package rpc

import (
	"context"
	"encoding/base64"
	"github.com/huaouo/t4k/common"
	"github.com/huaouo/t4k/t4k-rdbms-service/repository"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"log"
)

var Empty = &emptypb.Empty{}

type AccountHandler struct {
	UnimplementedAccountServer
	DB *gorm.DB
}

func (h *AccountHandler) CreateAccount(ctx context.Context, req *CreateAccountRequest) (*emptypb.Empty, error) {
	var exists bool
	err := h.DB.Model(&repository.Account{}).
		Select("count(*) > 0").
		Where("name = ?", req.GetName()).
		Find(&exists).
		Error
	if err != nil {
		log.Printf("%v", err)
		return Empty, common.ErrInternal
	} else if exists {
		return Empty, common.ErrUserAlreadyExist
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("%v", err)
		return Empty, common.ErrInternal
	}
	b64 := base64.StdEncoding.EncodeToString(hash)
	err = h.DB.Create(&repository.Account{Name: req.GetName(), Password: b64}).Error
	if err != nil {
		log.Printf("%v", err)
		return Empty, common.ErrInternal
	}
	return Empty, nil
}

func (h *AccountHandler) Authenticate(ctx context.Context, req *AuthenticateRequest) (*emptypb.Empty, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("%v", err)
		return Empty, common.ErrInternal
	}

	b64 := base64.StdEncoding.EncodeToString(hash)
	var b64InDB string
	err = h.DB.Model(&repository.Account{}).
		Select("password").
		Where("name = ?", req.GetName()).
		Find(&b64InDB).
		Error
	if err != nil {
		log.Printf("%v", err)
		return Empty, common.ErrInternal
	}
	if b64InDB != b64 {
		return Empty, common.ErrPasswordIncorrect
	}
	return Empty, nil
}
