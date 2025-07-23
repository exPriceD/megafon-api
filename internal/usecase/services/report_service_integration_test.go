//go:build integration
// +build integration

package services_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"

	"megafon-buisness-reports/internal/config"
	"megafon-buisness-reports/internal/domain/entities"
	infra "megafon-buisness-reports/internal/infrastructure/megafon"
	megsvc "megafon-buisness-reports/internal/infrastructure/megafon/services"
	"megafon-buisness-reports/internal/infrastructure/reporting"
	"megafon-buisness-reports/internal/interfaces"
	uc "megafon-buisness-reports/internal/usecase/services"
)

type noopLogger struct{}

func (noopLogger) Info(msg string, fields ...any)         {}
func (noopLogger) Warn(msg string, fields ...any)         {}
func (noopLogger) Error(msg string, fields ...any)        {}
func (noopLogger) Debug(msg string, fields ...any)        {}
func (l noopLogger) With(fields ...any) interfaces.Logger { return l }

func TestReportService_Generate_Live(t *testing.T) {
	_ = godotenv.Load("../../../.env")

	apiKey := os.Getenv("MEGAFON_API_KEY")
	if apiKey == "" {
		t.Skip("MEGAFON_API_KEY is not set; skipping live test")
	}

	baseURL := os.Getenv("MEGAFON_BASE_URL")
	if baseURL == "" {
		t.Skip("MEGAFON_BASE_URL is not set; skipping live test")
	}

	cfg := config.MegafonBuisnessConfig{BaseURL: baseURL, APIKey: apiKey}
	cl, err := infra.NewClient(cfg, noopLogger{})
	if err != nil {
		t.Fatalf("client: %v", err)
	}
	repo := megsvc.NewCallService(cl, noopLogger{})
	builder := reporting.NewExcelBuilder()
	svc := uc.NewReportService(repo, builder)

	now := time.Now().UTC()
	start := now.Add(-24 * 2 * time.Hour)
	end := now.Add(-24 * 1 * time.Hour)

	filter := entities.CallFilter{
		Start:         &start,
		End:           &end,
		Limit:         500,
		ProcessMissed: true,
	}

	buf, err := svc.GenerateCallReport(context.Background(), filter)
	if err != nil {
		t.Fatalf("GenerateCallReport: %v", err)
	}

	const out = "test_live_calls.xlsx"
	if err = os.WriteFile(out, buf.Bytes(), 0o644); err != nil {
		t.Fatalf("save file: %v", err)
	}
	t.Logf("отчёт сохранён: %s (%d bytes)", out, buf.Len())
}
