package cascade

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Connection соединение с Каскадом
type Connection struct {
	uri     string
	authURI string
	client  *http.Client
	login   *LoginResponse
}

// NewConnection возвращает настроенное соединение с Каскадом
func NewConnection(uri string, client *http.Client) *Connection {
	return &Connection{
		uri:    uri,
		client: client,
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

// Login авторизация пользователя в Каскаде
func (self *Connection) Login(authURI string, auth Auth) error {
	if self.client == nil {
		return ErrorHTTPClientNotSpecified
	}

	if self.login != nil {
		return ErrorUserAuthorized
	}

	if _, err := url.Parse(authURI); err != nil {
		return err
	}

	form := url.Values{}
	form.Add("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", authURI, strings.NewReader(form.Encode()))

	if err != nil {
		return fmt.Errorf("POST %s: %q", authURI, err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth.Secret()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	resp, err := self.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("POST %s: %s", authURI, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("POST %s: %q", authURI, err)
	}

	var login LoginResponse

	err = json.Unmarshal(body, &login)

	if err != nil {
		return fmt.Errorf("POST %s: %q", authURI, err)
	}

	self.authURI = authURI
	self.login = &login

	return nil
}

// Logout закрытие сессии пользователя
func (self *Connection) Logout() error {
	// TODO: необходимо уточнить порядок завершения сессии пользователя
	if self.login == nil {
		return nil
	}

	/*
		req, err := http.NewRequest("DELETE", self.authURI, nil)

		if err != nil {
			return fmt.Errorf("DELETE %s: %q", self.authURI, err)
		}

		req.Header.Set("Authorization", fmt.Sprintf("%s %s", self.TokenType(), self.AccessToken()))

		resp, err := self.client.Do(req)

		if err != nil {
			return fmt.Errorf("DELETE %s: %q", self.authURI, err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("DELETE %s: %s", self.authURI, resp.Status)
		}
	*/

	self.login = nil

	return nil
}
