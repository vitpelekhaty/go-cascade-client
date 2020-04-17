package cascade

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// DataArchive архив показаний прибора учета
type DataArchive byte

const (
	// UnknownArchive неизвестный тип архива
	UnknownArchive DataArchive = 0
	// HourArchive часовой архив
	HourArchive DataArchive = 1
	// DailyArchive суточный архив
	DailyArchive DataArchive = 2
)

func (self *Connection) Readings(deviceID int64, archive DataArchive, beginAt, endAt time.Time) ([]byte, error) {
	if err := self.checkConnection(); err != nil {
		return nil, fmt.Errorf("POST %s: %v", Readings, err)
	}

	methodURL, err := URLJoin(self.baseURL, Readings)

	if err != nil {
		return nil, fmt.Errorf("POST %s: %v", Readings, err)
	}

	readingsRequest := &ReadingsRequest{
		DeviceID: deviceID,
		Archive:  archive,
		BeginAt:  RequestTime(beginAt),
		EndAt:    RequestTime(endAt),
	}

	reqData, err := json.Marshal(readingsRequest)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", methodURL, bytes.NewReader(reqData))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", self.TokenType(), self.AccessToken()))
	req.Header.Set("Content-Type", "application/json")

	resp, err := self.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("POST %s: %v", Readings, err)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("POST %s: %v", Readings, err)
	}

	if resp.StatusCode != http.StatusOK {
		var (
			errorMessage     string
			errorDescription string
		)

		if len(data) > 0 {
			var message Message

			err = json.Unmarshal(data, &message)

			if err != nil {
				return nil, fmt.Errorf("POST %s %d: %v", Readings, resp.StatusCode, err)
			}

			errorMessage = message.Text
			errorDescription = message.Description

			return nil, fmt.Errorf("POST %s %d: %s: %s", Readings, resp.StatusCode, errorMessage,
				errorDescription)

		}

		return nil, fmt.Errorf("POST %s: %s", Readings, resp.Status)
	}

	return data, nil
}

// RequestTime описывает формат времени, принятый в запросах к АИСКУТЭ Каскад
type RequestTime time.Time

// ReadingsRequest запрос архива показаний прибора учета
type ReadingsRequest struct {
	// DeviceID идентификатор прибора учета
	DeviceID int64 `json:"deviceId"`
	// Archive тип архива показаний
	Archive DataArchive `json:"archiveType"`
	// BeginAt время начала периода показаний прибора учета
	BeginAt RequestTime `json:"beginAt"`
	// EndAt время окончания периода показаний прибора учета
	EndAt RequestTime `json:"endAt"`
}

const requestTimeLayout = `02.01.2006 15:04:05`

func (self *RequestTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse(requestTimeLayout, s)

	*self = RequestTime(t)

	return
}

func (self RequestTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, self.String())), nil
}

func (self *RequestTime) String() string {
	t := time.Time(*self)
	return t.Format(requestTimeLayout)
}

func (self *DataArchive) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)

	var i int
	i, err = strconv.Atoi(s)

	if err != nil {
		*self = UnknownArchive
		return
	}

	switch i {
	case int(HourArchive):
		*self = HourArchive
	case int(DailyArchive):
		*self = DailyArchive
	default:
		*self = UnknownArchive
		err = fmt.Errorf("unknown archive type %d", i)
	}

	return
}

// CounterHouseReadingDto элемент архива показаний
type CounterHouseReadingDto struct {
	// Archive тип архива
	Archive int32 `json:"archiveType"`
	// ChannelID идентификатор канала/трубы
	ChannelID int64 `json:"channelId"`
	// DT момент показания
	DT RequestTime `json:"dt"`
	// ID идентификатор показания
	ID int64 `json:"id"`
	// IsBadRow признак "плохой" строки показания (признак нештатной ситуации, зафиксированной на приборе учета)
	IsBadRow bool `json:"isBadRow"`
	// M расход теплоносителя в тоннах
	M float32 `json:"m"`
	// P давление
	P float32 `json:"p"`
	// Q расход тепла в Гкал
	Q float32 `json:"q"`
	// T температура теплоносителя
	T float32 `json:"t"`
	// TCW температура холодной воды
	TCW float32 `json:"tcw"`
	// TI время штатной работы прибора учета
	TI int `json:"ti"`
	// V расход теплоносителя в м3
	V float32 `json:"v"`
}

// DataArchive возвращает тип архива показаний прибора учета
func (self *CounterHouseReadingDto) DataArchive() DataArchive {
	switch int(self.Archive) {
	case int(HourArchive):
		return HourArchive
	case int(DailyArchive):
		return DailyArchive
	default:
		return UnknownArchive
	}
}

// Message сообщение сервера об ошибке
type Message struct {
	// Text текст ошибки
	Text string `json:"message"`
	// Description описание ошибки
	Description string `json:"description"`
}
