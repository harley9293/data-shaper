package converter

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
)

func ProtoToExcel(protoFilePath string, excelFilePath string) error {
	util := protoUtil{}
	err := util.Parse(protoFilePath)
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

	if err = f.SaveAs(excelFilePath + util.excelFileName); err != nil {
		return err
	}

	return nil
}

func ExcelToCfg(protoFilePath string, excelFilePath string, cfgOutputPath string) error {
	util := protoUtil{}
	err := util.Parse(protoFilePath)
	if err != nil {
		return err
	}

	err = util.LoadData(excelFilePath)
	if err != nil {
		return err
	}

	fileContents := ""
	for _, sheet := range util.sheetMap {
		fileContents += fmt.Sprintf("[%s]\n", sheet.codeName)
		fileContents += fmt.Sprintf("Size = %d\n\n", sheet.valueSize)

		fieldOrder := make([]string, len(sheet.fieldMap))
		for fieldName, field := range sheet.fieldMap {
			fieldOrder[field.index] = fieldName
		}
		for i := 0; i < sheet.valueSize; i++ {
			for _, fieldName := range fieldOrder {
				fileContents += fmt.Sprintf("%s_%d = %s\n", sheet.fieldMap[fieldName].codeName, i, sheet.fieldMap[fieldName].values[i])
			}
			fileContents += "\n"
		}
	}

	err = os.WriteFile(cfgOutputPath+util.codeName+".cfg", []byte(fileContents), 0644)
	if err != nil {
		return err
	}

	return nil
}
