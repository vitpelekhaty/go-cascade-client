package cascade

import (
	"github.com/vitpelekhaty/go-cascade-client/v2/archive"
)

// CurrentReadingsRequest запрос архива показаний прибора учета
type CurrentReadingsRequest struct {
	// DeviceID идентификатор прибора учета
	DeviceID int64 `json:"deviceId"`

	// InputNum номер теплового ввода
	InputNum byte `json:"inputNum,omitempty"`

	// Archive тип архива показаний
	Archive archive.DataArchive `json:"archiveType"`

	// BeginAt время начала периода показаний прибора учета
	BeginAt RequestTime `json:"beginAt"`

	// EndAt время окончания периода показаний прибора учета
	EndAt RequestTime `json:"endAt"`
}

// AlteredReadingsRequest запрос архива измененных показаний прибора учета за предыдущее время
type AlteredReadingsRequest struct {
	// DeviceID идентификатор прибора учета
	DeviceID int64 `json:"deviceId"`

	// InputNum номер теплового ввода
	InputNum byte `json:"inputNum,omitempty"`

	// Archive тип архива показаний
	Archive archive.DataArchive `json:"archiveType"`

	// BeginAt время начала периода изменения показаний прибора учета
	BeginCreateAt RequestTime `json:"beginCreateAt"`

	// EndAt время окончания периода изменения показаний прибора учета
	EndCreateAt RequestTime `json:"endCreateAt"`
}
