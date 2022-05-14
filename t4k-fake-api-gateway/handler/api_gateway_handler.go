package handler

import (
	"github.com/huaouo/t4k/common"
	"github.com/huaouo/t4k/t4k-fake-api-gateway/service"
	"io"
	"log"
	"net/http"
)

type ApiGatewayHandler struct {
	Verifier           service.JwtVerifier
	RouteTable         map[string]string
	JwtVerifyWhitelist map[string]bool
	HttpClient         http.Client
}

func (h *ApiGatewayHandler) ServeHTTP(respWriter http.ResponseWriter, req *http.Request) {
	var encodedJwtPayload string
	if _, ok := h.JwtVerifyWhitelist[req.URL.Path]; !ok {
		tokenStr := req.URL.Query().Get("token")
		var err error
		encodedJwtPayload, err = h.Verifier.VerifyAndEncode(tokenStr)
		if err != nil {
			http.Error(respWriter, "403 Forbidden", http.StatusForbidden)
			return
		}
	}

	newUrlHost, ok := h.RouteTable[req.URL.Path]
	if !ok {
		log.Printf("no route rules found: %s", req.URL.Path)
		http.Error(respWriter, "404 Not Found", http.StatusNotFound)
		return
	}

	url := req.URL
	url.Host = newUrlHost
	url.Scheme = "http"
	newReq, err := http.NewRequest(req.Method, url.String(), req.Body)
	newReq.Header = req.Header
	newReq.Header.Set(common.ExtractedJwtPayloadName, encodedJwtPayload)
	if err != nil {
		log.Printf("failed to create new request: %v", err)
		http.Error(respWriter, "502 Bad Gateway", http.StatusBadGateway)
		return
	}

	resp, err := h.HttpClient.Do(newReq)
	if err != nil {
		log.Printf("failed to forward request: %v", err)
		http.Error(respWriter, "502 Bad Gateway", http.StatusBadGateway)
		return
	}
	for k, vList := range resp.Header {
		for _, v := range vList {
			respWriter.Header().Add(k, v)
		}
	}
	respWriter.WriteHeader(resp.StatusCode)
	_, err = io.Copy(respWriter, resp.Body)
	if err != nil {
		log.Printf("failed to write upstream response: %v", err)
	}
	err = resp.Body.Close()
	if err != nil {
		log.Printf("failed to close the body of upstream response: %v", err)
	}
}
