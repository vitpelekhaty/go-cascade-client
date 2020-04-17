// +build integration

package cascade

import (
	"encoding/json"
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

func TestConnection_CounterHouse_Real(t *testing.T) {
	done := false

	_, exec, _, ok := runtime.Caller(0)

	if !ok {
		t.FailNow()
	}

	tracedata := filepath.Join(filepath.Dir(exec), "/testdata/trace/counterHouse")

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

	done = true

	devices, err := conn.CounterHouse()

	if err != nil {
		t.Fatal(err)
	}

	if len(devices) == 0 {
		t.Error("CounterHouse() failed!")
	}
}

var beginAt = time.Date(2020, 4, 2, 1, 0, 0, 0, time.UTC)

func TestConnection_Readings_Real_HourArchive(t *testing.T) {
	done := false

	_, exec, _, ok := runtime.Caller(0)

	if !ok {
		t.FailNow()
	}

	tracedata := filepath.Join(filepath.Dir(exec), "/testdata/trace/readings_hours")

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

	ch, err := conn.CounterHouse()

	if err != nil {
		t.Fatal(err)
	}

	if len(ch) == 0 {
		t.Error("CounterHouse() failed!")
	}

	var devices []CounterHouseDto

	err = json.Unmarshal(ch, &devices)

	if err != nil {
		t.Fatal(err)
	}

	done = true

	device := devices[0]

	data, err := conn.Readings(device.ID, HourArchive, beginAt, beginAt.Add(time.Hour*24))

	if err != nil {
		t.Fatal(err)
	}

	if len(data) == 0 {
		t.Error("hours Readings() failed!")
	}
}

func TestConnection_Readings_Real_DailyArchive(t *testing.T) {
	done := false

	_, exec, _, ok := runtime.Caller(0)

	if !ok {
		t.FailNow()
	}

	tracedata := filepath.Join(filepath.Dir(exec), "/testdata/trace/readings_daily")

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

	ch, err := conn.CounterHouse()

	if err != nil {
		t.Fatal(err)
	}

	if len(ch) == 0 {
		t.Error("CounterHouse() failed!")
	}

	var devices []CounterHouseDto

	err = json.Unmarshal(ch, &devices)

	if err != nil {
		t.Fatal(err)
	}

	done = true

	device := devices[0]

	data, err := conn.Readings(device.ID, DailyArchive, beginAt, beginAt.Add(time.Hour*72))

	if err != nil {
		t.Fatal(err)
	}

	if len(data) == 0 {
		t.Error("daily Readings() failed!")
	}
}
