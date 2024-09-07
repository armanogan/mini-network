package errors

type CodedError interface {
	error
	GetCode() int
}

type ErrorWithCode struct {
	Code    int
	message string
}

func NewErrorWithCode(code int, message string) *ErrorWithCode {
	return &ErrorWithCode{Code: code, message: message}
}

// Error возвращает сообщение ошибки.
func (e *ErrorWithCode) Error() string {
	return e.message
}
func (e *ErrorWithCode) GetCode() int {
	return e.Code
}
