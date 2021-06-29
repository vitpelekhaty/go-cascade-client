package cascade

import (
	"strings"

	"github.com/guregu/null"

	"github.com/vitpelekhaty/go-cascade-client/archive"
)

// CounterHouseDto элемент списка приборов учета
type CounterHouseDto struct {
	// ID
	ID int64 `json:"id"`
	// Name наименование прибора учета
	Name string `json:"name"`
	// Model модель прибора учета
	Model string `json:"modelName"`
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
	// Resource тип ресурса
	Resource Resource `json:"resourceType"`
	// Flow тип подключения - подача или обратка
	Flow Flow `json:"type"`
}

// CounterHouseReadingDto элемент архива показаний
type CounterHouseReadingDto struct {
	// Archive тип архива
	Archive archive.DataArchive `json:"archiveType"`
	// ChannelID идентификатор канала/трубы
	ChannelID null.Int `json:"channelId"`
	// CreateAt момент чтения показания
	CreateAt ReadingTime `json:"createAt"`
	// DeviceID идентификатор прибора учета
	DeviceID null.Int `json:"deviceId"`
	// Input номер теплового ввода
	Input null.Int `json:"inputNum"`
	// DT момент показания
	DT ReadingTime `json:"dt"`
	// ID идентификатор показания
	ID null.Int `json:"id"`
	// IsBadRow признак "плохой" строки показания (признак нештатной ситуации, зафиксированной на приборе учета)
	IsBadRow bool `json:"isBadRow"`
	// M расход теплоносителя в тоннах
	M null.Float `json:"m"`
	// P давление
	P null.Float `json:"p"`
	// Q тепловая энергия по всему вводу, Гкал
	Q null.Float `json:"q"`
	// Q1 тепловая энергия по отоплению, Гкал
	Q1 null.Float `json:"q1"`
	// Q2 тепловая энергия по ГВС, Гкал
	Q2 null.Float `json:"q2"`
	// T температура теплоносителя
	T null.Float `json:"t"`
	// TCW температура холодной воды
	TCW null.Float `json:"tcw"`
	// TI время штатной работы прибора учета
	TI null.Float `json:"ti"`
	// V расход теплоносителя в м3
	V null.Float `json:"v"`
	// Empty признак "пустой" строки показания
	Empty bool `json:"isEmpty"`
}

// message сообщение сервера об ошибке
type message struct {
	// Message текст ошибки
	Message null.String `json:"message"`
	// Description описание ошибки
	Description null.String `json:"description"`
	// Error текст ошибки
	Error null.String `json:"error"`
	// Exception тип исключения
	Exception null.String `json:"exception"`
	// Status статус исключения
	Status null.String `json:"status"`
}

func (m *message) Err() *ErrCascadeCall {
	e := &ErrCascadeCall{
		exception:   m.Exception.ValueOrZero(),
		message:     m.Message.ValueOrZero(),
		description: m.Description.ValueOrZero(),
		status:      m.Status.ValueOrZero(),
	}

	if len(strings.TrimSpace(e.message)) == 0 {
		e.message = m.Error.ValueOrZero()
	}

	return e
}
