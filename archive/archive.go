package archive

import (
	"fmt"
	"strconv"
	"strings"
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

// String возвращает строковое описание типа архива показаний
func (a DataArchive) String() string {
	switch a {
	case HourArchive:
		return "HourArchive"
	case DailyArchive:
		return "DailyArchive"
	default:
		return "UnknownArchive"
	}
}

func (a *DataArchive) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)

	var i int
	i, err = strconv.Atoi(s)

	if err != nil {
		*a = UnknownArchive
		return
	}

	switch i {
	case int(HourArchive):
		*a = HourArchive
	case int(DailyArchive):
		*a = DailyArchive
	default:
		*a = UnknownArchive
		err = fmt.Errorf("unknown archive type %d", i)
	}

	return
}

// Parse преобразование строки в значение DataArchive
func Parse(archive string) DataArchive {
	switch archive {
	case "HourArchive":
		return HourArchive
	case "DailyArchive":
		return DailyArchive
	default:
		return UnknownArchive
	}
}
