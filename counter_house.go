package cascade

// CounterHouseDto элемент списка приборов учета
type CounterHouseDto struct {
	// ID
	ID int64 `json:"id"`
	// Name тип (модель) прибора учета
	Name string `json:"name"`
	// SN серийный номер прибора
	SN string `json:"serialNumber"`
	// Title наименование прибора учета в АИСКУТЭ
	Title string `json:"title"`
	// Inputs тепловые вводы
	Inputs []CounterHouseEntryDto `json:"inputs"`
}

// CounterHouseEntryDto элемент списка тепловых вводов на приборе учета
type CounterHouseEntryDto struct {
	// Number номер теплового ввода
	Number int32 `json:"number"`
	// Channels каналы/трубы
	Channels []CounterHouseChannelDto `json:"channels"`
}

// CounterHouseChannelDto элемент списка каналов/труб на приборе учета
type CounterHouseChannelDto struct {
	// ID идентификатор канала/трубы
	ID int64 `json:"id"`
	// Number номер канала/трубы
	Number int32 `json:"number"`
	// ResourceType тип ресурса
	ResourceType string `json:"resourceType"`
	// Type тип подключения - подача или обратка
	Type string `json:"type"`
}

// inFlow тип подключения - прямое
const inFlow = "inFlow"

// outFlow тип подключения - обратное
const outFlow = "outFlow"

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

// Flow возвращает тип подключения
func (channel *CounterHouseChannelDto) Flow() Flow {
	switch channel.Type {
	case inFlow:
		return FlowDirect
	case outFlow:
		return FlowReverse
	default:
		return FlowUnknown
	}
}

// heat тип ресурса - отопление
const heat = "Heat"

// hotWater тип ресурса - горячая вода
const hotWater = "HotWater"

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

// Resource возвращает тип ресурса
func (channel *CounterHouseChannelDto) Resource() Resource {
	switch channel.ResourceType {
	case heat:
		return ResourceHeat
	case hotWater:
		return ResourceHotWater
	default:
		return ResourceUnknown
	}
}
