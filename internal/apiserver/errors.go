package apiserver

import (
	"database/sql"
	"net/http"
)

/*
These types are used to generalize errors

	and then used as responses to client
*/
type ApiErrorReply interface {
	Error() string
}

type ReceiveError struct {
	Err string `json:"error"`
}

func (e *ReceiveError) Error() string {
	return e.Err
}

func ReceiveErr(err error) error {
	return &ReceiveError{Err: err.Error()}
}

type ReplyError struct {
	Err string `json:"error"`
}

func (e *ReplyError) Error() string {
	return e.Err
}

func ReplyErr(err error) error {
	return &ReplyError{Err: err.Error()}
}

// NotFoundError represents a 404 Not Found error
type NotFoundError struct {
	Err string `json:"error"`
}

func (e *NotFoundError) Error() string {
	return e.Err
}

func NotFoundErr(err error) error {
	return &NotFoundError{Err: err.Error()}
}

// NoContentError represents a 204 No Content status
type NoContentError struct {
	Err string `json:"error"`
}

func NoContentErr(err error) error {
	return &NoContentError{Err: err.Error()}
}

func (e *NoContentError) Error() string {
	return e.Err
}

// BadRequestError represents a 400 Bad Request error
type BadRequestError struct {
	Err string `json:"error"`
}

func (e *BadRequestError) Error() string {
	return e.Err
}

func BadRequestErr(err error) error {
	return &BadRequestError{Err: err.Error()}
}

// InternalServerError represents a 500 Internal Server Error
type InternalServerError struct {
	Err string `json:"error"`
}

func (e *InternalServerError) Error() string {
	return e.Err
}

func InternalServerErr(err error) error {
	return &InternalServerError{Err: err.Error()}
}

func errToStatusCode(err ApiErrorReply) int {
	switch err.(type) {
	case *ReplyError:
		return http.StatusInternalServerError
	case *NoContentError:
		return http.StatusNoContent
	case *NotFoundError:
		return http.StatusNotFound
	case *ReceiveError:
		return http.StatusBadRequest
	case *BadRequestError:
		return http.StatusBadRequest
	case *InternalServerError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func parseDbErr(err error) error {
	switch err {
	case sql.ErrNoRows:
		return NoContentErr(err)
	case sql.ErrConnDone:
		fallthrough
	case sql.ErrTxDone:
		return InternalServerErr(err)
	}
	return InternalServerErr(err)
}
