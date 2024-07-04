package writer

import (
	"errors"
	"fmt"
	"github.com/harley9293/data-shaper/internal/pbz/core"
	"github.com/xuri/excelize/v2"
	"io/fs"
)

type ExcelWriter struct{}

func (w *ExcelWriter) Write(filePath string, protoSchema *core.ProtoExcelSchema) error {
	f, err := excelize.OpenFile(filePath)
	if errors.Is(err, fs.ErrNotExist) {
		f = excelize.NewFile()
	}

	textFmt := "@"
	style := &excelize.Style{
		CustomNumFmt: &textFmt,
	}

	styleID, err := f.NewStyle(style)
	if err != nil {
		return err
	}

	for _, sheet := range protoSchema.SheetList {
		if !hasSheet(f, sheet.Name) {
			_, err = f.NewSheet(sheet.Name)
			if err != nil {
				return err
			}
		}

		nameToCell := getExistFieldMap(f, sheet.Name)
		availableIndex := getAvailableIndex(f, sheet.Name)
		for _, field := range sheet.FieldList {
			if _, ok := nameToCell[field.Name]; !ok {
				if err = f.SetCellValue(sheet.Name, fmt.Sprintf("%s%d", string(rune('A'+availableIndex)), 1), field.Name); err != nil {
					return err
				}
				if err = f.SetCellStyle(sheet.Name, fmt.Sprintf("%s%d", string(rune('A'+availableIndex)), 2), fmt.Sprintf("%s%d", string(rune('A'+availableIndex)), 10), styleID); err != nil {
					return err
				}
				availableIndex++
			}
		}
	}

	_ = f.DeleteSheet("Sheet1")
	if err = f.SaveAs(protoSchema.FilePath); err != nil {
		return err
	}

	return nil
}

func hasSheet(f *excelize.File, sheetName string) bool {
	for _, sheet := range f.GetSheetList() {
		if sheet == sheetName {
			return true
		}
	}

	return false
}

func getExistFieldMap(f *excelize.File, sheetName string) map[string]int {
	rows, _ := f.GetRows(sheetName)
	nameToCell := make(map[string]int)
	if len(rows) > 0 {
		for index, field := range rows[0] {
			nameToCell[field] = index
		}
	}
	return nameToCell
}

func getAvailableIndex(f *excelize.File, sheetName string) int {
	rows, _ := f.GetRows(sheetName)
	nameToCell := make(map[string]int)
	availableIndex := 0
	if len(rows) > 0 {
		for index, field := range rows[0] {
			nameToCell[field] = index
			availableIndex = index + 1
		}
	}
	return availableIndex
}
