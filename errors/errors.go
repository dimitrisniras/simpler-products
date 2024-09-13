package errors

import "errors"

var (
	ErrProductNotFound        = errors.New("product not found")
	ErrInvalidProductID       = errors.New("invalid product id")
	ErrInvalidLimitParameter  = errors.New("invalid limit parameter, limit must be in the range of [1, 100]")
	ErrInvalidOffsetParameter = errors.New("invalid offset parameter, offest must be a positive number")
)
