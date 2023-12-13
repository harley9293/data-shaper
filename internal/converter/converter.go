package converter

import (
	"errors"
	"fmt"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/xuri/excelize/v2"
	"strings"
)

func ProtoToExcel(protoFilePath string, excelFilePath string) error {
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
			f := excelize.NewFile()
			excelName := comment[strings.Index(comment, "@wrapper")+9:] + ".xlsx"

			for _, sheet := range md.GetFields() {
				sheetComment := strings.TrimSpace(*sheet.GetSourceInfo().LeadingComments)
				if !strings.Contains(sheetComment, "@name") {
					return errors.New("no sheet name found in proto file")
				}

				sheetName := sheetComment[strings.Index(sheetComment, "@name")+6:]
				if _, err = f.NewSheet(sheetName); err != nil {
					return err
				}

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
					if err = f.SetCellValue(sheetName, fmt.Sprintf("%s%d", string(rune('A'+index)), 1), fieldName); err != nil {
						return err
					}
				}
			}

			_ = f.DeleteSheet("Sheet1")

			if err = f.SaveAs(excelFilePath + excelName); err != nil {
				return err
			}

			return nil
		}
	}

	return errors.New("no wrapper found in proto file")
}
