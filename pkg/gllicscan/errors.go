package gllicscan

import "errors"

var (
	ErrFileDoesNotExist    = errors.New("file does not exist")
	ErrCannotReadFile      = errors.New("cannot read file")
	ErrCannotUnmarshalJSON = errors.New("unmarshal JSON")
	ErrCannotStatFile      = errors.New("cannot stat file")
)
