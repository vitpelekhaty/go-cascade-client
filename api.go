package cascade

import (
	"net/url"
	"path"
)

var (
	// Login метод авторизации
	Login = "/oauth/token"
)

// URLJoin возвращает полный URI метода API
func URLJoin(baseURL, method string) (string, error) {
	u, err := url.Parse(baseURL)

	if err != nil {
		return method, err
	}

	u.Path = path.Join(u.Path, method)

	return u.String(), nil
}
