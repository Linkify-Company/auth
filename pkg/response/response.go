package response

import (
	"auth/pkg/errify"
	"auth/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"
)

type Send struct {
	Value   any    `json:"value"`
	Message string `json:"message"`
	code    int
}

type SendError struct {
	Error string `json:"error"`
	Read  string `json:"read"`
}

func NewSend(value any, message string, code int) Send {
	return Send{
		Value:   value,
		Message: message,
		code:    code,
	}
}

func sendError(err errify.IError) SendError {
	switch err.(type) {
	case *errify.UnauthorizedError:
		return SendError{
			Error: "Unauthorized",
			Read:  err.Read(),
		}
	case *errify.BadRequestError:
		return SendError{
			Error: err.Message(),
			Read:  err.Read(),
		}
	case *errify.InternalServerError:
		return SendError{
			Error: "Internal Server Error",
			Read:  err.Read(),
		}
	}
	return SendError{}
}

func Ok(w http.ResponseWriter, send Send, log logger.Logger) {
	w.Header().Set("Content-Type", "application/json")

	jsonBytes, err := json.MarshalIndent(send, "", "    ")
	if err != nil {
		log.Error(errify.NewInternalServerError(err.Error(), "Ok/MarshalIndent"))
		return
	}

	log.Infof(fmt.Sprintln(jsonBytes))

	w.WriteHeader(send.code)
	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Error(errify.NewInternalServerError(err.Error(), "Ok/Write"))
		return
	}
}

func Error(w http.ResponseWriter, err errify.IError, log logger.Logger) {
	w.Header().Set("Content-Type", "application/json")

	jsonBytes, e := json.MarshalIndent(sendError(err), "", "    ")
	if e != nil {
		log.Error(errify.NewInternalServerError(e.Error(), "Error/MarshalIndent"))
		return
	}
	log.Error(err)

	w.WriteHeader(err.Code())
	_, e = w.Write(jsonBytes)
	if e != nil {
		log.Error(errify.NewInternalServerError(e.Error(), "Error/Write"))
		return
	}
}
