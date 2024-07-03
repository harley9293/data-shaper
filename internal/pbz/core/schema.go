package core

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
	dataParser  Parser
	excelWriter Writer
}

func NewProtoExcelSchema(excelDirPath string, protoParser Parser, excelWriter Writer, dataParser Parser) *ProtoExcelSchema {
	return &ProtoExcelSchema{FilePath: excelDirPath, protoParser: protoParser, excelWriter: excelWriter, dataParser: dataParser}
}

func (schema *ProtoExcelSchema) ParseProto(filePath string) error {
	return schema.protoParser.Parse(filePath, schema)
}

func (schema *ProtoExcelSchema) SaveData() error {
	return schema.excelWriter.Write(schema.FilePath, schema)
}

func (schema *ProtoExcelSchema) LoadData() {
	_ = schema.dataParser.Parse(schema.FilePath, schema)
}
