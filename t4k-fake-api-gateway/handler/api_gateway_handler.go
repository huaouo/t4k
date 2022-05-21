package handler

import (
	"bytes"
	"fmt"
	"github.com/huaouo/t4k/common"
	"github.com/huaouo/t4k/t4k-fake-api-gateway/util"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

type ApiGatewayHandler struct {
	Verifier           util.JwtVerifier
	RouteTable         map[string]string
	JwtVerifyWhitelist map[string]bool
	HttpClient         http.Client
}

func createBodyWithMultiPart(form *multipart.Form) (io.Reader, string, error) {
	body := new(bytes.Buffer)
	mp := multipart.NewWriter(body)
	defer mp.Close()

	for key, vals := range form.File {
		for _, val := range vals {
			log.Printf("key = %v, value = %v", key, val.Filename)
			w, err := mp.CreateFormFile(key, val.Filename)
			if err != nil {
				log.Printf("failed to create file: %v", err)
				return nil, "", err
			}
			f, err := val.Open()
			if err != nil {
				log.Printf("failed to create file: %v", err)
				return nil, "", err
			}
			defer f.Close()
			_, err = io.Copy(w, f)
			if err != nil {
				log.Printf("failed to write file: %v", err)
				return nil, "", err
			}
		}
	}
	for key, vals := range form.Value {
		for _, val := range vals {
			err := mp.WriteField(key, val)
			if err != nil {
				log.Printf("failed to write field: %v", err)
				return nil, "", err
			}
		}
	}
	return body, mp.Boundary(), nil
}

func (h *ApiGatewayHandler) ServeHTTP(respWriter http.ResponseWriter, req *http.Request) {
	var tokenInMultipartBody bool
	if _, ok := h.JwtVerifyWhitelist[req.URL.Path]; !ok {
		tokenStr := req.URL.Query().Get("token")
		if tokenStr == "" {
			err := req.ParseMultipartForm(10 << 20)
			if err != nil {
				log.Printf("unable to get token: %v", err)
				http.Error(respWriter, "403 Forbidden", http.StatusForbidden)
				return
			}
			tokenStr = req.MultipartForm.Value["token"][0]
			tokenInMultipartBody = true
		}
		var err error
		encodedJwtPayload, err := h.Verifier.VerifyAndEncode(tokenStr)
		if err != nil {
			http.Error(respWriter, "403 Forbidden", http.StatusForbidden)
			return
		}
		req.Header.Set(common.ExtractedJwtPayloadName, encodedJwtPayload)
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
	body := req.Body
	var newReq *http.Request

	var err error
	if tokenInMultipartBody {
		var body io.Reader
		var boundary string
		body, boundary, err = createBodyWithMultiPart(req.MultipartForm)
		if err != nil {
			http.Error(respWriter, "502 Bad Gateway", http.StatusBadGateway)
			return
		}
		newReq, err = http.NewRequest(req.Method, url.String(), body)
		req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", boundary))
	} else {
		newReq, err = http.NewRequest(req.Method, url.String(), body)
	}
	if err != nil {
		log.Printf("failed to create new request: %v", err)
		http.Error(respWriter, "502 Bad Gateway", http.StatusBadGateway)
		return
	}
	newReq.Header = req.Header

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
