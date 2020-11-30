package cascade

import (
	"errors"
	"net/http"
)

// Connection соединение с Каскадом
type Connection struct {
	baseURL string
	authURL string
	client  *http.Client
	login   *LoginResponse

	OnError func(err error)
}

// NewConnection возвращает настроенное соединение с Каскадом
func NewConnection(baseURL string, client *http.Client) *Connection {
	return &Connection{
		baseURL: baseURL,
		client:  client,
	}
}

// Connected возвращает признак установленного соединения
func (c *Connection) Connected() bool {
	return c.login != nil
}

// AccessToken возвращает токен сессии
func (c *Connection) AccessToken() string {
	if c.login == nil {
		return ""
	}

	return c.login.AccessToken
}

// TokenType возвращает тип токена
func (c *Connection) TokenType() string {
	if c.login == nil {
		return ""
	}

	return c.login.TokenType
}

func (c *Connection) checkConnection() error {
	if c.login == nil {
		return errors.New("user not authorized")
	}

	if c.client == nil {
		return errors.New("no HTTP client")
	}

	return nil
}

func (c *Connection) errorCallbackFunc(err error) {
	if err != nil && c.OnError != nil {
		c.OnError(err)
	}
}
