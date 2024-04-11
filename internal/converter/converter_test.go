package converter

import (
	"os"
	"testing"
)

func TestProtoToExcel(t *testing.T) {
	_ = os.Remove("../../test/ActivityCfg.xlsx")
	err := ProtoToExcel("../../test/ActivityCfg.proto", "../../test/")
	if err != nil {
		t.Errorf("ProtoToExcel() error = %v", err)
	}
}

func TestExcelToJson(t *testing.T) {
	err := ExcelToJson("../../test/ActivityCfg.proto", "../../test/ActivityCfg.xlsx", "../../test/")
	if err != nil {
		t.Errorf("ExcelToJson() error = %v", err)
	}
}
