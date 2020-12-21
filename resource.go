package cascade

import (
	"fmt"
	"strings"
)

const (
	// heat тип ресурса - отопление
	heat = "Heat"
	// hotWater тип ресурса - горячая вода
	hotWater = "HotWater"
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
)

// UnmarshalJSON реализация интерфейса Unmarshaler для типа Resource
func (r *Resource) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)

	switch s {
	case heat:
		*r = ResourceHeat
	case hotWater:
		*r = ResourceHotWater
	default:
		*r = ResourceUnknown
		err = fmt.Errorf("unknown resource %s", s)
	}

	return
}
