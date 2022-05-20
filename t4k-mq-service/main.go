package main

import (
	"github.com/huaouo/t4k/t4k-mq-service/rpc"
	"github.com/huaouo/t4k/t4k-mq-service/util"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	port := os.Getenv("MQ_SERVICE_LISTEN_PORT")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	rpc.RegisterMqServer(grpcServer, &rpc.MqHandler{
		Conn: util.NewAmqpConnection(),
	})
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
