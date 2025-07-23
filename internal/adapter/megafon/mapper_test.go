package megafon

import (
	"testing"
	"time"

	dto "megafon-buisness-reports/internal/infrastructure/megafon/response"
)

func TestToEntity(t *testing.T) {
	const layout = time.RFC3339

	goodDTO := dto.CallDTO{
		UID:      "abc",
		Type:     "in",
		Status:   "ok",
		StartRaw: "2025-07-16T10:04:05Z",
		Duration: 30,
		Wait:     2,
	}
	wantStart, _ := time.Parse(layout, goodDTO.StartRaw)

	t.Run("success", func(t *testing.T) {
		got, err := ToEntity(goodDTO)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.UID != goodDTO.UID || !got.Start.Equal(wantStart) {
			t.Errorf("mismatch: %+v", got)
		}
	})

	t.Run("bad date", func(t *testing.T) {
		bad := goodDTO
		bad.StartRaw = "not-a-date"
		if _, err := ToEntity(bad); err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
