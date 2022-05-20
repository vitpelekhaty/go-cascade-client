package cascade

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
)

// pathJoin возвращает полный URL метода API
func pathJoin(u, p string) (string, error) {
	parsedURL, err := url.Parse(u)

	if err != nil {
		return p, err
	}

	parsedURL.Path = path.Join(parsedURL.Path, p)

	return parsedURL.String(), nil
}

// secret возвращает шифрованные параметры для базовой авторизации
func secret(username, passwd string) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, passwd)))
}
