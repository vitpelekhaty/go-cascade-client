package cascade

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// methodLogin авторизация пользователя в Каскаде
func (c *Connection) login(ctx context.Context, authURL string, secret string) error {
	if c.token != nil {
		return fmt.Errorf("POST %s: %v", authURL, errors.New("user is already authorized"))
	}

	if _, err := url.Parse(authURL); err != nil {
		return fmt.Errorf("POST %s: %v", authURL, err)
	}

	form := url.Values{}
	form.Add("grant_type", "client_credentials")

	req, err := http.NewRequestWithContext(ctx, "POST", authURL, strings.NewReader(form.Encode()))

	if err != nil {
		return fmt.Errorf("POST %s: %v", authURL, err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", secret))
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

	var login token

	err = json.Unmarshal(body, &login)

	if err != nil {
		return fmt.Errorf("POST %s: %v", authURL, err)
	}

	c.authURL = authURL
	c.token = &login

	return nil
}
