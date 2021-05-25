package cascade

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/vitpelekhaty/go-cascade-client/archive"
)

// Option опция соединения с API Cascade
type Option func(c *Connection)

func WithAuth(authURL *url.URL, auth Auth) Option {
	return func(c *Connection) {
		c.authURL = authURL.String()
		c.secret = auth.Secret()
	}
}

var (
	// methodCounterHouse метод получения списка приборов учета
	methodCounterHouse = "/api/cascade/counter-house"
	// methodReadings метод чтения архива показаний прибора учета
	methodReadings = "/api/cascade/counter-house/reading"
)

// Connection соединение с Каскадом
type Connection struct {
	baseURL string
	client  *http.Client

	token *token

	authURL string
	secret  string

	OnError func(err error)
}

// NewConnection возвращает настроенное соединение с Каскадом
func NewConnection(client *http.Client) (*Connection, error) {
	if client == nil {
		return nil, errors.New("undefined HTTP client")
	}

	return &Connection{
		client: client,
	}, nil
}

// Open открывает соединение с API Cascade
func (c *Connection) Open(rawURL string, options ...Option) error {
	_, err := url.Parse(rawURL)

	if err != nil {
		return err
	}

	c.baseURL = rawURL

	for _, option := range options {
		option(c)
	}

	return c.login(c.authURL, c.secret)
}

func (c *Connection) Close() error {
	c.token = nil
	c.secret = ""

	return nil
}

// Connected возвращает признак установленного соединения
func (c *Connection) Connected() bool {
	return c.token != nil
}

// methodCounterHouse возвращает список приборов учета
func (c *Connection) CounterHouse() ([]byte, error) {
	if err := c.checkConnection(); err != nil {
		return nil, fmt.Errorf("GET %s: %v", methodCounterHouse, err)
	}

	methodURL, err := join(c.baseURL, methodCounterHouse)

	if err != nil {
		return nil, fmt.Errorf("GET %s: %v", methodCounterHouse, err)
	}

	req, err := http.NewRequest("GET", methodURL, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", c.token.Type, c.token.Value))

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("GET %s: %v", methodCounterHouse, err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.errorCallbackFunc(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %s", methodCounterHouse, resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("GET %s: %v", methodCounterHouse, err)
	}

	return data, nil
}

func (c *Connection) Readings(deviceID int64, archive archive.DataArchive, beginAt, endAt time.Time,
	inputNum ...byte) ([]byte, error) {
	if err := c.checkConnection(); err != nil {
		return nil, fmt.Errorf("POST %s: %v", methodReadings, err)
	}

	methodURL, err := join(c.baseURL, methodReadings)

	if err != nil {
		return nil, fmt.Errorf("POST %s: %v", methodReadings, err)
	}

	readingsRequest := &ReadingsRequest{
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

	req, err := http.NewRequest("POST", methodURL, bytes.NewReader(reqData))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", c.token.Type, c.token.Value))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("POST %s: %v", methodReadings, err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.errorCallbackFunc(err)
		}
	}()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("POST %s: %v", methodReadings, err)
	}

	if resp.StatusCode != http.StatusOK {
		if len(data) > 0 {
			message := &message{}

			err = json.Unmarshal(data, message)

			if err != nil {
				return nil, fmt.Errorf("POST %s %d: %v", methodReadings, resp.StatusCode, err)
			}

			ce := message.Err()

			ce.path = methodReadings
			ce.method = "POST"
			ce.statusCode = resp.StatusCode

			return nil, ce
		}

		return nil, fmt.Errorf("POST %s: %s", methodReadings, resp.Status)
	}

	return data, nil
}

func (c *Connection) checkConnection() error {
	if c.token == nil {
		return errors.New("user not authorized")
	}

	if c.client == nil {
		return errors.New("no HTTP client")
	}

	return nil
}

func (c *Connection) errorCallbackFunc(err error) {
	if err != nil && c.OnError != nil {
		c.OnError(err)
	}
}

// join возвращает полный URI метода API
func join(baseURL, method string) (string, error) {
	u, err := url.Parse(baseURL)

	if err != nil {
		return method, err
	}

	u.Path = path.Join(u.Path, method)

	return u.String(), nil
}
