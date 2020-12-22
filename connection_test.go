// +build integration

package cascade

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/vitpelekhaty/httptracer"

	"github.com/vitpelekhaty/go-cascade-client/archive"
)

var (
	username           string
	password           string
	authURL            string
	cascadeURL         string
	insecureSkipVerify bool
	bodies             bool
	tracePath          string
	strTimeout         string
	limit              uint
	strStart           string
	strEnd             string
	strDataArchive     string
)

func init() {
	flag.StringVar(&username, "username", "", "username")
	flag.StringVar(&password, "password", "", "password")
	flag.StringVar(&cascadeURL, "api-url", "", "Cascade API URL")
	flag.StringVar(&authURL, "auth-url", "", "auth URL")
	flag.BoolVar(&insecureSkipVerify, "insecure-skip-verify", false, "insecure skip verify")
	flag.StringVar(&tracePath, "trace", "", "write trace into path")
	flag.StringVar(&strTimeout, "timeout", "30s", "timeout")
	flag.BoolVar(&bodies, "bodies", false, "write bodies into trace")
	flag.UintVar(&limit, "limit", 0, "limit number of devices")
	flag.StringVar(&strStart, "from", "", "a beginning of a measurement period")
	flag.StringVar(&strEnd, "to", "", "end of measurement period")
	flag.StringVar(&strDataArchive, "archive", "Hour", "type of archive")
}

const (
	layoutQuery = `02.01.2006 15`
)

func TestConnection(t *testing.T) {
	if strings.TrimSpace(cascadeURL) == "" {
		t.Fatal(errors.New("no API URL"))
	}

	if strings.TrimSpace(authURL) == "" {
		t.Fatal(errors.New("no auth URL"))
	}

	start, err := time.Parse(layoutQuery, strStart)

	if err != nil {
		t.Fatal(err)
	}

	end, err := time.Parse(layoutQuery, strEnd)

	if err != nil {
		t.Fatal(err)
	}

	archiveType, err := archive.Parse(strDataArchive)

	if err != nil {
		t.Fatal(err)
	}

	timeout, err := time.ParseDuration(strTimeout)

	if err != nil {
		t.Fatal(err)
	}

	client := setupHTTPClient(timeout*time.Second, insecureSkipVerify)

	if strings.TrimSpace(tracePath) != "" {
		f, err := os.Create(tracePath)

		if err != nil {
			t.Fatal(err)
		}

		defer func() {
			if err := f.WriteString("]"); err != nil {
				t.Error(err)
			}

			if err := f.Close(); err != nil {
				t.Error(err)
			}
		}()

		_, err = f.WriteString("[")

		if err != nil {
			t.Fatal(err)
		}

		emptyTrace := true

		callbackFunc := func(entry *httptracer.Entry) {
			if entry == nil {
				return
			}

			b, err := json.Marshal(entry)

			if err != nil {
				t.Fatal(err)
			}

			if !emptyTrace {
				_, err = tf.WriteString(",")

				if err != nil {
					t.Fatal(err)
				}
			}

			_, err = tf.Write(b)

			if err != nil {
				t.Fatal(err)
			}

			emptyTrace = false
		}

		client = setupTracer(client, setupTracerOptions(bodies, callbackFunc))
	}

	c, err := NewConnection(client)

	if err != nil {
		t.Fatal(err)
	}

	err = c.Open(cascadeURL, WithAuth(authURL, Auth{
		Username: username,
		Password: password,
	}))

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := c.Close(); err != nil {
			t.Error(err)
		}
	}()

	sensorBytes, err := c.CounterHouse()

	if err != nil {
		t.Fatal(err)
	}

	var sensorCount int

	sensors := ParseCounterHouseDto(sensorBytes)

	for sensor := range sensors {
		if sensor.error != nil {
			t.Fatal(sensor.error)
		}

		if limit > 0 && sensorCount > limit {
			continue
		}

		archiveBytes, err := c.Readings(sensor.ID, archiveType, start, end)

		if err != nil {
			t.Fatal(err)
		}

		readings := ParseCounterHouseReadingDto()

		for item := range readings {
			if item.error != nil {
				t.Fatal(item.error)
			}
		}

		sensorCount++
	}
}

func setupHTTPClient(timeout time.Duration, insecureSkipVerify bool) *http.Client {
	client := &http.Client{
		Timeout: timeout,
	}

	if insecureSkipVerify {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client.Transport = transport
	}
}

func setupTracer(client *http.Client, options ...httptracer.Option) *http.Client {
	return httptracer.Trace(client, options)
}

func setupTracerOptions(withBodies bool, withCallback httptracer.CallbackFunc) []httptracer.Option {
	options := make([]httptracer.Option, 0)

	if withBodies {
		options = append(options, httptracer.WithBodies(withBodies))
	}

	if withCallback != nil {
		options = append(options, httptracer.WithCallback(withCallback))
	}

	return options
}
