package parser

import (
	"testing"
)

func TestDataParser_Parse(t *testing.T) {
	schema, err := Proto(TestExcelPath + "test_all.proto")
	if err != nil {
		t.Errorf("ParseProto() error = %v", err)
		return
	}

	err = Data(TestExcelPath+"测试配置读取.xlsx", schema)
	if err != nil {
		t.Errorf("Parse() error = %v", err)
		return
	}

	if schema.SheetList[0].ValueSize != 2 {
		t.Errorf("Parse() schema.SheetList[0].ValueSize = %v, want %v", schema.SheetList[0].ValueSize, 2)
		return
	}

	for _, field := range schema.SheetList[0].FieldList {
		if field.Name == "测试字段1" {
			if field.Values[0] != "1" || field.Values[1] != "2" {
				t.Errorf("Parse() field.Values = %v, want %v", field.Values, []string{"1", "2"})
				return
			}
		} else if field.Name == "测试字段2" {
			if field.Values[0] != "3" || field.Values[1] != "4" {
				t.Errorf("Parse() field.Values = %v, want %v", field.Values, []string{"3", "4"})
				return
			}
		} else if field.Name == "test_field_3" {
			if field.Values[0] != "5" || field.Values[1] != "6" {
				t.Errorf("Parse() field.Values = %v, want %v", field.Values, []string{"5", "6"})
				return
			}
		}
	}

	if schema.SheetList[1].ValueSize != 2 {
		t.Errorf("Parse() schema.SheetList[1].ValueSize = %v, want %v", schema.SheetList[1].ValueSize, 2)
		return
	}

	for _, field := range schema.SheetList[1].FieldList {
		if field.Name == "测试字段1" {
			if field.Values[0] != "1" || field.Values[1] != "10" {
				t.Errorf("Parse() field.Values = %v, want %v", field.Values, []string{"1", "10"})
				return
			}
		} else if field.Name == "测试字段2" {
			if field.Values[0] != "100" || field.Values[1] != "1000" {
				t.Errorf("Parse() field.Values = %v, want %v", field.Values, []string{"100", "1000"})
				return
			}
		} else if field.Name == "test_field_3" {
			if field.Values[0] != "10000" || field.Values[1] != "100000" {
				t.Errorf("Parse() field.Values = %v, want %v", field.Values, []string{"10000", "100000"})
				return
			}
		}
	}
}
