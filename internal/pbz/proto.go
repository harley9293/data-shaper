package pbz

import (
	"errors"
	"fmt"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/xuri/excelize/v2"
	"io/fs"
	"strings"
)

type Field struct {
	name        string
	messageName string

	values []string
}

type Sheet struct {
	name        string
	messageName string
	fieldList   []Field

	valueSize int
}

type ProtoExcelSchema struct {
	filePath    string
	messageName string
	sheetList   []Sheet
}

func parseProto(protoFilePath, excelDirPath string, loadData bool) (*ProtoExcelSchema, error) {
	util := &ProtoExcelSchema{}
	util.filePath = excelDirPath

	parse := protoparse.Parser{
		IncludeSourceCodeInfo: true,
	}
	files, err := parse.ParseFiles(protoFilePath)
	if err != nil {
		return nil, err
	}

	fd := files[0]
	mds := fd.GetMessageTypes()

	for _, md := range mds {
		if md.GetSourceInfo() == nil || md.GetSourceInfo().LeadingComments == nil {
			continue
		}

		if hasKeyFromComments(md.GetSourceInfo().LeadingComments, "wrapper") {
			util.filePath += getValueFromComments(md.GetSourceInfo().LeadingComments, "wrapper", md.GetName()) + ".xlsx"
			util.messageName = md.GetName()
			for _, sheet := range md.GetFields() {
				newSheet := Sheet{name: getValueFromComments(sheet.GetSourceInfo().LeadingComments, "name", sheet.GetName()), messageName: sheet.GetName()}
				mmd := fd.FindMessage(sheet.GetMessageType().GetFullyQualifiedName())
				if mmd == nil {
					return nil, errors.New(fmt.Sprintf("%s message not found in proto file", sheet.GetMessageType().GetFullyQualifiedName()))
				}

				for _, field := range mmd.GetFields() {
					fieldName := getValueFromComments(field.GetSourceInfo().LeadingComments, "name", field.GetName())
					newSheet.fieldList = append(newSheet.fieldList, Field{name: fieldName, messageName: field.GetName()})
				}
				util.sheetList = append(util.sheetList, newSheet)
			}

			return util, nil
		}
	}

	if loadData {
		util.loadData()
	}

	return nil, errors.New("no wrapper found in proto file")
}

func (util *ProtoExcelSchema) saveData() error {
	f, err := excelize.OpenFile(util.filePath)
	if errors.Is(err, fs.ErrNotExist) {
		f = excelize.NewFile()
	}

	for _, sheet := range util.sheetList {
		if !hasSheet(f, sheet.name) {
			_, err = f.NewSheet(sheet.name)
			if err != nil {
				return err
			}
		}

		nameToCell := getExistFieldMap(f, sheet.name)
		availableIndex := getAvailableIndex(f, sheet.name)
		for _, field := range sheet.fieldList {
			if _, ok := nameToCell[field.name]; !ok {
				if err = f.SetCellValue(sheet.name, fmt.Sprintf("%s%d", string(rune('A'+availableIndex)), 1), field.name); err != nil {
					return err
				}
				availableIndex++
			}
		}
	}

	_ = f.DeleteSheet("Sheet1")
	if err = f.SaveAs(util.filePath); err != nil {
		return err
	}

	return nil
}

func (util *ProtoExcelSchema) loadData() {

}

func getValueFromComments(comments *string, key string, defaultValue string) (value string) {
	value = defaultValue
	if comments == nil {
		return
	}

	if strings.Contains(*comments, "@"+key) {
		value = strings.TrimSpace(strings.Split(*comments, "@"+key)[1])
	}

	if value == "" {
		value = defaultValue
	}

	return
}

func hasKeyFromComments(comments *string, key string) bool {
	if comments == nil {
		return false
	}

	return strings.Contains(*comments, "@"+key)
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
