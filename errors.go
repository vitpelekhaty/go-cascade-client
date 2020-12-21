package cascade

import (
	"strconv"
	"strings"
)

// ErrCascadeCall
type ErrCascadeCall struct {
	exception   string
	message     string
	description string
	status      string
	path        string
	method      string
	statusCode  int
}

// Error реализация интерфейса error
func (e *ErrCascadeCall) Error() string {
	var builder strings.Builder

	if strings.TrimSpace(e.exception) != "" {
		builder.WriteString("exception " + e.exception)

		if strings.TrimSpace(e.status) != "" {
			builder.WriteString(" (")
			builder.WriteString(e.status)
			builder.WriteString(")")
		}
	} else {
		builder.WriteString("error")
	}

	if strings.TrimSpace(e.method) != "" || strings.TrimSpace(e.path) != "" {
		builder.WriteString(" during call")

		if strings.TrimSpace(e.method) != "" {
			builder.WriteString(" " + e.method)
		}

		if strings.TrimSpace(e.path) != "" {
			builder.WriteString(" " + e.path)
		}

		if e.statusCode > 0 {
			builder.WriteString(" (" + strconv.Itoa(e.statusCode) + ")")
		}
	}

	if strings.TrimSpace(e.message) != "" {
		if builder.Len() > 0 {
			builder.WriteRune(':')
		}

		builder.WriteString(" " + e.message)

		if strings.TrimSpace(e.description) != "" {
			builder.WriteString(" (" + e.description + ")")
		}
	}

	if builder.Len() == 0 {
		builder.WriteString("unknown error")
	}

	return builder.String()
}

// Exception исключение API
func (e *ErrCascadeCall) Exception() string {
	return e.exception
}

// Message сообщение об ошибке API
func (e *ErrCascadeCall) Message() string {
	return e.message
}

// ExceptionStatus состояние API
func (e *ErrCascadeCall) ExceptionStatus() string {
	return e.status
}

// Description описание ошибки API
func (e *ErrCascadeCall) Description() string {
	return e.description
}

// Path метод API, вернувший ошибку
func (e *ErrCascadeCall) Path() string {
	return e.path
}

// Method метод HTTP, использованный при вызове метода API
func (e *ErrCascadeCall) Method() string {
	return e.method
}

// StatusCode HTTP код результата вызова метода API
func (e *ErrCascadeCall) StatusCode() int {
	return e.statusCode
}
