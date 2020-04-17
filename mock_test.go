package cascade

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var responses = map[string]string{
	Login:             "/testdata/responses/login.json",
	CounterHouse:      "/testdata/responses/counterHouse.json",
	Readings:          "/testdata/responses/readings200.json",
	Readings + " 422": "/testdata/responses/readings422.json",
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

	// /api/cascade/counter-house/reading
	if strings.HasPrefix(r.RequestURI, Readings) && r.Method == "POST" {
		token, err := accessToken()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		path = filepath.Join(filepath.Dir(exec), responses[Readings])

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

		data, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var query ReadingsRequest

		err = json.Unmarshal(data, &query)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		beginAt := time.Time(query.BeginAt)
		endAt := time.Time(query.EndAt)

		if endAt.Sub(beginAt) > (time.Hour * 24 * 7) {
			path = filepath.Join(filepath.Dir(exec), responses[Readings+" 422"])

			data, _ = ioutil.ReadFile(path)

			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write(data)

			return
		}

		data, err = ioutil.ReadFile(path)

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
