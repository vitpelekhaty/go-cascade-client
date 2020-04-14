package cascade

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
)

var responses = map[string]string{
	Login:        "/testdata/responses/login.json",
	CounterHouse: "/testdata/responses/counterHouse.json",
}

func accessToken() (string, error) {
	_, exec, _, ok := runtime.Caller(0)

	if !ok {
		return "", errors.New("undefined testdata path")
	}

	path := filepath.Join(filepath.Dir(exec), responses[Login])

	data, err := ioutil.ReadFile(path)

	if err != nil {
		return "", err
	}

	var login LoginResponse

	err = json.Unmarshal(data, &login)

	if err != nil {
		return "", err
	}

	return login.AccessToken, nil
}

func MockServerFunc(w http.ResponseWriter, r *http.Request) {
	_, exec, _, ok := runtime.Caller(0)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var path string

	// /oauth/token
	if strings.HasPrefix(r.RequestURI, Login) && r.Method == "POST" {
		path = filepath.Join(filepath.Dir(exec), responses[Login])

		authHeader := r.Header.Get("Authorization")
		values := strings.Split(authHeader, " ")

		if len(values) < 2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		credentials := values[1]

		if credentials != auth.Secret() {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		data, err := ioutil.ReadFile(path)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)

		return
	}

	// /api/cascade/counter-house
	if strings.HasPrefix(r.RequestURI, CounterHouse) && r.Method == "GET" {
		token, err := accessToken()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		path = filepath.Join(filepath.Dir(exec), responses[CounterHouse])

		authHeader := r.Header.Get("Authorization")
		values := strings.Split(authHeader, " ")

		if len(values) < 2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		credentials := values[1]

		if credentials != token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		data, err := ioutil.ReadFile(path)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)

		return
	}

	w.WriteHeader(http.StatusBadRequest)
}
