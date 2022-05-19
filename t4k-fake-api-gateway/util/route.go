package util

import "os"

func NewRouteTable() map[string]string {
	accountService := os.Getenv("ACCOUNT_SERVICE_ADDR") + ":" + os.Getenv("ACCOUNT_SERVICE_LISTEN_PORT")
	videoService := os.Getenv("VIDEO_SERVICE_ADDR") + ":" + os.Getenv("VIDEO_SERVICE_LISTEN_PORT")

	return map[string]string{
		"/douyin/user/register/":  accountService,
		"/douyin/user/login/":     accountService,
		"/douyin/user/":           accountService,
		"/douyin/publish/action/": videoService,
	}
}

func NewJwtVerifyWhitelist() map[string]bool {
	return map[string]bool{
		"/douyin/user/register/": true,
		"/douyin/user/login/":    true,
	}
}
