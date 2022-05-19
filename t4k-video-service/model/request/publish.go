package request

import "mime/multipart"

type Publish struct {
	Data  *multipart.FileHeader `form:"data" binding:"required"`
	Token string                `form:"token" binding:"required"`
	Title string                `form:"title" binding:"required"`
}
