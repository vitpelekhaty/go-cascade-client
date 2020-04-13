package cascade

import (
	"net/http"
)

// Connection соединение с Каскадом
type Connection struct {
	baseURL string
	authURL string
	client  *http.Client
	login   *LoginResponse
}

// NewConnection возвращает настроенное соединение с Каскадом
func NewConnection(baseURL string, client *http.Client) *Connection {
	return &Connection{
		baseURL: baseURL,
		client:  client,
	}
}

// Connected возвращает признак установленного соединения
func (self *Connection) Connected() bool {
	return self.login != nil
}

// AccessToken возвращает токен сессии
func (self *Connection) AccessToken() string {
	if self.login == nil {
		return ""
	}

	return self.login.AccessToken
}

// TokenType возвращает тип токена
func (self *Connection) TokenType() string {
	if self.login == nil {
		return ""
	}

	return self.login.TokenType
}

func (self *Connection) checkConnection() error {
	if self.login == nil {
		return ErrorUserUnauthorized
	}

	if self.client == nil {
		return ErrorHTTPClientNotSpecified
	}

	return nil
}
