package schema

import "google.golang.org/protobuf/types/descriptorpb"

type Field struct {
	Name        string
	MessageType descriptorpb.FieldDescriptorProto_Type
	MessageName string
	Note        string

	Values []string
}

type Sheet struct {
	Name        string
	MessageName string
	FieldList   []Field

	ValueSize int
}

type Proto struct {
	Name        string
	MessageName string
	SheetList   []Sheet
}
