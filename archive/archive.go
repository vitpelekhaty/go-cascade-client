package archive

import (
	"fmt"
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

const (
	hourArchive    = "Hour"
	dailyArchive   = "Day"
	unknownArchive = "Unknown"
)

// String возвращает строковое описание типа архива показаний
func (a DataArchive) String() string {
	switch a {
	case HourArchive:
		return hourArchive
	case DailyArchive:
		return dailyArchive
	default:
		return unknownArchive
	}
}

func (a *DataArchive) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)

	switch s {
	case hourArchive:
		*a = HourArchive
	case dailyArchive:
		*a = DailyArchive
	default:
		*a = UnknownArchive
		err = fmt.Errorf("unknown archive type %s", s)
	}

	return
}

// Parse преобразование строки в значение DataArchive
func Parse(archive string) DataArchive {
	switch archive {
	case hourArchive:
		return HourArchive
	case dailyArchive:
		return DailyArchive
	default:
		return UnknownArchive
	}
}
