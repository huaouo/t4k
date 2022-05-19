package rpc

import (
	"context"
	"github.com/huaouo/t4k/common"
	"github.com/huaouo/t4k/t4k-rdbms-service/repository"
	"github.com/u2takey/go-utils/uuid"
	"gorm.io/gorm"
	"log"
)

type VideoHandler struct {
	UnimplementedVideoServer
	DB *gorm.DB
}

func (h *VideoHandler) Create(ctx context.Context, req *CreateVideoRequest) (*CreateVideoResponse, error) {
	objectId := uuid.NewUUID()
	err := h.DB.Create(&repository.Video{
		UserId:   req.GetUserId(),
		ObjectId: objectId,
		Title:    req.GetTitle(),
	}).Error
	if err != nil {
		log.Printf("failed to create video item: %v", err)
		return &CreateVideoResponse{}, common.ErrInternal
	}

	return &CreateVideoResponse{
		ObjectId: objectId,
	}, nil
}
