package usecase

import "errors"

var (
	ErrNotFound    = errors.New("proxy not found")
	ErrInRepo      = errors.New("error in repo")
	ErrInvalidData = errors.New("invalid data")
)
