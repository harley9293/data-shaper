package core

type Parser interface {
	Parse(filePath string, protoSchema *ProtoExcelSchema) error
}

type Writer interface {
	Write(filePath string, protoSchema *ProtoExcelSchema) error
}
