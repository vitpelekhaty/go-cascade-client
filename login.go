package cascade

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Login авторизация пользователя в Каскаде
func (self *Connection) Login(authURL string, auth Auth) error {
	if self.client == nil {
		return fmt.Errorf("POST %s: %v", authURL, ErrorHTTPClientNotSpecified)
	}

	if self.login != nil {
		return fmt.Errorf("POST %s: %v", authURL, ErrorUserAuthorized)
	}

	if _, err := url.Parse(authURL); err != nil {
		return fmt.Errorf("POST %s: %v", authURL, err)
	}

	form := url.Values{}
	form.Add("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", authURL, strings.NewReader(form.Encode()))

	if err != nil {
		return fmt.Errorf("POST %s: %v", authURL, err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth.Secret()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	resp, err := self.client.Do(req)

	if err != nil {
		return fmt.Errorf("POST %s: %v", authURL, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("POST %s: %s", authURL, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("POST %s: %v", authURL, err)
	}

	var login LoginResponse

	err = json.Unmarshal(body, &login)

	if err != nil {
		return fmt.Errorf("POST %s: %v", authURL, err)
	}

	self.authURL = authURL
	self.login = &login

	return nil
}

// Logout закрытие сессии пользователя
func (self *Connection) Logout() error {
	if self.login == nil {
		return nil
	}

	self.login = nil

	return nil
}
