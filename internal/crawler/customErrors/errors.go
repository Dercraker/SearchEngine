package customErrors

import "errors"

var ErrMaxPagesReached = errors.New("[CRAWLER] : Max pages per run reached")
