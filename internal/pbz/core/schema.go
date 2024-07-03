package core

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io/fs"
)

type FieldSchema struct {
	Name        string
	MessageName string

	Values []string
}

type SheetSchema struct {
	Name        string
	MessageName string
	FieldList   []FieldSchema

	ValueSize int
}

type ProtoExcelSchema struct {
	FilePath    string
	MessageName string
	SheetList   []SheetSchema

	protoParser Parser
}

func NewProtoExcelSchema(excelDirPath string, protoParser Parser) *ProtoExcelSchema {
	return &ProtoExcelSchema{FilePath: excelDirPath, protoParser: protoParser}
}

func (schema *ProtoExcelSchema) ParseProto(filePath string) error {
	return schema.protoParser.Parse(filePath, schema)
}

func (schema *ProtoExcelSchema) SaveData() error {
	f, err := excelize.OpenFile(schema.FilePath)
	if errors.Is(err, fs.ErrNotExist) {
		f = excelize.NewFile()
	}

	for _, sheet := range schema.SheetList {
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
				availableIndex++
			}
		}
	}

	_ = f.DeleteSheet("Sheet1")
	if err = f.SaveAs(schema.FilePath); err != nil {
		return err
	}

	return nil
}

func (schema *ProtoExcelSchema) loadData() {

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
