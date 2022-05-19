package common

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
)

func ExtractSignInUserId(c *gin.Context) (uint64, error) {
	b64JwtPayload := c.GetHeader(ExtractedJwtPayloadName)
	jwtPayload, err := base64.StdEncoding.DecodeString(b64JwtPayload)
	if err != nil {
		log.Printf("cannot decode extracted jwt header: %v", err)
		return 0, ErrInternal
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(jwtPayload, &m)
	if err != nil {
		log.Printf("cannot unmarshal jwt json: %v", err)
		return 0, err
	}
	signInUserId := uint64(m[JwtPayloadUserIdName].(float64))
	return signInUserId, nil
}
