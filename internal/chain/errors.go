package chain

import (
	"errors"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrBadPassword = errors.New("invalid password")

	ErrInvalidToken            = errors.New("ivalid token")
	ErrKeyNotFound             = errors.New("key not found")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrExpired                 = errors.New("expired")
)
