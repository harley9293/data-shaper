package pbz

import (
	"errors"
	"fmt"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/xuri/excelize/v2"
	"io/fs"
	"strings"
)

type fieldStruct struct {
	fieldName   string
	messageName string

	values []string
}

type sheetStruct struct {
	sheetName   string
	messageName string
	fieldList   []*fieldStruct

	valueSize int
}

type excelStruct struct {
	dirPath     string
	fileName    string
	messageName string
	sheetList   []*sheetStruct
}

func parseProto(protoFilePath, excelDirPath string, loadData bool) (*excelStruct, error) {
	util := &excelStruct{}
	util.dirPath = excelDirPath

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
			util.fileName = getValueFromComments(md.GetSourceInfo().LeadingComments, "wrapper", md.GetName()) + ".xlsx"
			util.messageName = md.GetName()
			for _, sheet := range md.GetFields() {
				sheetName := getValueFromComments(sheet.GetSourceInfo().LeadingComments, "name", sheet.GetName())
				newSheet := &sheetStruct{sheetName: sheetName, messageName: sheet.GetName()}
				mmd := fd.FindMessage(sheet.GetMessageType().GetFullyQualifiedName())
				if mmd == nil {
					return nil, errors.New(fmt.Sprintf("%s message not found in proto file", sheet.GetMessageType().GetFullyQualifiedName()))
				}

				for _, field := range mmd.GetFields() {
					fieldName := getValueFromComments(field.GetSourceInfo().LeadingComments, "name", field.GetName())
					newSheet.fieldList = append(newSheet.fieldList, &fieldStruct{fieldName: fieldName, messageName: field.GetName()})
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

func (util *excelStruct) saveData() error {
	f, err := excelize.OpenFile(util.dirPath + util.fileName)
	if errors.Is(err, fs.ErrNotExist) {
		f = excelize.NewFile()
	}

	for _, sheet := range util.sheetList {
		exist := false
		for _, s := range f.GetSheetList() {
			if s == sheet.sheetName {
				exist = true
				break
			}
		}

		if !exist {
			_, err = f.NewSheet(sheet.sheetName)
			if err != nil {
				return err
			}
		}

		rows, _ := f.GetRows(sheet.sheetName)
		nameToCell := make(map[string]int)
		avilableIndex := 0
		if len(rows) > 0 {
			for index, field := range rows[0] {
				nameToCell[field] = index
				avilableIndex = index + 1
			}
		}

		for _, field := range sheet.fieldList {
			if _, ok := nameToCell[field.fieldName]; !ok {
				if err = f.SetCellValue(sheet.sheetName, fmt.Sprintf("%s%d", string(rune('A'+avilableIndex)), 1), field.fieldName); err != nil {
					return err
				}
				avilableIndex++
			}
		}
	}

	_ = f.DeleteSheet("Sheet1")
	if err = f.SaveAs(util.dirPath + util.fileName); err != nil {
		return err
	}

	return nil
}

func (util *excelStruct) loadData() {

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
