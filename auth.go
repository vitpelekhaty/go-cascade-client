package cascade

import (
	"encoding/base64"
	"fmt"
	"time"
)

// Auth параметры авторизации Каскада
type Auth struct {
	// Username имя пользователя
	Username string
	// Password пароль
	Password string
}

// Secret токен авторизации
func (self Auth) Secret() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", self.Username, self.Password)))
}

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
func (self *token) expired(loc *time.Location) time.Time {
	return time.Unix(self.expiresIn, 0).In(loc)
}
