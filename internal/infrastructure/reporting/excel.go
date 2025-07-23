package reporting

import (
	"bytes"
	"fmt"

	"github.com/xuri/excelize/v2"
	"megafon-buisness-reports/internal/usecase/callstats"
)

type ExcelBuilder struct{}

func NewExcelBuilder() *ExcelBuilder { return &ExcelBuilder{} }

func (ExcelBuilder) Build(sum callstats.Summary) (*bytes.Buffer, error) {
	const sheet = "Calls"

	f := excelize.NewFile()
	idx, _ := f.NewSheet(sheet)
	_ = f.DeleteSheet("Sheet1")

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#D9E1F2"}, Pattern: 1},
		Border:    borders(1),
	})
	bodyStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "top"},
		Border:    borders(1),
	})
	topLeftStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "top"},
		Border:    borders(1),
	})

	_ = f.SetSheetRow(sheet, "A1", &[]any{"Категория", "Количество", "Номера"})
	_ = f.SetCellStyle(sheet, "A1", "C1", headerStyle)

	order := []callstats.Bucket{
		callstats.BucketNoCallBack,
		callstats.BucketWeCalledBackFail,
		callstats.BucketAll,
		callstats.BucketMissed,
		callstats.BucketWeCalledBackOK,
		callstats.BucketClientCalledBack,
	}

	next := 2
	for _, b := range order {
		stat := sum[b]
		if stat == nil {
			stat = &callstats.Stat{}
		}

		if len(stat.Numbers) == 0 {
			rowCell, _ := excelize.CoordinatesToCellName(1, next)
			_ = f.SetSheetRow(sheet, rowCell, &[]any{bucketName(b), stat.Count, ""})
			_ = f.SetCellStyle(sheet, rowCell, cellRef(3, next), bodyStyle)

			sepRow := next + 1
			sepCell, _ := excelize.CoordinatesToCellName(1, sepRow)
			_ = f.SetSheetRow(sheet, sepCell, &[]any{"", "", ""})
			_ = f.SetCellStyle(sheet, sepCell, cellRef(3, sepRow), bodyStyle)

			next += 2
			continue
		}

		start := next
		for i, num := range stat.Numbers {
			row := next + i
			cell, _ := excelize.CoordinatesToCellName(1, row)
			if i == 0 {
				_ = f.SetSheetRow(sheet, cell, &[]any{bucketName(b), stat.Count, num})
			} else {
				_ = f.SetSheetRow(sheet, cell, &[]any{"", "", num})
			}
			_ = f.SetCellStyle(sheet, cell, cellRef(3, row), bodyStyle)
		}
		end := next + len(stat.Numbers) - 1
		if end > start {
			_ = f.MergeCell(sheet, cellRef(1, start), cellRef(1, end))
			_ = f.MergeCell(sheet, cellRef(2, start), cellRef(2, end))
			_ = f.SetCellStyle(sheet, cellRef(1, start), cellRef(1, start), topLeftStyle)
			_ = f.SetCellStyle(sheet, cellRef(2, start), cellRef(2, start), topLeftStyle)
		}

		sepRow := end + 1
		sepCell, _ := excelize.CoordinatesToCellName(1, sepRow)
		_ = f.SetSheetRow(sheet, sepCell, &[]any{"", "", ""})
		_ = f.SetCellStyle(sheet, sepCell, cellRef(3, sepRow), bodyStyle)

		next = end + 2
	}

	_ = f.SetColWidth(sheet, "A", "A", 35)
	_ = f.SetColWidth(sheet, "B", "B", 13)
	_ = f.SetColWidth(sheet, "C", "C", 22)
	_ = f.SetPanes(sheet, &excelize.Panes{Freeze: true, YSplit: 1, TopLeftCell: "A2"})
	f.SetActiveSheet(idx)

	buf := &bytes.Buffer{}
	if err := f.Write(buf); err != nil {
		return nil, fmt.Errorf("write xlsx: %w", err)
	}
	return buf, nil
}

func borders(style int) []excelize.Border {
	return []excelize.Border{
		{Type: "left", Style: style},
		{Type: "right", Style: style},
		{Type: "top", Style: style},
		{Type: "bottom", Style: style},
	}
}

func bucketName(b callstats.Bucket) string {
	switch b {
	case callstats.BucketAll:
		return "Звонки"
	case callstats.BucketMissed:
		return "Пропущенные"
	case callstats.BucketClientCalledBack:
		return "Клиент перезвонил и дозвонился"
	case callstats.BucketWeCalledBackOK:
		return "Клиенту перезвонили и дозвонились"
	case callstats.BucketWeCalledBackFail:
		return "Клиенту перезвонили и не дозвонились"
	case callstats.BucketNoCallBack:
		return "Клиенту не перезвонили"
	default:
		return "Неизвестно"
	}
}

func cellRef(col, row int) string {
	ref, _ := excelize.CoordinatesToCellName(col, row)
	return ref
}
