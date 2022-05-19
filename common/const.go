package common

import "errors"

var (
	ErrInternal          = errors.New("internal error")
	ErrUserAlreadyExist  = errors.New("user already exist")
	ErrUserNotExist      = errors.New("user not exist")
	ErrPasswordIncorrect = errors.New("password incorrect")
	ErrSignJwt           = errors.New("cannot sign jwt")
	ErrVerifyJwt         = errors.New("failed to verify jwt")
)

const (
	StatusSuccess = 0
	StatusFailure = 1
)

const (
	ExtractedJwtPayloadName = "Jwt-Payload"
	JwtPayloadUserIdName    = "uid"
)

const (
	S3CoverBucketName            = "t4k-cover"
	S3VideoBucketName            = "t4k-video"
	ObjectServiceCoverPathPrefix = "/s3/cover/"
	ObjectServiceVideoPathPrefix = "/s3/video/"
	ObjectServiceFilenameParam   = "filename"
)
