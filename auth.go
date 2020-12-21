package cascade

import (
	"encoding/base64"
	"fmt"
)

// Auth параметры авторизации Каскада
type Auth struct {
	// Username имя пользователя
	Username string
	// Password пароль
	Password string
}

// Secret токен авторизации
func (a Auth) Secret() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", a.Username, a.Password)))
}
