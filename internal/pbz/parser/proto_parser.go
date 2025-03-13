package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/harley9293/data-shaper/internal/pbz/schema"
	"github.com/jhump/protoreflect/desc/protoparse"
)

func Proto(filePath string) (*schema.Proto, error) {
	if !strings.HasSuffix(filePath, "proto") {
		return nil, errors.New("必须是proto文件")
	}

	parse := protoparse.Parser{
		IncludeSourceCodeInfo: true,
	}

	files, err := parse.ParseFiles(filePath)
	if err != nil {
		return nil, err
	}

	proto := &schema.Proto{}
	fd := files[0]
	mds := fd.GetMessageTypes()

	for _, md := range mds {
		if md.GetSourceInfo() == nil || md.GetSourceInfo().LeadingComments == nil {
			continue
		}

		if hasKeyFromComments(md.GetSourceInfo().LeadingComments, "wrapper") {
			proto.Name = getValueFromComments(md.GetSourceInfo().LeadingComments, "wrapper", md.GetName())
			proto.MessageName = md.GetName()
			for _, sheet := range md.GetFields() {
				newSheet := schema.Sheet{Name: getValueFromComments(sheet.GetSourceInfo().LeadingComments, "name", sheet.GetName()), MessageName: sheet.GetName()}
				mmd := fd.FindMessage(sheet.GetMessageType().GetFullyQualifiedName())
				if mmd == nil {
					return nil, errors.New(fmt.Sprintf("%s message not found in proto file", sheet.GetMessageType().GetFullyQualifiedName()))
				}

				for _, field := range mmd.GetFields() {
					fieldName := getValueFromComments(field.GetSourceInfo().LeadingComments, "name", field.GetName())
					note := getValueFromComments(field.GetSourceInfo().LeadingComments, "note", "")
					newSheet.FieldList = append(newSheet.FieldList, schema.Field{Name: fieldName, MessageName: field.GetName(), MessageType: field.GetType(), Note: note})
				}
				proto.SheetList = append(proto.SheetList, newSheet)
			}

			return proto, nil
		}
	}

	return nil, errors.New("没有wrapper found in proto file")
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
