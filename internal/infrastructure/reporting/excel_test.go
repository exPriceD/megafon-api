package reporting_test

import (
	"bytes"
	"testing"

	"github.com/xuri/excelize/v2"
	"megafon-buisness-reports/internal/infrastructure/reporting"
	"megafon-buisness-reports/internal/usecase/callstats"
)

func TestExcelBuilder_VerticalNumbers(t *testing.T) {
	sum := callstats.Summary{
		callstats.BucketMissed: &callstats.Stat{
			Count:   3,
			Numbers: []string{"111", "222", "333"},
		},
	}

	buf, err := reporting.NewExcelBuilder().Build(sum)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	f, err := excelize.OpenReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("open xlsx: %v", err)
	}
	defer f.Close()

	if val, _ := f.GetCellValue("Calls", "A2"); val != "Пропущенные" {
		t.Errorf("A2 want 'Пропущенные', got %q", val)
	}
	if val, _ := f.GetCellValue("Calls", "B2"); val != "3" {
		t.Errorf("B2 want '3', got %q", val)
	}
	if val, _ := f.GetCellValue("Calls", "C2"); val != "111" {
		t.Errorf("C2 want '111', got %q", val)
	}
	if val, _ := f.GetCellValue("Calls", "C3"); val != "222" {
		t.Errorf("C3 want '222', got %q", val)
	}
	if val, _ := f.GetCellValue("Calls", "C4"); val != "333" {
		t.Errorf("C4 want '333', got %q", val)
	}

	if val, _ := f.GetCellValue("Calls", "A3"); val != "Пропущенные" {
		t.Errorf("A3 expected merged value, got %q", val)
	}

	merges, err := f.GetMergeCells("Calls")
	if err != nil {
		t.Fatalf("GetMergeCells: %v", err)
	}
	hasMerge := func(ref string) bool {
		for _, m := range merges {
			if m.GetStartAxis()+":"+m.GetEndAxis() == ref {
				return true
			}
		}
		return false
	}
	if !hasMerge("A2:A4") || !hasMerge("B2:B4") {
		t.Errorf("ожидались объединения A2:A4 и B2:B4, получили %+v", merges)
	}
}
