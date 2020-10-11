package userstore

import (
	"errors"
)

// ErrAlreadyExist is returned if the email is already present in the store
var ErrAlreadyExist = errors.New("user already exist")

// ErrNotFound is returned if the id of the user isn't found in the store
var ErrNotFound = errors.New("user not found")

// ErrQueryNotCompatible is returned if the email is already present in the store
var ErrQueryNotCompatible = errors.New("the provided query is not compatible")
