package customErrors

import "errors"

var ErrTooManyRedirects = errors.New("[CRAWLER] : Too many redirects")
var ErrMaxPagesReached = errors.New("[CRAWLER] : Max pages per run reached")
var ErrBodyTooLarge = errors.New("[CRAWLER] : Body too large")
