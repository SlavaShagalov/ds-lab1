package errors

import (
	"errors"
)

var (
	// Common repository
	ErrDb = errors.New("db error")

	// Persons
	ErrPersonNotFound      = errors.New("person not found")
	ErrPersonAlreadyExists = errors.New("person already exists")

	// HTTP
	ErrReadBody = errors.New("read request body error")
)
