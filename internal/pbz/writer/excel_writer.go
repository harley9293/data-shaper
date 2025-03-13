package writer

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/harley9293/data-shaper/internal/pbz/schema"
	"github.com/xuri/excelize/v2"
)

type ExcelWriter struct{}

func (w *ExcelWriter) Write(filePath string, proto *schema.Proto) error {
	f, err := excelize.OpenFile(filePath)
	if errors.Is(err, fs.ErrNotExist) {
		f = excelize.NewFile()
	}

	for _, sheet := range proto.SheetList {
		if !hasSheet(f, sheet.Name) {
			_, err = f.NewSheet(sheet.Name)
			if err != nil {
				return err
			}
		}

		nameToCell := getExistFieldMap(f, sheet.Name)
		availableIndex := getAvailableIndex(f, sheet.Name)
		for _, field := range sheet.FieldList {
			cellIndex := ""
			if existIndex, ok := nameToCell[field.Name]; !ok {
				cellIndex, err = setRepeatedCell(f, sheet.Name, availableIndex, field)
				if err != nil {
					return err
				}
				availableIndex++
			} else {
				cellIndex = fmt.Sprintf("%s%d", string(rune('A'+existIndex)), 1)
			}
			err = setCellComment(f, sheet.Name, cellIndex, field)
			if err != nil {
				return err
			}
		}
	}

	_ = f.DeleteSheet("Sheet1")
	if err = f.SaveAs(filePath); err != nil {
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
	if len(rows) > 0 {
		return len(rows[0])
	} else {
		return 0
	}
}

func setCellComment(f *excelize.File, sheetName string, cellIndex string, field schema.Field) error {
	_ = f.DeleteComment(sheetName, cellIndex)
	err := f.AddComment(sheetName, excelize.Comment{
		Cell:   cellIndex,
		Text:   field.MessageName + "\n" + field.MessageType.String() + "\n\n" + field.Note,
		Height: 400,
		Width:  400,
	})
	if err != nil {
		return err
	}

	return nil
}

func setRepeatedCell(f *excelize.File, sheetName string, availableIndex int, field schema.Field) (string, error) {
	styleID, err := getTextStyleID(f)
	if err != nil {
		return "", err
	}
	cellIndex := fmt.Sprintf("%s%d", string(rune('A'+availableIndex)), 1)
	if err = f.SetCellValue(sheetName, cellIndex, field.Name); err != nil {
		return "", err
	}
	if err = f.SetCellStyle(sheetName, fmt.Sprintf("%s2", string(rune('A'+availableIndex))), fmt.Sprintf("%s10", string(rune('A'+availableIndex))), styleID); err != nil {
		return "", err
	}

	return cellIndex, nil
}

func setConstCell(f *excelize.File, sheetName string, availableIndex int, field schema.Field) (string, error) {
	styleID, err := getTextStyleID(f)
	if err != nil {
		return "", err
	}
	cellIndex := fmt.Sprintf("A%d", availableIndex+1)
	if err = f.SetCellValue(sheetName, cellIndex, field.Name); err != nil {
		return "", err
	}
	if err = f.SetCellStyle(sheetName, fmt.Sprintf("B1"), fmt.Sprintf("B10"), styleID); err != nil {
		return "", err
	}

	return cellIndex, nil
}

func getTextStyleID(f *excelize.File) (int, error) {
	textFmt := "@"
	style := &excelize.Style{
		CustomNumFmt: &textFmt,
	}

	styleID, err := f.NewStyle(style)
	if err != nil {
		return 0, err
	}

	return styleID, nil
}
