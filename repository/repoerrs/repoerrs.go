package repoerrs

import "errors"

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotExists     = errors.New("not exists")
)
