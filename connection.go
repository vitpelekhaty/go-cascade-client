package cascade

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/vitpelekhaty/go-cascade-client/v2/archive"
)

// IConnection интерфейс соединения с API Каскада
type IConnection interface {
	// Open открывает соединение с API Каскада
	Open(ctx context.Context, rawURL, username, passwd string, options ...OpenOption) error

	// Close закрывает соединение с API Каскада
	Close(ctx context.Context) error

	// Connected возвращает текущее состояние соединения с API Каскада
	Connected() bool

	// Gauges возвращает список доступных приборов учета с тепловыми вводами и каналами
	Gauges(ctx context.Context) ([]byte, error)

	// CurrentReadings возвращает текущие показания прибора учета за указанный период. Если указан номер теплового
	// ввода, то возвращаются показания по этому вводу прибора учета
	CurrentReadings(ctx context.Context, deviceID int64, archive archive.DataArchive, beginAt, endAt time.Time,
		inputNum ...byte) ([]byte, error)

	// AlteredReadings возвращает измененные показания прибора учета за указанный период. Если указан номер теплового
	// ввода, то возвращаются показания по этому вводу прибора учета
	AlteredReadings(ctx context.Context, deviceID int64, archive archive.DataArchive, beginCreateAt,
		endCreateAt time.Time, inputNum ...byte) ([]byte, error)
}

// NewConnection возвращает настроенное соединение с Каскадом
func NewConnection(options ...Option) (IConnection, error) {
	opts := &connOptions{}

	for _, option := range options {
		option(opts)
	}

	conn := &connection{
		client: opts.client,
	}

	if conn.client == nil {
		conn.client = &http.Client{}
	}

	return conn, nil
}

var _ IConnection = (*connection)(nil)

type connection struct {
	rawURL, authURL string
	secret          string
	client          *http.Client
	token           *token
}

// Open открывает соединение с API Каскада
func (conn *connection) Open(ctx context.Context, rawURL, username, passwd string, options ...OpenOption) error {
	_, err := url.Parse(rawURL)

	if err != nil {
		return err
	}

	conn.rawURL = rawURL
	conn.authURL = rawURL

	opts := &openOptions{}

	for _, option := range options {
		option(opts)
	}

	if opts.authURL != conn.authURL && opts.authURL != "" {
		conn.authURL = opts.authURL
	}

	conn.secret = secret(username, passwd)

	return conn.login(ctx, conn.authURL, conn.secret)
}

func (conn *connection) login(ctx context.Context, authURL string, secret string) error {
	if _, err := url.Parse(authURL); err != nil {
		return fmt.Errorf("POST %s: %v", authURL, err)
	}

	form := url.Values{}
	form.Add("grant_type", "client_credentials")

	req, err := http.NewRequestWithContext(ctx, "POST", authURL, strings.NewReader(form.Encode()))

	if err != nil {
		return fmt.Errorf("POST %s: %v", authURL, err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", secret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	resp, err := conn.client.Do(req)

	if err != nil {
		return fmt.Errorf("POST %s: %v", authURL, err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("POST %s: %s", authURL, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("POST %s: %v", authURL, err)
	}

	var t token

	err = json.Unmarshal(body, &t)

	if err != nil {
		return fmt.Errorf("POST %s: %v", authURL, err)
	}

	conn.token = &t

	return nil
}

// Close закрывает соединение с API Каскада
func (conn *connection) Close(_ context.Context) error {
	conn.token = nil
	conn.secret = ""

	return nil
}

// Connected возвращает признак установленного соединения
func (conn *connection) Connected() bool {
	return conn.token != nil
}

// methodGauges метод получения списка приборов учета
const methodGauges = "/api/cascade/counter-house"

// Gauges возвращает список доступных приборов учета с тепловыми вводами и каналами
func (conn *connection) Gauges(ctx context.Context) ([]byte, error) {
	if err := conn.checkConnection(); err != nil {
		return nil, fmt.Errorf("GET %s: %v", methodGauges, err)
	}

	methodURL, err := pathJoin(conn.rawURL, methodGauges)

	if err != nil {
		return nil, fmt.Errorf("GET %s: %v", methodGauges, err)
	}

	var headers = map[string]string{
		"Authorization": fmt.Sprintf("%s %s", conn.token.Type, conn.token.Value),
	}

	data, statusCode, err := conn.gauges(ctx, methodURL, headers)

	if err != nil {
		if statusCode == http.StatusUnauthorized {
			err = conn.login(ctx, conn.authURL, conn.secret)

			if err != nil {
				return nil, fmt.Errorf("GET %s: %v", methodGauges, err)
			}

			data, _, err = conn.gauges(ctx, methodURL, headers)

			if err != nil {
				return nil, fmt.Errorf("GET %s: %v", methodGauges, err)
			}
		} else {
			return nil, fmt.Errorf("GET %s: %v", methodGauges, err)
		}
	}

	return data, nil
}

func (conn *connection) gauges(ctx context.Context, rawURL string, headers map[string]string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", rawURL, nil)

	if err != nil {
		return nil, http.StatusOK, err
	}

	if len(headers) > 0 {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := conn.client.Do(req)

	if err != nil {
		return nil, -1, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, errors.New(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, resp.StatusCode, err
	}

	return data, resp.StatusCode, nil
}

// methodCurrentReadings метод чтения архива показаний прибора учета
const methodCurrentReadings = "/api/cascade/counter-house/reading"

// CurrentReadings возвращает текущие показания прибора учета за указанный период. Если указан номер теплового
// ввода, то возвращаются показания по этому вводу прибора учета
func (conn *connection) CurrentReadings(ctx context.Context, deviceID int64, archive archive.DataArchive, beginAt,
	endAt time.Time, inputNum ...byte) ([]byte, error) {
	if err := conn.checkConnection(); err != nil {
		return nil, fmt.Errorf("POST %s: %v", methodCurrentReadings, err)
	}

	methodURL, err := pathJoin(conn.rawURL, methodCurrentReadings)

	if err != nil {
		return nil, fmt.Errorf("POST %s: %v", methodCurrentReadings, err)
	}

	readingsRequest := &CurrentReadingsRequest{
		DeviceID: deviceID,
		Archive:  archive,
		BeginAt:  RequestTime(beginAt),
		EndAt:    RequestTime(endAt),
	}

	if len(inputNum) > 0 {
		readingsRequest.InputNum = inputNum[0]
	}

	reqData, err := json.Marshal(readingsRequest)

	if err != nil {
		return nil, err
	}

	var headers = map[string]string{
		"Authorization": fmt.Sprintf("%s %s", conn.token.Type, conn.token.Value),
		"Content-Type":  "application/json",
	}

	data, statusCode, err := conn.readings(ctx, methodURL, headers, reqData)

	if err != nil {
		if statusCode == http.StatusUnauthorized {
			err = conn.login(ctx, conn.authURL, conn.secret)

			if err != nil {
				return nil, fmt.Errorf("GET %s: %v", methodCurrentReadings, err)
			}

			data, statusCode, err = conn.readings(ctx, methodURL, headers, reqData)

			if err != nil {
				return nil, fmt.Errorf("GET %s: %v", methodCurrentReadings, err)
			}
		} else {
			return nil, fmt.Errorf("GET %s: %v", methodCurrentReadings, err)
		}
	}

	if statusCode != http.StatusOK {
		if len(data) > 0 {
			var m errorMessage

			err = json.Unmarshal(data, &m)

			if err != nil {
				return nil, fmt.Errorf("POST %s %d: %v", methodCurrentReadings, statusCode, err)
			}

			return nil, NewCascadeError(&m, "POST", methodCurrentReadings, statusCode)
		}

		return nil, fmt.Errorf("POST %s: %s", methodCurrentReadings, http.StatusText(statusCode))
	}

	return data, nil
}

func (conn *connection) readings(ctx context.Context, rawURL string, headers map[string]string,
	payload []byte) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", rawURL, bytes.NewReader(payload))

	if err != nil {
		return nil, http.StatusOK, err
	}

	if len(headers) > 0 {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := conn.client.Do(req)

	if err != nil {
		return nil, -1, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, resp.StatusCode, errors.New(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, resp.StatusCode, err
	}

	return data, resp.StatusCode, nil
}

// methodAlteredReadings метод чтения архива измененных показаний прибора учета за предыдущие даты опроса
const methodAlteredReadings = "/api/cascade/counter-house/reading/created"

// AlteredReadings возвращает измененные показания прибора учета за указанный период. Если указан номер теплового
// ввода, то возвращаются показания по этому вводу прибора учета
func (conn *connection) AlteredReadings(ctx context.Context, deviceID int64, archive archive.DataArchive,
	beginCreateAt, endCreateAt time.Time, inputNum ...byte) ([]byte, error) {
	if err := conn.checkConnection(); err != nil {
		return nil, fmt.Errorf("POST %s: %v", methodAlteredReadings, err)
	}

	methodURL, err := pathJoin(conn.rawURL, methodAlteredReadings)

	if err != nil {
		return nil, fmt.Errorf("POST %s: %v", methodAlteredReadings, err)
	}

	readingsRequest := &AlteredReadingsRequest{
		DeviceID:      deviceID,
		Archive:       archive,
		BeginCreateAt: RequestTime(beginCreateAt),
		EndCreateAt:   RequestTime(endCreateAt),
	}

	if len(inputNum) > 0 {
		readingsRequest.InputNum = inputNum[0]
	}

	reqData, err := json.Marshal(readingsRequest)

	if err != nil {
		return nil, err
	}

	var headers = map[string]string{
		"Authorization": fmt.Sprintf("%s %s", conn.token.Type, conn.token.Value),
		"Content-Type":  "application/json",
	}

	data, statusCode, err := conn.readings(ctx, methodURL, headers, reqData)

	if err != nil {
		if statusCode == http.StatusUnauthorized {
			err = conn.login(ctx, conn.authURL, conn.secret)

			if err != nil {
				return nil, fmt.Errorf("GET %s: %v", methodAlteredReadings, err)
			}

			data, statusCode, err = conn.readings(ctx, methodURL, headers, reqData)

			if err != nil {
				return nil, fmt.Errorf("GET %s: %v", methodAlteredReadings, err)
			}
		} else {
			return nil, fmt.Errorf("GET %s: %v", methodAlteredReadings, err)
		}
	}

	if statusCode != http.StatusOK {
		if len(data) > 0 {
			var m errorMessage

			err = json.Unmarshal(data, &m)

			if err != nil {
				return nil, fmt.Errorf("POST %s %d: %v", methodAlteredReadings, statusCode, err)
			}

			return nil, NewCascadeError(&m, "POST", methodAlteredReadings, statusCode)
		}

		return nil, fmt.Errorf("POST %s: %s", methodAlteredReadings, http.StatusText(statusCode))
	}

	return data, nil
}

func (conn *connection) checkConnection() error {
	if conn.token == nil {
		return errors.New("user not authorized")
	}

	if conn.client == nil {
		return errors.New("no HTTP client")
	}

	return nil
}
