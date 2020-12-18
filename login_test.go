package cascade

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/vitpelekhaty/httptracer"
)

var auth = Auth{
	Username: "username",
	Password: "password",
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

	err = conn.login(login, auth)

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
