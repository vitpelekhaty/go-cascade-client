//go:build integration

package cascade

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/guregu/null"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vitpelekhaty/go-cascade-client/v2/archive"
	"github.com/vitpelekhaty/httptracer"
)

var (
	username              string
	password              string
	authURL               string
	cascadeURL            string
	strInsecureSkipVerify string
	strBodies             string
	tracePath             string
	strTimeout            string
	strDataArchive        string
	casesPath             string
)

var envFile = flag.String("env", "", "параметры теста")

func TestConnection_Gauges(t *testing.T) {
	flag.Parse()

	err := godotenv.Load(*envFile)

	require.NoError(t, err, "ошибка загрузки параметров теста")

	cascadeURL = os.Getenv("CASCADE_URL")
	authURL = os.Getenv("CASCADE_AUTH_URL")
	username = os.Getenv("CASCADE_USERNAME")
	password = os.Getenv("CASCADE_USER_PASSWD")
	strTimeout = os.Getenv("CASCADE_TEST_PARAM_TIMEOUT")
	strInsecureSkipVerify = os.Getenv("CASCADE_TEST_PARAM_INSECURE_SKIP_VERIFY")
	tracePath = os.Getenv("CASCADE_TRACE_PATH")
	strBodies = os.Getenv("CASCADE_TRACE_PARAM_BODIES")

	require.NotEmpty(t, cascadeURL, "не указан адрес API")

	_, err = url.Parse(cascadeURL)

	require.NoError(t, err, "некорректный URL API Каскада")

	if authURL != "" {
		_, err = url.Parse(authURL)

		require.NoError(t, err, "некорректный URL авторизации в API Каскад")
	}

	var timeout = time.Second * 30

	if strTimeout != "" {
		timeout, err = time.ParseDuration(strTimeout)

		require.NoError(t, err, "некорректный формат таймаута")
	}

	var insecureSkipVerify = false

	if strInsecureSkipVerify != "" {
		insecureSkipVerify, err = strconv.ParseBool(strInsecureSkipVerify)

		require.NoError(t, err, "insecure-skip-verify")
	}

	var bodies = true

	if strBodies != "" {
		bodies, err = strconv.ParseBool(strBodies)

		require.NoError(t, err, "bodies")
	}

	client := setupHTTPClient(timeout, insecureSkipVerify)

	if strings.TrimSpace(tracePath) != "" {
		f, err := os.Create(tracePath)

		require.NoError(t, err)

		defer func() {
			_, _ = f.WriteString("]")
			_ = f.Close()
		}()

		_, _ = f.WriteString("[")

		emptyTrace := true

		callbackFunc := func(entry *httptracer.Entry) {
			if b, err := json.Marshal(entry); err == nil {
				if !emptyTrace {
					_, _ = f.WriteString(",")
				}

				_, _ = f.Write(b)

				emptyTrace = false
			}
		}

		client = setupTracer(client, setupTracerOptions(bodies, callbackFunc)...)
	}

	c, err := NewConnection(WithHTTPClient(client))

	require.NoError(t, err)

	ctx := context.TODO()

	err = c.Open(ctx, cascadeURL, username, password, WithAuthURL(authURL))

	require.NoError(t, err)

	defer func() {
		err := c.Close(ctx)
		assert.NoError(t, err)
	}()

	gb, err := c.Gauges(ctx)

	require.NoError(t, err)
	assert.NotEmpty(t, gb)
}

type gaugeTestCase struct {
	// GaugeID идентификатор прибора учета
	GaugeID int `json:"deviceID"`

	// Input тепловой ввод
	Input null.Int `json:"input"`

	// From начало периода запрашиваемых показаний
	From time.Time `json:"from"`

	// To окончание периода запрашиваемых показаний
	To null.Time `json:"to"`
}

func loadTestCases(path string) ([]gaugeTestCase, error) {
	b, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var cases = make([]gaugeTestCase, 0, 10)

	err = json.Unmarshal(b, &cases)

	if err != nil {
		return nil, err
	}

	return cases, nil
}

func TestConnection_CurrentReadings(t *testing.T) {
	flag.Parse()

	err := godotenv.Load(*envFile)

	require.NoError(t, err, "ошибка загрузки параметров теста")

	cascadeURL = os.Getenv("CASCADE_URL")
	authURL = os.Getenv("CASCADE_AUTH_URL")
	username = os.Getenv("CASCADE_USERNAME")
	password = os.Getenv("CASCADE_USER_PASSWD")
	strDataArchive = os.Getenv("CASCADE_TEST_ARCHIVE")
	strTimeout = os.Getenv("CASCADE_TEST_PARAM_TIMEOUT")
	strInsecureSkipVerify = os.Getenv("CASCADE_TEST_PARAM_INSECURE_SKIP_VERIFY")
	tracePath = os.Getenv("CASCADE_TRACE_PATH")
	strBodies = os.Getenv("CASCADE_TRACE_PARAM_BODIES")
	casesPath = os.Getenv("CASCADE_TEST_CASES_PATH")

	require.NotEmpty(t, casesPath, "не указан файл тестовых случаев")

	cases, err := loadTestCases(casesPath)

	require.NoError(t, err)
	require.NotEmpty(t, cases, "нет тестовых случаев")

	require.NotEmpty(t, cascadeURL, "не указан адрес API")

	_, err = url.Parse(cascadeURL)

	require.NoError(t, err, "некорректный URL API Каскада")

	if authURL != "" {
		_, err = url.Parse(authURL)

		require.NoError(t, err, "некорректный URL авторизации в API Каскад")
	}

	archiveType := archive.Parse(strDataArchive)

	require.NotEqual(t, archive.UnknownArchive, archiveType, "неизвестный тип архива показаний")

	var timeout = time.Second * 30

	if strTimeout != "" {
		timeout, err = time.ParseDuration(strTimeout)

		require.NoError(t, err, "некорректный формат таймаута")
	}

	var insecureSkipVerify = false

	if strInsecureSkipVerify != "" {
		insecureSkipVerify, err = strconv.ParseBool(strInsecureSkipVerify)

		require.NoError(t, err, "insecure-skip-verify")
	}

	var bodies = true

	if strBodies != "" {
		bodies, err = strconv.ParseBool(strBodies)

		require.NoError(t, err, "bodies")
	}

	client := setupHTTPClient(timeout, insecureSkipVerify)

	if strings.TrimSpace(tracePath) != "" {
		f, err := os.Create(tracePath)

		require.NoError(t, err)

		defer func() {
			_, _ = f.WriteString("]")
			_ = f.Close()
		}()

		_, _ = f.WriteString("[")

		emptyTrace := true

		callbackFunc := func(entry *httptracer.Entry) {
			if b, err := json.Marshal(entry); err == nil {
				if !emptyTrace {
					_, _ = f.WriteString(",")
				}

				_, _ = f.Write(b)

				emptyTrace = false
			}
		}

		client = setupTracer(client, setupTracerOptions(bodies, callbackFunc)...)
	}

	c, err := NewConnection(WithHTTPClient(client))

	require.NoError(t, err)

	ctx := context.TODO()

	err = c.Open(ctx, cascadeURL, username, password, WithAuthURL(authURL))

	require.NoError(t, err)

	defer func() {
		err := c.Close(ctx)
		assert.NoError(t, err)
	}()

	for _, test := range cases {
		var to = time.Now()

		if test.To.Valid {
			to = test.To.ValueOrZero()
		}

		var inputs = make([]byte, 0, 1)

		if test.Input.Valid {
			inputs = append(inputs, byte(test.Input.ValueOrZero()))
		}

		b, err := c.CurrentReadings(ctx, int64(test.GaugeID), archiveType, test.From, to, inputs...)

		assert.NoError(t, err)

		if err == nil {
			assert.NotEmpty(t, b)
		}
	}
}

func TestConnection_AlteredReadings(t *testing.T) {
	flag.Parse()

	err := godotenv.Load(*envFile)

	require.NoError(t, err, "ошибка загрузки параметров теста")

	cascadeURL = os.Getenv("CASCADE_URL")
	authURL = os.Getenv("CASCADE_AUTH_URL")
	username = os.Getenv("CASCADE_USERNAME")
	password = os.Getenv("CASCADE_USER_PASSWD")
	strDataArchive = os.Getenv("CASCADE_TEST_ARCHIVE")
	strTimeout = os.Getenv("CASCADE_TEST_PARAM_TIMEOUT")
	strInsecureSkipVerify = os.Getenv("CASCADE_TEST_PARAM_INSECURE_SKIP_VERIFY")
	tracePath = os.Getenv("CASCADE_TRACE_PATH")
	strBodies = os.Getenv("CASCADE_TRACE_PARAM_BODIES")
	casesPath = os.Getenv("CASCADE_TEST_CASES_PATH")

	require.NotEmpty(t, casesPath, "не указан файл тестовых случаев")

	cases, err := loadTestCases(casesPath)

	require.NoError(t, err)
	require.NotEmpty(t, cases, "нет тестовых случаев")

	require.NotEmpty(t, cascadeURL, "не указан адрес API")

	_, err = url.Parse(cascadeURL)

	require.NoError(t, err, "некорректный URL API Каскада")

	if authURL != "" {
		_, err = url.Parse(authURL)

		require.NoError(t, err, "некорректный URL авторизации в API Каскад")
	}

	archiveType := archive.Parse(strDataArchive)

	require.NotEqual(t, archive.UnknownArchive, archiveType, "неизвестный тип архива показаний")

	var timeout = time.Second * 30

	if strTimeout != "" {
		timeout, err = time.ParseDuration(strTimeout)

		require.NoError(t, err, "некорректный формат таймаута")
	}

	var insecureSkipVerify = false

	if strInsecureSkipVerify != "" {
		insecureSkipVerify, err = strconv.ParseBool(strInsecureSkipVerify)

		require.NoError(t, err, "insecure-skip-verify")
	}

	var bodies = true

	if strBodies != "" {
		bodies, err = strconv.ParseBool(strBodies)

		require.NoError(t, err, "bodies")
	}

	client := setupHTTPClient(timeout, insecureSkipVerify)

	if strings.TrimSpace(tracePath) != "" {
		f, err := os.Create(tracePath)

		require.NoError(t, err)

		defer func() {
			_, _ = f.WriteString("]")
			_ = f.Close()
		}()

		_, _ = f.WriteString("[")

		emptyTrace := true

		callbackFunc := func(entry *httptracer.Entry) {
			if b, err := json.Marshal(entry); err == nil {
				if !emptyTrace {
					_, _ = f.WriteString(",")
				}

				_, _ = f.Write(b)

				emptyTrace = false
			}
		}

		client = setupTracer(client, setupTracerOptions(bodies, callbackFunc)...)
	}

	c, err := NewConnection(WithHTTPClient(client))

	require.NoError(t, err)

	ctx := context.TODO()

	err = c.Open(ctx, cascadeURL, username, password, WithAuthURL(authURL))

	require.NoError(t, err)

	defer func() {
		err := c.Close(ctx)
		assert.NoError(t, err)
	}()

	for _, test := range cases {
		var to = time.Now()

		if test.To.Valid {
			to = test.To.ValueOrZero()
		}

		var inputs = make([]byte, 0, 1)

		if test.Input.Valid {
			inputs = append(inputs, byte(test.Input.ValueOrZero()))
		}

		b, err := c.AlteredReadings(ctx, int64(test.GaugeID), archiveType, test.From, to, inputs...)

		assert.NoError(t, err)

		if err == nil {
			assert.NotEmpty(t, b)
		}
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
