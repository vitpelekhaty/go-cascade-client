package cascade

import (
	"errors"
)

// ErrorHTTPClientNotSpecified ошибка "Не указан HTTP клиент"
var ErrorHTTPClientNotSpecified = errors.New("HTTP client not specified")

// ErrorUserAuthorized ошибка "Пользователь уже авторизован"
var ErrorUserAuthorized = errors.New("user is already authorized")
