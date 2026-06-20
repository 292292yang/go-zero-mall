package errorx

type CodeError struct {
	Code int64  `json:"code"`
	Msg  string `json:"message"`
}

func NewCodeError(code int64, msg string) error {
	return &CodeError{
		Code: code,
		Msg:  msg,
	}
}

func (e *CodeError) Error() string {
	return e.Msg
}

func (e *CodeError) Data() *CodeError {
	return e
}
