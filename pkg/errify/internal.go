package errify

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type InternalServerError struct {
	err      string
	location string
	details  string
	code     int
	read     string
}

func (m *InternalServerError) Error() string {
	return m.err
}

func (m *InternalServerError) Location() string {
	return m.location
}

func (m *InternalServerError) Details() string {
	return m.details
}

func (m *InternalServerError) Code() int {
	return http.StatusInternalServerError
}

func (m *InternalServerError) Read() string {
	if m.read == "" {
		m.read = uuid.NewString()
	}
	return m.read
}

func (m *InternalServerError) SetRead(read string) IError {
	m.read = read
	return m
}

func (m *InternalServerError) SetDetails(details string) IError {
	m.details = details
	return m
}

func (m *InternalServerError) JoinLoc(location string) IError {
	m.location = fmt.Sprintf("%s/%s", location, m.location)
	return m
}

func NewInternalServerError(err string, location string) IError {
	return &InternalServerError{
		err:      err,
		location: fmt.Sprintf("(%s)", location),
	}
}

func (m *InternalServerError) Message() (msg string) { return }
