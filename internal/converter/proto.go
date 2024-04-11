package converter

import (
	"errors"
	"fmt"
	"github.com/jhump/protoreflect/desc/protoparse"
	"strings"
)

type FieldUtil struct {
	index int
}

type sheetUtil struct {
	fieldMap map[string]FieldUtil
}

type protoUtil struct {
	excelFileName string
	sheetMap      map[string]sheetUtil
}

func (util *protoUtil) Parse(protoFilePath string) error {
	util.sheetMap = make(map[string]sheetUtil)

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
			util.excelFileName = comment[strings.Index(comment, "@wrapper")+9:] + ".xlsx"

			for _, sheet := range md.GetFields() {
				sheetComment := strings.TrimSpace(*sheet.GetSourceInfo().LeadingComments)
				if !strings.Contains(sheetComment, "@name") {
					return errors.New("no sheet name found in proto file")
				}

				sheetName := sheetComment[strings.Index(sheetComment, "@name")+6:]
				util.sheetMap[sheetName] = sheetUtil{fieldMap: make(map[string]FieldUtil)}

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
					util.sheetMap[sheetName].fieldMap[fieldName] = FieldUtil{index: index}
				}
			}

			return nil
		}
	}

	return errors.New("no wrapper found in proto file")
}
