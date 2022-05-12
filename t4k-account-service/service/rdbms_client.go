package service

import (
	"github.com/huaouo/t4k/t4k-rdbms-service/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func NewRDBMSAccountClient() rpc.AccountClient {
	option := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial("127.0.0.1:8888", option)
	if err != nil {
		log.Fatalf("failed to connect RDBMS service: %v", err)
	}
	return rpc.NewAccountClient(conn)
}
