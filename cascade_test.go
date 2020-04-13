package cascade

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/vitpelekhaty/httptracer"
)

var responses = map[string]string{
	Login: "/testdata/responses/login.json",
}

var auth = Auth{
	Username: "username",
	Password: "password",
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

func TestConnection_LoginLogout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(MockServerFunc))
	defer ts.Close()

	done := true

	_, exec, _, ok := runtime.Caller(0)

	if !ok {
		t.FailNow()
	}

	tracedata := filepath.Join(filepath.Dir(exec), "/testdata/trace/login_test")

	f, err := os.Create(tracedata)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		f.WriteString("]")
		f.Close()
	}()

	f.WriteString("[")

	client := &http.Client{Timeout: time.Second * 10}
	client = httptracer.Trace(client, httptracer.WithBodies(true), httptracer.WithWriter(f),
		httptracer.WithCallback(func(entry *httptracer.Entry) {
			if !done {
				if entry != nil {
					f.WriteString(",")
				}
			}
		}))

	conn := NewConnection(ts.URL, client)

	defer func() {
		done = true

		if err := conn.Logout(); err != nil {
			t.Error(err)
		}
	}()

	login, err := URLJoin(ts.URL, Login)

	if err != nil {
		t.Fatal(err)
	}

	err = conn.Login(login, auth)

	if err != nil {
		t.Fatal(err)
	}

	if conn.AccessToken() == "" {
		t.Error("unauthorized")
	}

	if conn.TokenType() == "" {
		t.Error("unauthorized")
	}
}
