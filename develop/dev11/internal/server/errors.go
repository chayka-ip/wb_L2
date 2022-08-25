package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	errMethodNotAllowed = errors.New("method is not allowed")
	errBadRequestData   = errors.New("invalid request data")
)

const (
	buisnessLogicErrorCode  = http.StatusServiceUnavailable
	badRequestErrorCode     = http.StatusBadRequest
	internalServerErrorCode = http.StatusInternalServerError
)

func errNotEnoughArguments(expectedArgs ...string) error {
	return fmt.Errorf("Not enough arguments. Expected at least %s", strings.Join(expectedArgs, ", "))
}

func errInvalidArgument(name string) error {
	return fmt.Errorf("argument is invalid: %s", name)
}
