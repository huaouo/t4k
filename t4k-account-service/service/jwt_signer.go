package service

import (
	"crypto/ecdsa"
	"github.com/golang-jwt/jwt"
	"github.com/huaouo/t4k/common"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type JwtSigner struct {
	prvKey *ecdsa.PrivateKey
}

func NewJwtSigner() JwtSigner {
	keyPath := os.Getenv("PRIVATE_KEY_PATH")
	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("failed to load private key: %v", err)
	}
	prvKey, err := jwt.ParseECPrivateKeyFromPEM(keyBytes)
	if err != nil {
		log.Fatalf("failed to parse private key: %v", err)
	}
	return JwtSigner{prvKey: prvKey}
}

func (s *JwtSigner) Sign(userId uint64) (string, error) {
	now := time.Now().UTC()
	claims := make(jwt.MapClaims)
	claims["uid"] = userId
	claims["exp"] = now.Add(time.Hour * 24 * 7).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString(s.prvKey)
	if err != nil {
		log.Printf("cannot sign jwt: %v", err)
		return "", common.ErrSignJwt
	}

	return token, nil
}
