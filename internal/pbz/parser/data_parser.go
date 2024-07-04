package parser

import (
	"github.com/harley9293/data-shaper/internal/pbz/core"
	"github.com/xuri/excelize/v2"
)

type DataParser struct{}

func (parser *DataParser) Parse(filePath string, protoSchema *core.ProtoExcelSchema) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}

	for sheetIndex, sheet := range protoSchema.SheetList {
		rows, err := f.GetRows(sheet.Name)
		if err != nil {
			continue
		}

		if sheet.Repeated {
			if len(rows) <= 1 {
				continue
			}

			protoSchema.SheetList[sheetIndex].ValueSize = len(rows) - 1
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
		} else {
			if len(rows) < 1 {
				continue
			}

			protoSchema.SheetList[sheetIndex].ValueSize = 1
			for _, row := range rows {
				if len(row) < 1 {
					continue
				}
				fieldName := row[0]
				fieldValue := ""
				if len(row) > 1 {
					fieldValue = row[1]
				}

				for fieldIndex, field := range sheet.FieldList {
					if field.Name != fieldName {
						continue
					}

					sheet.FieldList[fieldIndex].Values = append(sheet.FieldList[fieldIndex].Values, fieldValue)
					break
				}
			}
		}
	}

	return nil
}
