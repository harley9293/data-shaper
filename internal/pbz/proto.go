package pbz

import (
	"errors"
	"fmt"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/xuri/excelize/v2"
	"strings"
	"unicode"
)

type FieldUtil struct {
	index    int
	values   []string
	codeName string
}

type sheetUtil struct {
	fieldMap  map[string]*FieldUtil
	valueSize int
	codeName  string
}

type protoUtil struct {
	excelFileName string
	codeName      string
	sheetMap      map[string]*sheetUtil
}

func (util *protoUtil) Parse(protoFilePath string) error {
	util.sheetMap = make(map[string]*sheetUtil)

	parse := protoparse.Parser{
		IncludeSourceCodeInfo: true,
	}
	files, err := parse.ParseFiles(protoFilePath)
	if err != nil {
		return err
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
				util.excelFileName = md.GetName() + ".xlsx"
			} else {
				util.excelFileName = comment[strings.Index(comment, "@wrapper")+9:] + ".xlsx"
			}
			util.codeName = md.GetName()

			for _, sheet := range md.GetFields() {
				sheetComment := strings.TrimSpace(*sheet.GetSourceInfo().LeadingComments)
				if !strings.Contains(sheetComment, "@name") {
					return errors.New("no sheet name found in proto file")
				}

				sheetName := "#" + sheetComment[strings.Index(sheetComment, "@name")+6:]
				util.sheetMap[sheetName] = &sheetUtil{fieldMap: make(map[string]*FieldUtil), codeName: strings.Split(sheet.GetName(), "_")[1]}

				msgName := sheet.GetMessageType().GetFullyQualifiedName()
				mmd := fd.FindMessage(msgName)
				if mmd == nil {
					return errors.New(fmt.Sprintf("%s message not found in proto file", msgName))
				}

				for index, field := range mmd.GetFields() {
					fieldComment := strings.TrimSpace(*field.GetSourceInfo().LeadingComments)
					if !strings.Contains(fieldComment, "@name") {
						return errors.New("no field name found in proto file")
					}

					fieldName := fieldComment[strings.Index(fieldComment, "@name")+6:]
					firstSpaceIndex := strings.IndexFunc(fieldName, unicode.IsSpace)
					if firstSpaceIndex != -1 {
						fieldName = fieldName[:firstSpaceIndex]
					}
					util.sheetMap[sheetName].fieldMap[fieldName] = &FieldUtil{index: index, codeName: field.GetName()}
				}
			}

			return nil
		}
	}

	return errors.New("no wrapper found in proto file")
}

func (util *protoUtil) LoadData(excelFilePath string) error {
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
