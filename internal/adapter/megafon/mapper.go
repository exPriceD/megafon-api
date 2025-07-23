package megafon

import (
	"fmt"
	dto "megafon-buisness-reports/internal/infrastructure/megafon/response"
	"time"

	"megafon-buisness-reports/internal/domain/entities"
)

const layout = time.RFC3339

func ToEntity(dto dto.CallDTO) (entities.Call, error) {
	start, err := time.Parse(layout, dto.StartRaw)
	if err != nil {
		return entities.Call{}, fmt.Errorf("parse start: %w", err)
	}

	return entities.Call{
		UID:          dto.UID,
		Type:         entities.CallType(dto.Type),
		Status:       entities.CallStatus(dto.Status),
		Client:       dto.Client,
		Diversion:    dto.Diversion,
		TelnumName:   dto.TelnumName,
		Destination:  dto.Destination,
		User:         dto.User,
		UserName:     dto.UserName,
		GroupName:    dto.GroupName,
		Start:        start,
		Wait:         dto.Wait,
		Duration:     dto.Duration,
		Record:       dto.Record,
		Rating:       dto.Rating,
		Note:         dto.Note,
		MissedStatus: entities.MissedStatusCode(dto.MissedStatus),
	}, nil
}
