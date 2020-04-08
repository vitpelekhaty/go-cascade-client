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

// LoginResponse ответ сервера авторизации
type LoginResponse struct {
	// AccessToken токен сессии
	AccessToken string `json:"access_token"`
	// TokenType тип токена (bearer etc)
	TokenType string `json:"token_type"`
	// ExpiresIn timestamp времени окончания действия токена
	ExpiresIn int64 `json:"expires_in"`
	// Scope ???
	Scope string `json:"scope"`
	// UserID идентификатор пользователя в Каскаде
	UserID int `json:"userid"`
	// Login имя пользователя
	Login string `json:"login"`
	// Name наименование соединения
	Name string `json:"name"`
	// ServerType тип сервера (development etc)
	ServerType string `json:"server_type"`
}

// Expired возвращает время окончания действия токена
func (self *LoginResponse) Expired(loc *time.Location) time.Time {
	return time.Unix(self.ExpiresIn, 0).In(loc)
}
