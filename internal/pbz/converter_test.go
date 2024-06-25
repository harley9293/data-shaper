package pbz

import (
	"os"
	"testing"
)

func TestProtoToExcel(t *testing.T) {
	_ = os.Remove("./test/测试配置.xlsx")
	err := ProtoToExcel("./test/test.proto", "./test/")
	if err != nil {
		t.Errorf("ProtoToExcel() error = %v", err)
	}
}

func TestExcelToCfg(t *testing.T) {
	err := ExcelToCfg("./test/ActivityCfg.proto", "./test/ActivityCfg.xlsx", "./test/")
	if err != nil {
		t.Errorf("ExcelToCfg() error = %v", err)
	}
}
