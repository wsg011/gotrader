package types

import "errors"

var (
	ErrRateLimit  = errors.New("rate limit")
	ErrRateBanned = errors.New("rate banned")
)
