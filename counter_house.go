package cascade

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// CounterHouse возвращает список приборов учета
func (self *Connection) CounterHouse() ([]byte, error) {
	if err := self.checkConnection(); err != nil {
		return nil, err
	}

	methodURL, err := URLJoin(self.baseURL, CounterHouse)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", methodURL, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", self.TokenType(), self.AccessToken()))

	resp, err := self.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %s", CounterHouse, resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return data, nil
}

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
func (self *CounterHouseChannelDto) Flow() Flow {
	switch self.Type {
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
func (self *CounterHouseChannelDto) Resource() Resource {
	switch self.ResourceType {
	case heat:
		return ResourceHeat
	case hotWater:
		return ResourceHotWater
	default:
		return ResourceUnknown
	}
}
