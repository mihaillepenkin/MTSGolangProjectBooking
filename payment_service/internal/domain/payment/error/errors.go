package error

import "errors"

var (
	ErrPaymentTokenIsInvalid = errors.New("payment token is invalid")
	ErrPaymentIsInvalid      = errors.New("payment is invalid")
	ErrPaymentFailed         = errors.New("payment failed")
)
