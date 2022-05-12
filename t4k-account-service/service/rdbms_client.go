package service

import (
	"github.com/huaouo/t4k/t4k-rdbms-service/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

func NewRdbmsAccountClient() rpc.AccountClient {
	option := grpc.WithTransportCredentials(insecure.NewCredentials())
	addr := os.Getenv("RDBMS_SERVICE_ADDR")
	port := os.Getenv("RDBMS_SERVICE_LISTEN_PORT")
	conn, err := grpc.Dial(addr+":"+port, option)
	if err != nil {
		log.Fatalf("failed to connect RDBMS service: %v", err)
	}
	return rpc.NewAccountClient(conn)
}
