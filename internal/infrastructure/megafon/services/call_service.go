package services

import (
	"context"
	"fmt"
	adapter "megafon-buisness-reports/internal/adapter/megafon"
	"megafon-buisness-reports/internal/infrastructure/megafon"
	"megafon-buisness-reports/internal/infrastructure/megafon/request"
	"megafon-buisness-reports/internal/infrastructure/megafon/response"
	"net/http"

	"megafon-buisness-reports/internal/domain/entities"
	"megafon-buisness-reports/internal/domain/ports"
	"megafon-buisness-reports/internal/interfaces"
)

const historyEndpoint = "/crmapi/v1/history/json"

var _ ports.CallRepository = (*CallService)(nil)

type CallService struct {
	cl  *megafon.Client
	log interfaces.Logger
}

func NewCallService(cl *megafon.Client, lg interfaces.Logger) *CallService {
	return &CallService{cl: cl, log: lg}
}

// History реализует порт CallRepository.
func (s *CallService) History(ctx context.Context, f entities.CallFilter) ([]entities.Call, error) {
	params := request.HistoryParams{
		Start:         f.Start,
		End:           f.End,
		Period:        request.Period(f.Period),
		Type:          request.CallType(f.Type),
		Limit:         f.Limit,
		User:          f.User,
		Diversion:     f.Diversion,
		Client:        f.Client,
		Groups:        f.Groups,
		FirstAnswered: f.FirstAnswered,
		ProcessMissed: f.ProcessMissed,
		MissedStatus:  f.MissedStatus,
	}
	fmt.Println(params.ToQuery())
	var dto []response.CallDTO
	if err := s.cl.Do(ctx, http.MethodGet, historyEndpoint, params.ToQuery(), nil, &dto); err != nil {
		s.log.Error("history request failed", "err", err)
		return nil, err
	}

	calls := make([]entities.Call, 0, len(dto))
	for _, d := range dto {
		ent, err := adapter.ToEntity(d)
		if err != nil {
			s.log.Warn("bad call record", "err", err, "uid", d.UID)
			continue
		}
		calls = append(calls, ent)
	}

	s.log.Info("history fetched", "count", len(calls))
	return calls, nil
}
