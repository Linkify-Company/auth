package errify

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type UnauthorizedError struct {
	err      string
	message  string
	location string
	details  string
	read     string
}

func (m *UnauthorizedError) Error() string {
	return m.err
}

func (m *UnauthorizedError) Message() string {
	return m.message
}

func (m *UnauthorizedError) Location() string {
	return m.location
}

func (m *UnauthorizedError) Details() string {
	return m.details
}

func (m *UnauthorizedError) Code() int {
	return http.StatusUnauthorized
}

func (m *UnauthorizedError) Read() string {
	if m.read == "" {
		m.read = uuid.NewString()
	}
	return m.read
}

func (m *UnauthorizedError) SetDetails(details string) IError {
	m.details = details
	return m
}

func (m *UnauthorizedError) JoinLoc(location string) IError {
	m.location = fmt.Sprintf("%s/%s", location, m.location)
	return m
}

func NewUnauthorizedError(err string, message string, location string) IError {
	return &UnauthorizedError{
		err:      err,
		message:  message,
		location: fmt.Sprintf("(%s)", location),
	}
}
