package reports

import "github.com/google/uuid"

type SqsMessage struct {
	UserId   uuid.UUID `json:"userId"`
	ReportId uuid.UUID `json:"reportId"`
}
