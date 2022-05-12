package rpc

import (
	"context"
	"encoding/base64"
	"github.com/huaouo/t4k/common"
	"github.com/huaouo/t4k/t4k-rdbms-service/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
)

var Empty = &AuthNResponse{}

type AccountHandler struct {
	UnimplementedAccountServer
	DB *gorm.DB
}

func (h *AccountHandler) CreateAccount(ctx context.Context, req *AuthNRequest) (*AuthNResponse, error) {
	var exists bool
	err := h.DB.Model(&repository.Account{}).
		Select("count(*) > 0").
		Where("name = ?", req.GetName()).
		First(&exists).
		Error
	if err != nil {
		log.Printf("rdbms access error: %v", err)
		return Empty, common.ErrInternal
	} else if exists {
		return Empty, common.ErrUserAlreadyExist
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to generate password hash: %v", err)
		return Empty, common.ErrInternal
	}
	b64 := base64.StdEncoding.EncodeToString(hash)
	err = h.DB.Create(&repository.Account{Name: req.GetName(), Password: b64}).Error
	if err != nil {
		log.Printf("rdbms access error: %v", err)
		return Empty, common.ErrInternal
	}

	var userId uint64
	err = h.DB.Model(&repository.Account{}).
		Select("id").
		Where("name = ?", req.GetName()).
		First(&userId).Error
	if err != nil {
		log.Printf("rdbms access error: %v", err)
		return Empty, common.ErrInternal
	}

	return &AuthNResponse{UserId: userId}, nil
}

func (h *AccountHandler) Authenticate(ctx context.Context, req *AuthNRequest) (*AuthNResponse, error) {
	var account repository.Account
	err := h.DB.Model(&repository.Account{}).
		Select("*").
		Where("name = ?", req.GetName()).
		FirstOrInit(&account).
		Error
	if err != nil {
		log.Printf("rdbms access error: %v", err)
		return Empty, common.ErrInternal
	}
	if account.Name == "" {
		return Empty, common.ErrUserNotExist
	}
	hashedPassword, err := base64.StdEncoding.DecodeString(account.Password)
	if err != nil {
		log.Printf("invalid password hash in rdbms: %v", err)
		return Empty, common.ErrInternal
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(req.GetPassword()))
	if err != nil {
		return Empty, common.ErrPasswordIncorrect
	}

	return &AuthNResponse{UserId: account.Id}, nil
}
