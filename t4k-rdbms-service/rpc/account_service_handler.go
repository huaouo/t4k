package rpc

import (
	"context"
	"encoding/base64"
	"github.com/huaouo/t4k/common"
	"github.com/huaouo/t4k/t4k-rdbms-service/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"sync"
)

var (
	EmptyAuthNResp = &AuthNResponse{}
	EmptyInfoResp  = &InfoResponse{}
)

type AccountHandler struct {
	UnimplementedAccountServer
	DB *gorm.DB
}

func (h *AccountHandler) CreateAccount(ctx context.Context, req *AuthNRequest) (*AuthNResponse, error) {
	var exists bool
	err := h.DB.Model(&repository.Account{}).
		Select("count(1) > 0").
		Where("name = ?", req.GetName()).
		First(&exists).
		Error
	if err != nil {
		log.Printf("rdbms access error: %v", err)
		return EmptyAuthNResp, common.ErrInternal
	} else if exists {
		return EmptyAuthNResp, common.ErrUserAlreadyExist
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to generate password hash: %v", err)
		return EmptyAuthNResp, common.ErrInternal
	}
	b64 := base64.StdEncoding.EncodeToString(hash)
	err = h.DB.Create(&repository.Account{Name: req.GetName(), Password: b64}).Error
	if err != nil {
		log.Printf("rdbms access error: %v", err)
		return EmptyAuthNResp, common.ErrInternal
	}

	var userId uint64
	err = h.DB.Model(&repository.Account{}).
		Select("id").
		Where("name = ?", req.GetName()).
		First(&userId).Error
	if err != nil {
		log.Printf("rdbms access error: %v", err)
		return EmptyAuthNResp, common.ErrInternal
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
		return EmptyAuthNResp, common.ErrInternal
	}
	if account.Name == "" {
		return EmptyAuthNResp, common.ErrUserNotExist
	}
	hashedPassword, err := base64.StdEncoding.DecodeString(account.Password)
	if err != nil {
		log.Printf("invalid password hash in rdbms: %v", err)
		return EmptyAuthNResp, common.ErrInternal
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(req.GetPassword()))
	if err != nil {
		return EmptyAuthNResp, common.ErrPasswordIncorrect
	}

	return &AuthNResponse{UserId: account.Id}, nil
}

func (h *AccountHandler) GetUserInfo(ctx context.Context, in *InfoRequest) (*InfoResponse, error) {
	resp := &InfoResponse{}

	err := h.DB.Model(&repository.Account{}).
		Select("name").
		Where("id = ?", in.GetUserId()).
		FirstOrInit(&resp.Name).
		Error
	if err != nil {
		log.Printf("rdbms access error: %v", err)
		return EmptyInfoResp, common.ErrInternal
	}
	if resp.Name == "" {
		return EmptyInfoResp, common.ErrUserNotExist
	}

	var wg sync.WaitGroup
	wg.Add(3)
	errChan := make(chan error, 3)
	followQuery := func(fieldName string, result *uint64) {
		errChan <- h.DB.Model(&repository.Follow{}).
			Select("count(1)").
			Where(fieldName+" = ?", in.GetUserId()).
			First(result).
			Error
		wg.Done()
	}
	go followQuery("user_id", &resp.FollowCount)
	go followQuery("to_user_id", &resp.FollowerCount)
	go func() {
		var isFollow int
		errChan <- h.DB.Model(&repository.Follow{}).
			Select("count(1)").
			Where("user_id = ? and to_user_id = ?", in.GetSignInUserId(), in.GetUserId()).
			First(&isFollow).
			Error
		if isFollow != 0 {
			resp.IsFollow = true
		}
		wg.Done()
	}()
	wg.Wait()
	close(errChan)
	var containError bool
	for e := range errChan {
		if e != nil {
			log.Printf("rdbms access error: %v", e)
			containError = true
		}
	}
	if containError {
		return EmptyInfoResp, common.ErrInternal
	}
	return resp, nil
}
