package cascade

import (
	"github.com/vitpelekhaty/go-cascade-client/archive"
)

// ReadingsRequest запрос архива показаний прибора учета
type ReadingsRequest struct {
	// DeviceID идентификатор прибора учета
	DeviceID int64 `json:"deviceId"`
	// Archive тип архива показаний
	Archive archive.DataArchive `json:"archiveType"`
	// BeginAt время начала периода показаний прибора учета
	BeginAt RequestTime `json:"beginAt"`
	// EndAt время окончания периода показаний прибора учета
	EndAt RequestTime `json:"endAt"`
}
