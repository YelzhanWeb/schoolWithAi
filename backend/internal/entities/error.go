package entities

import "errors"

var (
	ErrAlreadyExists      = errors.New("already exists")
	ErrNotFound           = errors.New("not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
