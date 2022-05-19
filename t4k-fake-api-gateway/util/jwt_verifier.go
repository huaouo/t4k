package util

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/huaouo/t4k/common"
	"io/ioutil"
	"log"
	"os"
)

type JwtVerifier struct {
	pubKey *ecdsa.PublicKey
}

func NewJwtVerifier() JwtVerifier {
	keyPath := os.Getenv("PUBLIC_KEY_PATH")
	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("failed to load public key: %v", err)
	}
	pubKey, err := jwt.ParseECPublicKeyFromPEM(keyBytes)
	if err != nil {
		log.Fatalf("failed to parse public key: %v", err)
	}
	return JwtVerifier{pubKey: pubKey}
}

func (s *JwtVerifier) VerifyAndEncode(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(*jwt.Token) (interface{}, error) {
		return s.pubKey, nil
	})
	if err != nil {
		log.Printf("failed to parse JWT token: %v", err)
		return "", common.ErrVerifyJwt
	}

	err = token.Claims.Valid()
	if err != nil {
		log.Printf("failed to verify JWT token: %v", err)
		return "", common.ErrVerifyJwt
	}

	jsonBytes, err := json.Marshal(token.Claims.(jwt.MapClaims))
	if err != nil {
		log.Printf("failed to encode JWT claims: %v", err)
		return "", common.ErrVerifyJwt
	}

	return base64.StdEncoding.EncodeToString(jsonBytes), nil
}
