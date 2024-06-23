package resolvers

import (
	"errors"
	"fmt"
)

func ErrNotFound(resource string) error {
	errNotFound := fmt.Sprintf("resource %s not found", resource)
	return errors.New(errNotFound)
}

func ErrInvalidInput(field string, reason string) error {
	errInvalidInput := fmt.Sprintf("invalid input provided for field %s: %s", field, reason)
	return errors.New(errInvalidInput)
}

func ErrStringInput(field string) error {
	return ErrInvalidInput(field, "must be a STRING")
}

func ErrArrayInput(field string) error {
	return ErrInvalidInput(field, "must be an ARRAY")
}

func ErrBooleanInput(field string) error {
	return ErrInvalidInput(field, "must be an BOOLEAN")
}

func ErrExistInput(field string) error {
	return ErrInvalidInput(field, "DOESN'T EXIST")
}

var (
	ErrInternal         = errors.New("internal error, please try again later")
	ErrPermission       = errors.New("permission denied")
	ErrClassroomNotPaid = errors.New("classroom not paid")
)
