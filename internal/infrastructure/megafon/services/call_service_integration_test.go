package services_test

import (
	"context"
	"encoding/json"
	"megafon-buisness-reports/internal/domain/entities"
	"megafon-buisness-reports/internal/interfaces"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"

	"megafon-buisness-reports/internal/config"
	infra "megafon-buisness-reports/internal/infrastructure/megafon"
	"megafon-buisness-reports/internal/infrastructure/megafon/services"
)

type noopLogger struct{}

func (noopLogger) Info(msg string, fields ...any)         {}
func (noopLogger) Warn(msg string, fields ...any)         {}
func (noopLogger) Error(msg string, fields ...any)        {}
func (noopLogger) Debug(msg string, fields ...any)        {}
func (l noopLogger) With(fields ...any) interfaces.Logger { return l }

func TestCallServiceHistory_Live(t *testing.T) {
	_ = godotenv.Load("../../../../.env")

	apiKey := os.Getenv("MEGAFON_API_KEY")
	if apiKey == "" {
		t.Skip("MEGAFON_API_KEY is not set; skipping live test")
	}

	baseURL := os.Getenv("MEGAFON_BASE_URL")
	if baseURL == "" {
		t.Skip("MEGAFON_BASE_URL is not set; skipping live test")
	}

	cfg := config.MegafonBuisnessConfig{
		BaseURL: baseURL,
		APIKey:  apiKey,
	}
	cl, err := infra.NewClient(cfg, noopLogger{})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	svc := services.NewCallService(cl, noopLogger{})

	now := time.Now().UTC()
	start := now.Add(-24 * 1 * time.Hour)
	end := now.Add(-24 * 0 * time.Hour)

	params := entities.CallFilter{
		Start:         &start,
		End:           &end,
		Period:        "",
		Type:          entities.CallAll,
		Limit:         500,
		User:          "",
		Diversion:     "",
		Client:        "",
		Groups:        nil,
		FirstAnswered: false,
		ProcessMissed: true,
		MissedStatus:  nil,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	calls, err := svc.History(ctx, params)
	if err != nil {
		t.Fatalf("History error: %v", err)
	}
	t.Logf("получено %d звонков", len(calls))

	if data, err := json.MarshalIndent(calls, "", "  "); err == nil {
		t.Logf("Calls dump:\n%s", string(data))
	} else {
		t.Logf("не смог сериализовать calls: %v", err)
	}
}
