package rpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

func NewMqServiceClient() *grpc.ClientConn {
	option := grpc.WithTransportCredentials(insecure.NewCredentials())
	addr := os.Getenv("MQ_SERVICE_ADDR")
	port := os.Getenv("MQ_SERVICE_LISTEN_PORT")
	conn, err := grpc.Dial(addr+":"+port, option)
	if err != nil {
		log.Fatalf("failed to connect MQ service: %v", err)
	}
	return conn
}
