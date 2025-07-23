package services

import (
	"bytes"
	"context"
	"megafon-buisness-reports/internal/domain/entities"
	"megafon-buisness-reports/internal/domain/ports"
	"megafon-buisness-reports/internal/usecase/callstats"
)

type ReportBuilder interface {
	Build(sum callstats.Summary) (*bytes.Buffer, error)
}

// ReportService создаёт Excel-отчёты по звонкам.
type ReportService struct {
	callServiceRepo ports.CallRepository
	cityPhoneRepo   ports.CityPhoneRepository
	builder         ReportBuilder
}

func NewReportService(
	callServiceRepo ports.CallRepository,
	cityPhoneRepo ports.CityPhoneRepository,
	b ReportBuilder,
) *ReportService {
	return &ReportService{callServiceRepo: callServiceRepo, cityPhoneRepo: cityPhoneRepo, builder: b}
}

// GenerateCallReport загружает звонки, агрегирует их и строит XLSX.
func (s *ReportService) GenerateCallReport(ctx context.Context, filter entities.CallFilter, city string) (*bytes.Buffer, error) {
	cityPhone, err := s.cityPhoneRepo.GetByCity(ctx, city)
	if err != nil {
		return nil, err
	}

	var calls []entities.Call
	for _, phoneNumeber := range cityPhone.DiversionPhone {
		filter.Diversion = phoneNumeber
		history, err := s.callServiceRepo.History(ctx, filter)
		if err != nil {
			return nil, err
		}
		calls = append(calls, history...)
	}

	sum := callstats.Aggregate(calls)
	return s.builder.Build(sum)
}
