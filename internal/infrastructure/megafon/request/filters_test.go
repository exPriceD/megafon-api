package request

import (
	"net/url"
	"testing"
	"time"
)

func TestHistoryParams_ToQuery(t *testing.T) {
	start := time.Date(2025, 7, 15, 12, 0, 0, 0, time.UTC)
	end := start.Add(2 * time.Hour)

	p := HistoryParams{
		Start:  &start,
		End:    &end,
		Type:   In,
		Limit:  100,
		Groups: []string{"sales", "support"},
	}

	got := p.ToQuery()

	exp := url.Values{
		"start":  []string{"20250715T120000Z"},
		"end":    []string{"20250715T140000Z"},
		"type":   []string{"in"},
		"limit":  []string{"100"},
		"groups": []string{"sales,support"},
	}
	if got.Encode() != exp.Encode() {
		t.Fatalf("want %q, got %q", exp.Encode(), got.Encode())
	}
}
