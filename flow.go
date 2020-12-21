package cascade

import (
	"fmt"
	"strings"
)

// Flow тип подключения
type Flow byte

const (
	// FlowUnknown неизвестный тип подключения
	FlowUnknown Flow = iota
	// FlowDirect прямое подключение
	FlowDirect
	// FlowReverse обратное подключение
	FlowReverse
)

const (
	// inFlow тип подключения - прямое
	inFlow = "inFlow"
	// outFlow тип подключения - обратное
	outFlow = "outFlow"
)

// UnmarshalJSON реализация интерфейса Unmarshaler для типа Flow
func (f *Flow) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)

	switch s {
	case inFlow:
		*f = FlowDirect
	case outFlow:
		*f = FlowReverse
	default:
		*f = FlowUnknown
		err = fmt.Errorf("unknown flow %s", s)
	}

	return
}
