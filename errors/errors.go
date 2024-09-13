package errors

import "errors"

var (
	ErrProductNotFound            = errors.New("product not found")
	ErrInvalidProductID           = errors.New("invalid product id")
	ErrInvalidLimitParameter      = errors.New("invalid limit parameter, limit must be in the range of [1, 100]")
	ErrInvalidOffsetParameter     = errors.New("invalid offset parameter, offest must be a positive number")
	ErrAuthorizationHeaderMissing = errors.New("authorization header is missing")
	ErrAuthorizationHeaderFormat  = errors.New("invalid Authorization header format")
	ErrDecodingPublicKey          = errors.New("error decoding public key")
	ErrParsingPublicKey           = errors.New("error parsing public key")
	ErrInvalidToken               = errors.New("invalid token")
)
