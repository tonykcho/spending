package utils

import "errors"

var ErrNotFound = errors.New("not found")

var ErrConflict = errors.New("conflict")

var ErrInvalidInput = errors.New("invalid input")

var ErrResourceExists = errors.New("resource already exists")

func MapErrorToStatusCode(err error) int {
	if errors.Is(err, ErrNotFound) {
		return 404
	} else if errors.Is(err, ErrConflict) {
		return 409
	} else if errors.Is(err, ErrInvalidInput) {
		return 400
	} else {
		return 500
	}
}
