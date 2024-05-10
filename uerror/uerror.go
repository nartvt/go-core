package uerror

import (
	"errors"
	"fmt"
	"net/http"
)

type StatusError struct {
	Code    int64
	Key     string
	Message string
	Err     error
}

var (
	BAD_REQUEST           = "BAD_REQUEST_ERROR"
	SUCCESS               = "SUCCESS"
	ERROR                 = "ERROR"
	INTERNAL_SERVER_ERROR = "INTERNAL_SERVER_ERROR"
	NOT_FOUND             = "NOT_FOUND_ERROR"
	INPUT_PARAM_INVALID   = "INPUT_PARAM_INVALID"
	FORBIDEN_ERROR        = "FORBIDEN_ERROR"
	UNAUTHORIZED_ERROR    = "UNAUTHORIZED_ERROR"

	ErrorCode   = int64(1)
	SuccessCode = int64(0)
)

func newError(status int64, key string, err error) *StatusError {
	if err == nil {
		return &StatusError{
			Code:    status,
			Key:     key,
			Message: ERROR,
		}
	}
	return &StatusError{
		Code:    status,
		Message: err.Error(),
		Key:     key,
		Err:     err,
	}
}

func BadRequestError(message string) *StatusError {
	return newError(http.StatusBadRequest, BAD_REQUEST, errors.New(message))
}

func ParamInvalidError(message string) *StatusError {
	return newError(http.StatusBadRequest, INPUT_PARAM_INVALID, errors.New(message))
}

func InteralServerError(message string) *StatusError {
	return newError(http.StatusInternalServerError, INTERNAL_SERVER_ERROR, errors.New(message))
}

func NotFoundError(message string) *StatusError {
	return newError(http.StatusNotFound, NOT_FOUND, errors.New(message))
}

func ForbidenError(message string) *StatusError {
	return newError(http.StatusForbidden, FORBIDEN_ERROR, errors.New(message))
}

func UnAuthorizeError(message string) *StatusError {
	return newError(http.StatusUnauthorized, UNAUTHORIZED_ERROR, errors.New(message))
}

func (s *StatusError) Error() string {
	if s.Err != nil {
		return s.Err.Error()
	}
	return ""
}

func Recover() {
	if r := recover(); r != nil {
		fmt.Println("Program is panicking with value {}", r)
	}
}
