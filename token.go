package cascade

import (
	"time"
)

// token ответ сервера авторизации
type token struct {
	// value токен сессии
	value string `json:"access_token"`
	// tokenType тип токена (bearer etc)
	tokenType string `json:"token_type"`
	// expiresIn timestamp времени окончания действия токена
	expiresIn int64 `json:"expires_in"`
	// scope ???
	scope string `json:"scope"`
	// userID идентификатор пользователя в Каскаде
	userID int `json:"userid"`
	// user имя пользователя
	user string `json:"token"`
	// connectionName наименование соединения
	connectionName string `json:"name"`
	// serverType тип сервера (development etc)
	serverType string `json:"server_type"`
}

// expired возвращает время окончания действия токена
func (t *token) expired(loc *time.Location) time.Time {
	return time.Unix(t.expiresIn, 0).In(loc)
}
