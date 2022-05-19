package rpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

func NewRdbmsClient() *grpc.ClientConn {
	option := grpc.WithTransportCredentials(insecure.NewCredentials())
	addr := os.Getenv("RDBMS_SERVICE_ADDR")
	port := os.Getenv("RDBMS_SERVICE_LISTEN_PORT")
	conn, err := grpc.Dial(addr+":"+port, option)
	if err != nil {
		log.Fatalf("failed to connect RDBMS service: %v", err)
	}
	return conn
}
