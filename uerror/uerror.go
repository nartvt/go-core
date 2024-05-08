package uerror

import (
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

func BadRequestError(err error) *StatusError {
	return newError(http.StatusBadRequest, BAD_REQUEST, err)
}

func ParamInvalidError(err error) *StatusError {
	return newError(http.StatusBadRequest, INPUT_PARAM_INVALID, err)
}

func InteralServerError(err error) *StatusError {
	return newError(http.StatusInternalServerError, INTERNAL_SERVER_ERROR, err)
}

func NotFoundError(err error) *StatusError {
	return newError(http.StatusNotFound, NOT_FOUND, err)
}

func ForbidenError(err error) *StatusError {
	return newError(http.StatusForbidden, FORBIDEN_ERROR, err)
}

func UnAuthorizeError(err error) *StatusError {
	return newError(http.StatusUnauthorized, UNAUTHORIZED_ERROR, err)
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
