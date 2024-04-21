package filesystem

import "errors"

var ErrInvalidCharacter = errors.New("invalid input string")
var ErrDataNotFound = errors.New("not exist")
var ErrDataAlreadyExists = errors.New("already exists")
