package cascade

import (
	"strconv"
	"strings"
)

// Error ошибка метода API Каскад
type Error struct {
	exception   string
	message     string
	description string
	status      string
	path        string
	method      string
	statusCode  int
}

// NewCascadeError возвращает новый экземпляр ошибки выполнения метода API
func NewCascadeError(err *errorMessage, method, path string, statusCode int) *Error {
	e := &Error{
		exception:   err.Exception,
		message:     err.Message,
		description: err.Description,
		status:      err.Status,
		path:        path,
		method:      method,
		statusCode:  statusCode,
	}

	if err.Error != "" {
		e.message = err.Error
	}

	return e
}

// Error реализация интерфейса error
func (e *Error) Error() string {
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

// Exception тип исключения
func (e *Error) Exception() string {
	return e.exception
}

// Message сообщение об ошибке
func (e *Error) Message() string {
	return e.message
}

// ExceptionStatus состояние исключения
func (e *Error) ExceptionStatus() string {
	return e.status
}

// Description описание ошибки
func (e *Error) Description() string {
	return e.description
}

// Path метод API, вернувший ошибку
func (e *Error) Path() string {
	return e.path
}

// Method метод HTTP, использованный при вызове метода API
func (e *Error) Method() string {
	return e.method
}

// StatusCode HTTP код результата вызова метода API
func (e *Error) StatusCode() int {
	return e.statusCode
}

// message сообщение сервера об ошибке
type errorMessage struct {
	// Message текст ошибки
	Message string `json:"message"`

	// Description описание ошибки
	Description string `json:"description"`

	// Error текст ошибки
	Error string `json:"error"`

	// Exception тип исключения
	Exception string `json:"exception"`

	// Status статус исключения
	Status string `json:"status"`
}
