package errorx

import (
	"errors"
	"fmt"
)

type CodeError struct {
	Code int64
	Msg  string
}

func NewCodeError(code int64, msg string) error {
	return &CodeError{
		Code: code,
		Msg:  msg,
	}
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("%d:%s", e.Code, e.Msg)
}

func FromError(err error) (*CodeError, bool) {
	if err == nil {
		return nil, false
	}

	var e *CodeError
	ok := errors.As(err, &e)
	return e, ok
}
