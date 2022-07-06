package rudp

import (
	"errors"
)

var (
	ErrMaxMTU      = errors.New("packet size larger than MTU")
	ErrWriteClosed = errors.New("attempt to write to closed connection")
	ErrReadClosed  = errors.New("attempt to read from closed connection")
)
