package parsers

import (
	"fmt"
	"strings"
)

const (
	// heat тип ресурса - отопление
	heat = "Heat"

	// hotWater тип ресурса - горячая вода
	hotWater = "HotWater"

	// none тип ресурса - не указан
	none = "None"
)

// Resource тип ресурса
type Resource byte

const (
	// ResourceUnknown неизвестный тип ресурса
	ResourceUnknown Resource = iota

	// ResourceHeat тип ресурса - отопление
	ResourceHeat

	// ResourceHotWater тип ресурса - горячая вода
	ResourceHotWater

	// ResourceNone ресурс не указан (для общего потребления тепловой энергии по тепловому вводу прибора учета)
	ResourceNone
)

// UnmarshalJSON реализация интерфейса Unmarshaler для типа Resource
func (r *Resource) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)

	switch s {
	case heat:
		*r = ResourceHeat

	case hotWater:
		*r = ResourceHotWater

	case none:
		*r = ResourceNone

	default:
		*r = ResourceUnknown
		err = fmt.Errorf("unknown resource %s", s)
	}

	return
}
