package pbz

import (
	"os"
	"testing"
)

func Test_parseProto_all(t *testing.T) {
	file, err := parseProto("./test/test_all.proto", "./test/", false)
	if err != nil {
		t.Errorf("parseProto() error = %v", err)
		return
	}

	if file.fileName != "测试配置.xlsx" {
		t.Errorf("parseProto() file.excelFileName = %v, want %v", file.fileName, "测试配置.xlsx")
		return
	}

	if file.messageName != "TestConfig" {
		t.Errorf("parseProto() file.messageName = %v, want %v", file.messageName, "TestConfig")
		return
	}

	if len(file.sheetList) != 2 {
		t.Errorf("parseProto() len(file.sheetList) = %v, want %v", len(file.sheetList), 2)
		return
	}

	sheet1 := file.sheetList[0]
	if sheet1.sheetName != "测试页签1" {
		t.Errorf("parseProto() file.sheetList[0].sheetName = %v, want %v", file.sheetList[0].sheetName, "测试页签1")
		return
	}

	if file.sheetList[1].sheetName != "test_sheet_2" {
		t.Errorf("parseProto() file.sheetList[1].sheetName = %v, want %v", file.sheetList[1].sheetName, "test_sheet_2")
		return
	}

	if sheet1.messageName != "test_sheet_1" {
		t.Errorf("parseProto() sheet1.messageName = %v, want %v", sheet1.messageName, "test_sheet_1")
		return
	}

	if len(sheet1.fieldList) != 3 {
		t.Errorf("parseProto() len(sheet1.fieldList) = %v, want %v", len(sheet1.fieldList), 3)
		return
	}

	field1 := sheet1.fieldList[0]
	if field1.fieldName != "测试字段1" {
		t.Errorf("parseProto() sheet1.fieldList[0].fieldName = %v, want %v", field1.fieldName, "测试字段1")
		return
	}

	if field1.messageName != "test_field_1" {
		t.Errorf("parseProto() field1.messageName = %v, want %v", field1.messageName, "test_field_1")
		return
	}

	field2 := sheet1.fieldList[1]
	if field2.fieldName != "测试字段2" {
		t.Errorf("parseProto() sheet1.fieldList[1].fieldName = %v, want %v", field2.fieldName, "测试字段2")
		return
	}

	if field2.messageName != "test_field_2" {
		t.Errorf("parseProto() field2.messageName = %v, want %v", field2.messageName, "test_field_2")
		return
	}

	field3 := sheet1.fieldList[2]
	if field3.fieldName != "test_field_3" {
		t.Errorf("parseProto() sheet1.fieldList[2].fieldName = %v, want %v", field3.fieldName, "test_field_3")
		return
	}

	if field3.messageName != "test_field_3" {
		t.Errorf("parseProto() field3.messageName = %v, want %v", field3.messageName, "test_field_3")
		return
	}
}

func Test_parseProto_defaultWrapper(t *testing.T) {
	file, err := parseProto("./test/test_wrapper_default.proto", "./test/", false)
	if err != nil {
		t.Errorf("parseProto() error = %v", err)
		return
	}

	if file.fileName != "TestConfig.xlsx" {
		t.Errorf("parseProto() file.excelFileName = %v, want %v", file.fileName, "TestConfig.xlsx")
		return
	}
}

func Test_parseProto_noWrapper(t *testing.T) {
	_, err := parseProto("./test/test_wrapper_no.proto", "./test/", false)
	if err == nil {
		t.Errorf("parseProto() error = nil, want no wrapper found in proto file")
		return
	}
}

func Test_saveData(t *testing.T) {
	file, err := parseProto("./test/test_all.proto", "./test/", false)
	if err != nil {
		t.Errorf("saveData() parseProto error = %v", err)
		return
	}

	err = file.saveData()
	if err != nil {
		t.Errorf("saveData() err = %v", err)
		return
	}
	_ = os.Remove("./test/测试配置.xlsx")
}
