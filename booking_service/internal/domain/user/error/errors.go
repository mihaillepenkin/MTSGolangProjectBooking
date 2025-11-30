package error

import "errors"

var (
	ErrFailedToAuthorizeUser = errors.New("failed to authorize user")
	ErrUserValidationFailed  = errors.New("user validation failed")
)
