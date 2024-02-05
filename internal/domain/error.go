package domain

import "errors"

var (
	ErrNotFound = errors.New("proxy not found")
	ErrOnCreate = errors.New("can't create")
	ErrOnUpdate = errors.New("can't update")
	ErrOnDelete = errors.New("can't delete")
)
