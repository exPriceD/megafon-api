package megafon

import "errors"

var (
	ErrUnauthorized = errors.New("megafon: unauthorized")
	ErrRateLimited  = errors.New("megafon: rate limited")
)
