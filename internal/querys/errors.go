package querys

import (
	"errors"
)

var (
	ErrNoRecord = errors.New("models: no matching found")

	ErrInvalidCredentials = errors.New("models: invalid credential")

	ErrDuplicateEmail = errors.New("models: duplicates email")
)