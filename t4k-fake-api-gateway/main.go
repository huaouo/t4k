package main

import (
	"github.com/huaouo/t4k/t4k-fake-api-gateway/handler"
	"github.com/huaouo/t4k/t4k-fake-api-gateway/util"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("FAKE_API_SERVER_LISTEN_PORT")
	err := http.ListenAndServe(":"+port, &handler.ApiGatewayHandler{
		Verifier:           util.NewJwtVerifier(),
		RouteTable:         util.NewRouteTable(),
		JwtVerifyWhitelist: util.NewJwtVerifyWhitelist(),
	})
	if err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
