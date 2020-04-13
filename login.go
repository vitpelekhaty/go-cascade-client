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
		return ErrorHTTPClientNotSpecified
	}

	if self.login != nil {
		return ErrorUserAuthorized
	}

	if _, err := url.Parse(authURL); err != nil {
		return err
	}

	form := url.Values{}
	form.Add("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", authURL, strings.NewReader(form.Encode()))

	if err != nil {
		return fmt.Errorf("POST %s: %q", authURL, err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth.Secret()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	resp, err := self.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("POST %s: %s", authURL, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("POST %s: %q", authURL, err)
	}

	var login LoginResponse

	err = json.Unmarshal(body, &login)

	if err != nil {
		return fmt.Errorf("POST %s: %q", authURL, err)
	}

	self.authURL = authURL
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
		req, err := http.NewRequest("DELETE", self.authURL, nil)

		if err != nil {
			return fmt.Errorf("DELETE %s: %q", self.authURL, err)
		}

		req.Header.Set("Authorization", fmt.Sprintf("%s %s", self.TokenType(), self.AccessToken()))

		resp, err := self.client.Do(req)

		if err != nil {
			return fmt.Errorf("DELETE %s: %q", self.authURL, err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("DELETE %s: %s", self.authURL, resp.Status)
		}
	*/

	self.login = nil

	return nil
}
