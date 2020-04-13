package cascade

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/vitpelekhaty/httptracer"
)

const (
	username = `sibgenco`
	password = `AmXn9GZZd@#xGB!`
	authURI  = `http://10.253.128.1/oauth/token`
	uri      = `http://10.253.128.1/api/cascade`
)

func TestConnection_LoginLogout(t *testing.T) {
	done := false

	_, exec, _, ok := runtime.Caller(0)

	if !ok {
		t.FailNow()
	}

	tracedata := filepath.Join(filepath.Dir(exec), "/testdata/trace/login")

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

	conn := NewConnection(uri, client)

	defer func() {
		done = true

		if err := conn.Logout(); err != nil {
			t.Error(err)
		}
	}()

	err = conn.Login(authURI, Auth{Username: username, Password: password})

	if err != nil {
		t.Fatal(err)
	}
}
