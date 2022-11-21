package errors

import (
	"fmt"

	ut "github.com/go-playground/universal-translator"
)

var (
	// ErrParamsMissing ...
	ErrParamsMissing = NewWithStatus("缺少参数", 400)
	// ErrNotFound ...
	ErrNotFound = NewWithStatus("未找到", 404)
	// ErrSignature ...
	ErrSignature = NewWithStatus("签名错误", 400)
	// ErrSignTimestamp ...
	ErrSignTimestamp = NewWithStatus("时间戳错误", 400)
	// ErrTokenInvalid ...
	ErrTokenInvalid = NewWithStatus("token 错误", 400)
	// ErrUserGone ...
	ErrUserGone = NewWithStatus("用户已不存在", 410)
	// ErrInvalidStatus ...
	ErrInvalidStatus = NewWithStatus("状态错误", 400)
)

// Translator ...
type Translator interface {
	Translate(ut ut.Translator) string
}

// StatusCodeGetter ...
type StatusCodeGetter interface {
	HTTPStatusCode() int
}

// NotFoundChecker ...
type NotFoundChecker interface {
	IsNotFound() bool
}

// DetailGetter ...
type DetailGetter interface {
	Details() []string
}

// Error ...
type Error struct {
	status  int
	message string
	details []string
	err     error
}

// Errorf ...
func Errorf(format string, args ...interface{}) error {
	return &Error{
		message: fmt.Sprintf(format, args...),
		status:  500,
	}
}

// Wrap ...
func Wrap(e error, message string) error {
	status := 500
	if g, ok := e.(StatusCodeGetter); ok {
		status = g.HTTPStatusCode()
	}
	return &Error{
		message: e.Error(),
		status:  status,
		err:     e,
		details: []string{message},
	}
}

// WrapError ...
func WrapError(e error, detailErr error) error {
	status := 500
	if g, ok := e.(StatusCodeGetter); ok {
		status = g.HTTPStatusCode()
	}

	var details []string
	{
		if derr, ok := detailErr.(*Error); ok {
			details = append(details, derr.Details()...)
		}
	}

	{
		if de, ok := e.(*Error); ok {
			de.details = append(de.details, details...)
			de.err = detailErr
			return de
		}
	}

	return &Error{
		message: e.Error(),
		status:  status,
		details: details,
		err:     detailErr,
	}
}

// Wrapf ...
func Wrapf(e error, format string, args ...interface{}) error {
	status := 500
	if g, ok := e.(StatusCodeGetter); ok {
		status = g.HTTPStatusCode()
	}

	{
		if e, ok := e.(*Error); ok {
			e.details = append(e.details, fmt.Sprintf(format, args...))
			return e
		}
	}

	return &Error{
		message: e.Error(),
		status:  status,
		err:     e,
		details: []string{fmt.Sprintf(format, args...)},
	}
}

// New ...
func New(msg string) error {
	return &Error{
		message: msg,
		status:  500,
	}
}

// NewWithStatus ...
func NewWithStatus(msg string, status int) error {
	return &Error{
		message: msg,
		status:  status,
	}
}

// NewfWithStatus ...
func NewfWithStatus(status int, format string, args ...interface{}) error {
	return &Error{
		message: fmt.Sprintf(format, args...),
		status:  status,
	}
}

// NewWithStatusDetail ...
func NewWithStatusDetail(msg string, status int, details []string) error {
	return &Error{
		message: msg,
		details: details,
		status:  status,
	}
}

// Error ...
func (e Error) Error() string {
	return e.message
}

// Details ...
func (e Error) Details() []string {
	return e.details
}

// Unwrap ...
func (e Error) Unwrap() error {
	return e.err
}

// HTTPStatusCode ...
func (e Error) HTTPStatusCode() int {
	return e.status
}

// IsNotFound ...
func (e Error) IsNotFound() bool {
	return e.status == 404
}
