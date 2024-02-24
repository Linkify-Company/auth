package errify

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type BadRequestError struct {
	err      string
	message  string
	location string
	details  string
	read     string
}

func (m *BadRequestError) Error() string {
	return m.err
}

func (m *BadRequestError) Message() string {
	return m.message
}

func (m *BadRequestError) Location() string {
	return m.location
}

func (m *BadRequestError) Details() string {
	return m.details
}

func (m *BadRequestError) Code() int {
	return http.StatusBadRequest
}

func (m *BadRequestError) Read() string {
	if m.read == "" {
		m.read = uuid.NewString()
	}
	return m.read
}

func (m *BadRequestError) SetDetails(details string) IError {
	m.details = details
	return m
}

func (m *BadRequestError) JoinLoc(location string) IError {
	m.location = fmt.Sprintf("%s/%s", location, m.location)
	return m
}

func NewBadRequestError(err string, message string, location string) IError {
	return &BadRequestError{
		err:      err,
		message:  message,
		location: fmt.Sprintf("(%s)", location),
	}
}
