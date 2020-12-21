package cascade

import (
	"fmt"
	"strings"
	"time"
)

// RequestTime описывает формат времени, принятый в запросах к АИСКУТЭ Каскад
type RequestTime time.Time

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
