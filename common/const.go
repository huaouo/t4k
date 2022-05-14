package common

import "errors"

var (
	ErrInternal          = errors.New("internal error")
	ErrUserAlreadyExist  = errors.New("user already exist")
	ErrUserNotExist      = errors.New("user not exist")
	ErrPasswordIncorrect = errors.New("password incorrect")
	ErrSignJwt           = errors.New("cannot sign jwt")
	ErrVerifiyJwt        = errors.New("failed to verify jwt")
)

var (
	StatusSuccess = 0
	StatusFailure = 1
)

var (
	ExtractedJwtPayloadName = "Jwt-Payload"
	JwtPayloadUserIdName    = "uid"
)
