package error

import "errors"

var (
	ErrPaymentTokenIsInvalid = errors.New("http token is invalid")
	ErrPaymentIsInvalid      = errors.New("http is invalid")
	ErrPaymentFailed         = errors.New("http failed")
)
