package pbz

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func protoToExcel(protoFilePath string, excelFilePath string) error {
	util, err := parseProto(protoFilePath)
	if err != nil {
		return err
	}

	f := excelize.NewFile()
	for sheetName, sheet := range util.sheetMap {
		if _, err = f.NewSheet(sheetName); err != nil {
			return err
		}

		for fieldName, field := range sheet.fieldMap {
			if err = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", string(rune('A'+field.index)), 1), fieldName); err != nil {
				return err
			}
		}
	}

	_ = f.DeleteSheet("Sheet1")

	if err = f.SaveAs(excelFilePath + util.fileName); err != nil {
		return err
	}

	return nil
}
