package main

import (
	"github.com/huaouo/t4k/t4k-rdbms-service/repository"
	"github.com/huaouo/t4k/t4k-rdbms-service/rpc"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	port := os.Getenv("RDBMS_SERVICE_LISTEN_PORT")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	db := repository.InitDB()
	rpc.RegisterAccountServer(grpcServer, &rpc.AccountHandler{DB: db})
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
