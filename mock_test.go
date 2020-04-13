package cascade

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
)

var responses = map[string]string{
	Login: "/testdata/responses/login.json",
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

	w.WriteHeader(http.StatusBadRequest)
}
