package cascade

import (
	"fmt"
	"strings"
	"time"

	"github.com/vitpelekhaty/go-cascade-client/archive"
)

// RequestTime описывает формат времени, принятый в запросах к АИСКУТЭ Каскад
type RequestTime time.Time

// ReadingsRequest запрос архива показаний прибора учета
type ReadingsRequest struct {
	// DeviceID идентификатор прибора учета
	DeviceID int64 `json:"deviceId"`
	// Archive тип архива показаний
	Archive archive.DataArchive `json:"archiveType"`
	// BeginAt время начала периода показаний прибора учета
	BeginAt RequestTime `json:"beginAt"`
	// EndAt время окончания периода показаний прибора учета
	EndAt RequestTime `json:"endAt"`
}

const requestTimeLayout = `02.01.2006 15:04:05`

func (rt *RequestTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse(requestTimeLayout, s)

	*rt = RequestTime(t)

	return
}

func (rt RequestTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, rt.String())), nil
}

func (rt *RequestTime) String() string {
	t := time.Time(*rt)
	return t.Format(requestTimeLayout)
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
func (counter *CounterHouseReadingDto) DataArchive() archive.DataArchive {
	switch int(counter.Archive) {
	case int(archive.HourArchive):
		return archive.HourArchive
	case int(archive.DailyArchive):
		return archive.DailyArchive
	default:
		return archive.UnknownArchive
	}
}

// Message сообщение сервера об ошибке
type Message struct {
	// Text текст ошибки
	Text string `json:"message"`
	// Description описание ошибки
	Description string `json:"description"`
}
