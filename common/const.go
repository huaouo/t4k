package common

import "errors"

var (
	ErrInternal          = errors.New("internal error")
	ErrUserAlreadyExist  = errors.New("user already exist")
	ErrUserNotExist      = errors.New("user not exist")
	ErrPasswordIncorrect = errors.New("password incorrect")
)

var (
	StatusSuccess = 0
	StatusFailure = 1
)
