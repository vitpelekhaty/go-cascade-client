// +build integration

package cascade

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/vitpelekhaty/httptracer"
)

var (
	username   string
	password   string
	authURL    string
	cascadeURL string
)

func init() {
	flag.StringVar(&username, "username", "", "username")
	flag.StringVar(&password, "password", "", "password")
	flag.StringVar(&cascadeURL, "api-url", "", "Cascade API URL")
	flag.StringVar(&authURL, "auth-url", "", "Auth URL")
}

func TestConnection_LoginLogout_Real(t *testing.T) {
	done := true

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

	conn := NewConnection(cascadeURL, client)

	defer func() {
		done = true

		if err := conn.Logout(); err != nil {
			t.Error(err)
		}
	}()

	err = conn.Login(authURL, Auth{Username: username, Password: password})

	if err != nil {
		t.Fatal(err)
	}
}
