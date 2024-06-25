package pbz

import (
	"errors"
	"fmt"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/xuri/excelize/v2"
	"strings"
	"unicode"
)

type fieldStruct struct {
	messageName string
	index       int

	values []string
}

type sheetStruct struct {
	sheetName   string
	messageName string
	fieldMap    map[string]*fieldStruct

	valueSize int
}

type excelStruct struct {
	fileName    string
	messageName string
	sheetMap    map[string]*sheetStruct
}

func parseProto(protoFilePath string) (*excelStruct, error) {
	util := &excelStruct{}
	util.sheetMap = make(map[string]*sheetStruct)

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

		comment := strings.TrimSpace(*md.GetSourceInfo().LeadingComments)

		if strings.Contains(comment, "@wrapper") {
			if comment == "@wrapper" {
				util.fileName = md.GetName() + ".xlsx"
			} else {
				util.fileName = comment[strings.Index(comment, "@wrapper")+9:] + ".xlsx"
			}
			util.messageName = md.GetName()

			for _, sheet := range md.GetFields() {
				sheetName := sheet.GetName()
				if sheet.GetSourceInfo().LeadingComments != nil {
					tmp := *sheet.GetSourceInfo().LeadingComments
					if strings.Contains(tmp, "@name") {
						sheetName = strings.TrimSpace(strings.Split(tmp, "@name")[1])
					}
				}

				util.sheetMap[sheetName] = &sheetStruct{sheetName: "#" + sheetName, fieldMap: make(map[string]*fieldStruct), messageName: sheet.GetName()}

				msgName := sheet.GetMessageType().GetFullyQualifiedName()
				mmd := fd.FindMessage(msgName)
				if mmd == nil {
					return nil, errors.New(fmt.Sprintf("%s message not found in proto file", msgName))
				}

				for index, field := range mmd.GetFields() {
					fieldName := field.GetName()
					if field.GetSourceInfo().LeadingComments != nil {
						tmp := *field.GetSourceInfo().LeadingComments
						if strings.Contains(tmp, "@name") {
							fieldName = strings.TrimSpace(strings.Split(tmp, "@name")[1])
						}
					}

					firstSpaceIndex := strings.IndexFunc(fieldName, unicode.IsSpace)
					if firstSpaceIndex != -1 {
						fieldName = fieldName[:firstSpaceIndex]
					}
					util.sheetMap[sheetName].fieldMap[fieldName] = &fieldStruct{index: index, messageName: field.GetName()}
				}
			}

			return util, nil
		}
	}

	return nil, errors.New("no wrapper found in proto file")
}

func (util *excelStruct) LoadData(excelFilePath string) error {
	f, err := excelize.OpenFile(excelFilePath)
	if err != nil {
		return err
	}
	defer func(f *excelize.File) {
		_ = f.Close()
	}(f)

	// for util.sheetMap
	for sheetName, _ := range util.sheetMap {
		if _, ok := util.sheetMap[sheetName]; !ok {
			continue
		}

		rows, err := f.GetRows(sheetName)
		if err != nil {
			return err
		}

		if len(rows) < 1 {
			return errors.New("no data found in excel file")
		}

		fieldIndexToName := make(map[int]string)
		for index, cell := range rows[0] {
			if _, ok := util.sheetMap[sheetName].fieldMap[cell]; !ok {
				continue
			}

			fieldIndexToName[index] = cell
		}

		for rowIndex, row := range rows {
			if rowIndex == 0 {
				continue
			}

			for index, cell := range row {
				if fieldName, ok := fieldIndexToName[index]; ok {
					util.sheetMap[sheetName].fieldMap[fieldName].values = append(util.sheetMap[sheetName].fieldMap[fieldName].values, cell)
				}
			}
			util.sheetMap[sheetName].valueSize++
		}
	}
	return nil
}
