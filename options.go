package cascade

import (
	"net/http"
)

type connOptions struct {
	client *http.Client
}

type openOptions struct {
	authURL string
}

// Option опция соединения с API Каскад
type Option func(options *connOptions)

// WithHTTPClient устанавливает пользовательский экзепляр HTTP клиента взамен клиента по умолчанию
func WithHTTPClient(client *http.Client) Option {
	return func(options *connOptions) {
		options.client = client
	}
}

// OpenOption опция открытия соединения с API Каскад
type OpenOption func(options *openOptions)

// WithAuthURL устанавливает альтернативный URL авторизации в API Каскад
func WithAuthURL(authURL string) OpenOption {
	return func(options *openOptions) {
		options.authURL = authURL
	}
}
