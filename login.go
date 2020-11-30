package cascade

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Login авторизация пользователя в Каскаде
func (c *Connection) Login(authURL string, auth Auth) error {
	if c.client == nil {
		return fmt.Errorf("POST %s: %v", authURL, errors.New("no HTTP client"))
	}

	if c.login != nil {
		return fmt.Errorf("POST %s: %v", authURL, errors.New("user is already authorized"))
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

	resp, err := c.client.Do(req)

	if err != nil {
		return fmt.Errorf("POST %s: %v", authURL, err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.errorCallbackFunc(err)
		}
	}()

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

	c.authURL = authURL
	c.login = &login

	return nil
}

// Logout закрытие сессии пользователя
func (c *Connection) Logout() error {
	if c.login == nil {
		return nil
	}

	c.login = nil

	return nil
}
