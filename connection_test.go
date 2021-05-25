// +build integration

package cascade

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dustin/go-humanize/english"
	"github.com/vitpelekhaty/httptracer"

	"github.com/vitpelekhaty/go-cascade-client/archive"
)

var (
	username           string
	password           string
	authURL            string
	cascadeURL         string
	insecureSkipVerify bool
	separatedInputs    bool
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
	flag.BoolVar(&separatedInputs, "separated-inputs", false, "use separated inputs")
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

	aURL, err := url.Parse(authURL)

	if err != nil {
		t.Fatal(err)
	}

	start, err := time.Parse(layoutQuery, strStart)

	if err != nil {
		t.Fatal(err)
	}

	end, err := time.Parse(layoutQuery, strEnd)

	if err != nil {
		t.Fatal(err)
	}

	archiveType := archive.Parse(strDataArchive)

	if archiveType == archive.UnknownArchive {
		t.Fatal(fmt.Errorf("unknown type of archive %s", strDataArchive))
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
			if _, err := f.WriteString("]"); err != nil {
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
				_, err = f.WriteString(",")

				if err != nil {
					t.Fatal(err)
				}
			}

			_, err = f.Write(b)

			if err != nil {
				t.Fatal(err)
			}

			emptyTrace = false
		}

		client = setupTracer(client, setupTracerOptions(bodies, callbackFunc)...)
	}

	c, err := NewConnection(client)

	if err != nil {
		t.Fatal(err)
	}

	err = c.Open(cascadeURL, WithAuth(aURL, Auth{
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

	var readingFunc func(sensor struct {
		*CounterHouseDto
		error
	}, archive archive.DataArchive, start, end time.Time)

	var nonSeparatedInputNumFunc = func(sensor struct {
		*CounterHouseDto
		error
	}, archive archive.DataArchive, start, end time.Time) {
		if sensor.error != nil {
			t.Fatal(sensor.error)
		}

		t.Logf("check readings for sensor id = %d (%s)", sensor.ID, sensor.Name)

		archiveBytes, err := c.Readings(sensor.ID, archiveType, start, end)

		if err != nil {
			t.Fatal(err)
		}

		var itemCount int
		readings := ParseCounterHouseReadingDto(context.Background(), archiveBytes)

		for item := range readings {
			if item.error != nil {
				t.Fatal(item.error)
			}

			itemCount++
		}

		t.Logf("\ttotal items: %d", itemCount)
	}

	var SeparatedInputNumFunc = func(sensor struct {
		*CounterHouseDto
		error
	}, archive archive.DataArchive, start, end time.Time) {
		if sensor.error != nil {
			t.Fatal(sensor.error)
		}

		t.Logf("check readings for sensor id = %d (%s)", sensor.ID, sensor.Name)

		var inputNum byte

		for _, entry := range sensor.Inputs {
			inputNum = byte(entry.Number)

			t.Logf("\tinput %d", inputNum)

			archiveBytes, err := c.Readings(sensor.ID, archiveType, start, end, inputNum)

			if err != nil {
				t.Fatal(err)
			}

			var itemCount int
			readings := ParseCounterHouseReadingDto(context.Background(), archiveBytes)

			for item := range readings {
				if item.error != nil {
					t.Fatal(item.error)
				}

				itemCount++
			}

			t.Logf("\t\ttotal items: %d", itemCount)
		}
	}

	readingFunc = nonSeparatedInputNumFunc

	if separatedInputs {
		readingFunc = SeparatedInputNumFunc
	}

	var sensorCount int

	sensors := ParseCounterHouseDto(context.Background(), sensorBytes)

	for sensor := range sensors {
		if sensor.error != nil {
			t.Fatal(sensor.error)
		}

		if limit > 0 && sensorCount > int(limit) {
			continue
		}

		readingFunc(sensor, archiveType, start, end)

		sensorCount++
	}

	sensorCount--
	t.Logf("TestConnection: total %s", english.Plural(sensorCount, "sensor", ""))
}

func TestConnection2(t *testing.T) {
	if strings.TrimSpace(cascadeURL) == "" {
		t.Fatal(errors.New("no API URL"))
	}

	aURL, err := url.Parse(authURL)

	if err != nil {
		t.Fatal(err)
	}

	start, err := time.Parse(layoutQuery, strStart)

	if err != nil {
		t.Fatal(err)
	}

	end, err := time.Parse(layoutQuery, strEnd)

	if err != nil {
		t.Fatal(err)
	}

	archiveType := archive.Parse(strDataArchive)

	if archiveType == archive.UnknownArchive {
		t.Fatal(fmt.Errorf("unknown type of archive %s", strDataArchive))
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
			if _, err := f.WriteString("]"); err != nil {
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
				_, err = f.WriteString(",")

				if err != nil {
					t.Fatal(err)
				}
			}

			_, err = f.Write(b)

			if err != nil {
				t.Fatal(err)
			}

			emptyTrace = false
		}

		client = setupTracer(client, setupTracerOptions(bodies, callbackFunc)...)
	}

	c, err := NewConnection(client)

	if err != nil {
		t.Fatal(err)
	}

	err = c.Open(cascadeURL, WithAuth(aURL, Auth{
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

	var readingFunc func(sensor struct {
		*CounterHouseDto
		error
	}, archive archive.DataArchive, start, end time.Time)

	var nonSeparatedInputNumFunc = func(sensor struct {
		*CounterHouseDto
		error
	}, archive archive.DataArchive, start, end time.Time) {
		if sensor.error != nil {
			t.Fatal(sensor.error)
		}

		t.Logf("check changed readings for sensor id = %d (%s)", sensor.ID, sensor.Name)

		archiveBytes, err := c.ChangedReadings(sensor.ID, archiveType, start, end)

		if err != nil {
			t.Fatal(err)
		}

		var itemCount int
		readings := ParseCounterHouseReadingDto(context.Background(), archiveBytes)

		for item := range readings {
			if item.error != nil {
				t.Fatal(item.error)
			}

			itemCount++
		}

		t.Logf("\ttotal items: %d", itemCount)
	}

	var SeparatedInputNumFunc = func(sensor struct {
		*CounterHouseDto
		error
	}, archive archive.DataArchive, start, end time.Time) {
		if sensor.error != nil {
			t.Fatal(sensor.error)
		}

		t.Logf("check changed readings for sensor id = %d (%s)", sensor.ID, sensor.Name)

		var inputNum byte

		for _, entry := range sensor.Inputs {
			inputNum = byte(entry.Number)

			t.Logf("\tinput %d", inputNum)

			archiveBytes, err := c.ChangedReadings(sensor.ID, archiveType, start, end, inputNum)

			if err != nil {
				t.Fatal(err)
			}

			var itemCount int
			readings := ParseCounterHouseReadingDto(context.Background(), archiveBytes)

			for item := range readings {
				if item.error != nil {
					t.Fatal(item.error)
				}

				itemCount++
			}

			t.Logf("\t\ttotal items: %d", itemCount)
		}
	}

	readingFunc = nonSeparatedInputNumFunc

	if separatedInputs {
		readingFunc = SeparatedInputNumFunc
	}

	var sensorCount int

	sensors := ParseCounterHouseDto(context.Background(), sensorBytes)

	for sensor := range sensors {
		if sensor.error != nil {
			t.Fatal(sensor.error)
		}

		if limit > 0 && sensorCount > int(limit) {
			continue
		}

		readingFunc(sensor, archiveType, start, end)

		sensorCount++
	}

	sensorCount--
	t.Logf("TestConnection2: total %s", english.Plural(sensorCount, "sensor", ""))
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

	return client
}

func setupTracer(client *http.Client, options ...httptracer.Option) *http.Client {
	return httptracer.Trace(client, options...)
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
