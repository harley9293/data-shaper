package pbz

import "testing"

func Test_parseProto_all(t *testing.T) {
	file, err := parseProto("./test/test_all.proto")
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

	sheet1, ok := file.sheetMap["测试页签1"]
	if !ok {
		t.Errorf("parseProto() file.sheetMap[\"#测试页签1\"] = %v, want %v", ok, true)
		return
	}

	_, ok = file.sheetMap["test_sheet_2"]
	if !ok {
		t.Errorf("parseProto() file.sheetMap[\"test_sheet_2\"] = %v, want %v", ok, true)
		return
	}

	if sheet1.sheetName != "#测试页签1" {
		t.Errorf("parseProto() sheet1.sheetName = %v, want %v", sheet1.sheetName, "#测试页签1")
		return
	}

	if sheet1.messageName != "test_sheet_1" {
		t.Errorf("parseProto() sheet1.messageName = %v, want %v", sheet1.messageName, "test_sheet_1")
		return
	}

	field1, ok := sheet1.fieldMap["测试字段1"]
	if !ok {
		t.Errorf("parseProto() sheet1.fieldMap[\"测试字段1\"] = %v, want %v", ok, true)
		return
	}

	if field1.messageName != "test_field_1" {
		t.Errorf("parseProto() field1.messageName = %v, want %v", field1.messageName, "test_field_1")
		return
	}

	if field1.index != 0 {
		t.Errorf("parseProto() field1.index = %v, want %v", field1.index, 0)
		return
	}

	field2, ok := sheet1.fieldMap["测试字段2"]
	if !ok {
		t.Errorf("parseProto() sheet1.fieldMap[\"测试字段2\"] = %v, want %v", ok, true)
		return
	}

	if field2.messageName != "test_field_2" {
		t.Errorf("parseProto() field2.messageName = %v, want %v", field2.messageName, "test_field_2")
		return
	}

	if field2.index != 1 {
		t.Errorf("parseProto() field2.index = %v, want %v", field2.index, 1)
		return
	}

	field3, ok := sheet1.fieldMap["test_field_3"]
	if !ok {
		t.Errorf("parseProto() sheet1.fieldMap[\"test_field_3\"] = %v, want %v", ok, true)
		return
	}

	if field3.messageName != "test_field_3" {
		t.Errorf("parseProto() field3.messageName = %v, want %v", field3.messageName, "test_field_3")
		return
	}

	if field3.index != 2 {
		t.Errorf("parseProto() field3.index = %v, want %v", field3.index, 2)
		return
	}
}

func Test_parseProto_defaultWrapper(t *testing.T) {
	file, err := parseProto("./test/test_wrapper_default.proto")
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
	_, err := parseProto("./test/test_wrapper_no.proto")
	if err == nil {
		t.Errorf("parseProto() error = nil, want no wrapper found in proto file")
		return
	}
}
