package parsers

import (
	"strings"
	"time"
)

// ReadingTime описывает формат времени, принятый в показаниях АИСКУТЭ Каскад
type ReadingTime time.Time

const readingTimeLayout = `2006-01-02T15:04:05.999`

// UnmarshalJSON реализация интерфейса Unmarshaler для типа ReadingTime
func (rt *ReadingTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse(readingTimeLayout, s)

	*rt = ReadingTime(t)

	return
}
