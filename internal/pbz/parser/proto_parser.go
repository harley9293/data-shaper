package parser

import (
	"errors"
	"fmt"
	"github.com/harley9293/data-shaper/internal/pbz/core"
	"github.com/jhump/protoreflect/desc/protoparse"
	"regexp"
	"strings"
)

type ProtoParser struct{}

func (p *ProtoParser) Parse(filePath string, protoSchema *core.ProtoExcelSchema) error {
	if !strings.HasSuffix(filePath, "proto") {
		return errors.New("file must be a proto file")
	}

	parse := protoparse.Parser{
		IncludeSourceCodeInfo: true,
	}

	files, err := parse.ParseFiles(filePath)
	if err != nil {
		return err
	}

	fd := files[0]
	mds := fd.GetMessageTypes()

	for _, md := range mds {
		if md.GetSourceInfo() == nil || md.GetSourceInfo().LeadingComments == nil {
			continue
		}

		if hasKeyFromComments(md.GetSourceInfo().LeadingComments, "wrapper") {
			protoSchema.FilePath += getValueFromComments(md.GetSourceInfo().LeadingComments, "wrapper", md.GetName()) + ".xlsx"
			protoSchema.MessageName = md.GetName()
			for _, sheet := range md.GetFields() {
				newSheet := core.SheetSchema{Name: getValueFromComments(sheet.GetSourceInfo().LeadingComments, "name", sheet.GetName()), MessageName: sheet.GetName(), Repeated: sheet.IsRepeated()}
				mmd := fd.FindMessage(sheet.GetMessageType().GetFullyQualifiedName())
				if mmd == nil {
					return errors.New(fmt.Sprintf("%s message not found in proto file", sheet.GetMessageType().GetFullyQualifiedName()))
				}

				for _, field := range mmd.GetFields() {
					fieldName := getValueFromComments(field.GetSourceInfo().LeadingComments, "name", field.GetName())
					note := getValueFromComments(field.GetSourceInfo().LeadingComments, "note", "")
					newSheet.FieldList = append(newSheet.FieldList, core.FieldSchema{Name: fieldName, MessageName: field.GetName(), MessageType: field.GetType(), Note: note})
				}
				protoSchema.SheetList = append(protoSchema.SheetList, newSheet)
			}

			return nil
		}
	}

	return errors.New("no wrapper found in proto file")
}

func hasKeyFromComments(comments *string, key string) bool {
	if comments == nil {
		return false
	}

	return strings.Contains(*comments, "@"+key)
}

func getValueFromComments(comments *string, key string, defaultValue string) (value string) {
	value = defaultValue
	if comments == nil {
		return
	}

	re := regexp.MustCompile(`@` + key + `\s+(.*)`)
	matches := re.FindAllStringSubmatch(*comments, -1)

	if matches != nil {
		var values []string
		for _, match := range matches {
			if len(match) > 1 {
				values = append(values, strings.TrimSpace(match[1]))
			}
		}
		value = strings.Join(values, "\n")
	}

	if value == "" {
		value = defaultValue
	}

	return
}
