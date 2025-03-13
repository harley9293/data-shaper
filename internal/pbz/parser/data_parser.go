package parser

import (
	"github.com/harley9293/data-shaper/internal/pbz/schema"
	"github.com/xuri/excelize/v2"
)

func Data(filePath string, schema *schema.Proto) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}

	for sheetIndex, sheet := range schema.SheetList {
		rows, err := f.GetRows(sheet.Name)
		if err != nil {
			continue
		}

		if len(rows) <= 1 {
			continue
		}

		schema.SheetList[sheetIndex].ValueSize = len(rows) - 1
		for excelIndex, cell := range rows[0] {
			for fieldIndex, field := range sheet.FieldList {
				if field.Name != cell {
					continue
				}

				for i := 1; i < len(rows); i++ {
					cellValue := rows[i][excelIndex]
					sheet.FieldList[fieldIndex].Values = append(sheet.FieldList[fieldIndex].Values, cellValue)
				}

				break
			}
		}
	}

	return nil
}
