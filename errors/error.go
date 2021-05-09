package errors

import (
	"github.com/pkg/errors"
)

var (
	Unknown  = errors.New("unknown error.")
	NotFound = errors.New("not found.")
	DBError  = errors.New("db unknown error.")
)
