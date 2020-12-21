package cascade

// token ответ сервера авторизации
type token struct {
	// Value токен сессии
	Value string `json:"access_token"`
	// Type тип токена (bearer etc)
	Type string `json:"token_type"`
	// ExpiresIn timestamp времени окончания действия токена
	ExpiresIn int64 `json:"expires_in"`
	// Scope ???
	Scope string `json:"Scope"`
	// UserID идентификатор пользователя в Каскаде
	UserID int `json:"userid"`
	// User имя пользователя
	User string `json:"token"`
	// Connection наименование соединения
	Connection string `json:"name"`
	// ServerType тип сервера (development etc)
	ServerType string `json:"server_type"`
}
