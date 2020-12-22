package cascade

import (
	"fmt"
	"strings"
	"time"
)

// RequestTime описывает формат времени, принятый в запросах к АИСКУТЭ Каскад
type RequestTime time.Time

// ReadingTime описывает формат времени, принятый в показаниях АИСКУТЭ Каскад
type ReadingTime time.Time

const (
	requestTimeLayout = `02.01.2006 15:04:05`
	readingTimeLayout = `2006-01-02T15:04:05.999`
)

// UnmarshalJSON реализация интерфейса Unmarshaler для типа RequestTime
func (rt *RequestTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse(requestTimeLayout, s)

	*rt = RequestTime(t)

	return
}

// MarshalJSON реализация интерфейса Marshaler для типа RequestTime
func (rt RequestTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, rt.String())), nil
}

// String возвращает строковое представление типа RequestTime
func (rt *RequestTime) String() string {
	t := time.Time(*rt)
	return t.Format(requestTimeLayout)
}

// UnmarshalJSON реализация интерфейса Unmarshaler для типа ReadingTime
func (rt *ReadingTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse(readingTimeLayout, s)

	*rt = ReadingTime(t)

	return
}
