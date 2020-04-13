package cascade

import (
	"net/url"
	"path"
)

var (
	// Login метод авторизации
	Login = "/oauth/token"
	// CounterHouse метод получения списка приборов учета
	CounterHouse = "/api/cascade/counter-house"
	// Readings метод чтения архива показаний прибора учета
	Readings = "/api/cascade/counter-house/reading"
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
