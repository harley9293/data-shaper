package parser

import (
	"testing"

	"google.golang.org/protobuf/types/descriptorpb"
)

const TestExcelPath = "../test/"

func Test_parseProto_all(t *testing.T) {
	proto, err := Proto(TestExcelPath + "test_all.proto")
	if err != nil {
		t.Errorf("ParseProto() error = %v", err)
		return
	}

	if proto.Name != "测试配置" {
		t.Errorf("parseProto() proto.Name = %v, want %v", proto.Name, "测试配置")
		return
	}

	if proto.MessageName != "TestConfig" {
		t.Errorf("parseProto() proto.MessageName = %v, want %v", proto.MessageName, "TestConfig")
		return
	}

	if len(proto.SheetList) != 2 {
		t.Errorf("parseProto() len(proto.SheetList) = %v, want %v", len(proto.SheetList), 2)
		return
	}

	sheet1 := proto.SheetList[0]
	if sheet1.Name != "测试页签1" {
		t.Errorf("parseProto() sheet1.sheetName = %v, want %v", sheet1.Name, "测试页签1")
		return
	}

	if sheet1.MessageName != "test_sheet_1" {
		t.Errorf("parseProto() sheet1.MessageName = %v, want %v", sheet1.MessageName, "test_sheet_1")
		return
	}

	if len(sheet1.FieldList) != 3 {
		t.Errorf("parseProto() len(sheet1.FieldList) = %v, want %v", len(sheet1.FieldList), 3)
		return
	}

	field1 := sheet1.FieldList[0]
	if field1.Name != "测试字段1" {
		t.Errorf("parseProto() sheet1.FieldList[0].fieldName = %v, want %v", field1.Name, "测试字段1")
		return
	}

	if field1.MessageName != "test_field_1" {
		t.Errorf("parseProto() field1.MessageName = %v, want %v", field1.MessageName, "test_field_1")
		return
	}

	if field1.MessageType != descriptorpb.FieldDescriptorProto_TYPE_STRING {
		t.Errorf("parseProto() field1.MessageType = %v, want %v", field1.MessageType, int(descriptorpb.FieldDescriptorProto_TYPE_STRING))
		return
	}

	field2 := sheet1.FieldList[1]
	if field2.Name != "测试字段2" {
		t.Errorf("parseProto() sheet1.FieldList[1].fieldName = %v, want %v", field2.Name, "测试字段2")
		return
	}

	if field2.MessageName != "test_field_2" {
		t.Errorf("parseProto() field2.MessageName = %v, want %v", field2.MessageName, "test_field_2")
		return
	}

	if field2.MessageType != descriptorpb.FieldDescriptorProto_TYPE_INT32 {
		t.Errorf("parseProto() field2.MessageType = %v, want %v", field2.MessageType, int(descriptorpb.FieldDescriptorProto_TYPE_INT32))
		return
	}

	if field2.Note != "测试注释第一行\n测试注释第二行" {
		t.Errorf("parseProto() field2.Note = %v, want %v", field2.Note, "测试注释第一行\n测试注释第二行")
		return
	}

	field3 := sheet1.FieldList[2]
	if field3.Name != "test_field_3" {
		t.Errorf("parseProto() sheet1.FieldList[2].fieldName = %v, want %v", field3.Name, "test_field_3")
		return
	}

	if field3.MessageName != "test_field_3" {
		t.Errorf("parseProto() field3.MessageName = %v, want %v", field3.MessageName, "test_field_3")
		return
	}

	if field3.MessageType != descriptorpb.FieldDescriptorProto_TYPE_DOUBLE {
		t.Errorf("parseProto() field3.MessageType = %v, want %v", field3.MessageType, int(descriptorpb.FieldDescriptorProto_TYPE_DOUBLE))
		return
	}

	if field3.Note != "测试注释单行" {
		t.Errorf("parseProto() field3.Note = %v, want %v", field3.Note, "测试注释单行")
		return
	}

	sheet2 := proto.SheetList[1]

	if sheet2.Name != "test_sheet_2" {
		t.Errorf("parseProto() sheet2.Name = %v, want %v", sheet2.Name, "test_sheet_2")
		return
	}
}

func Test_parseProto_defaultWrapper(t *testing.T) {
	proto, err := Proto(TestExcelPath + "test_wrapper_default.proto")
	if err != nil {
		t.Errorf("ParseProto() error = %v", err)
		return
	}

	if proto.Name != "TestConfig" {
		t.Errorf("parseProto() proto.Name = %v, want %v", proto.Name, "TestConfig")
		return
	}
}

func Test_parseProto_noWrapper(t *testing.T) {
	_, err := Proto(TestExcelPath + "test_wrapper_no.proto")
	if err == nil {
		t.Errorf("ParseProto() error = %v, want %v", err, "no wrapper found in proto file")
		return
	}
}
