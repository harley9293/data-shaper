package writer

import (
	"os"
	"testing"

	"github.com/harley9293/data-shaper/internal/pbz/parser"
)

const TestExcelPath = "../test/"

func Test_saveData(t *testing.T) {

	schema, err := parser.Proto(TestExcelPath + "test_all.proto")
	if err != nil {
		t.Errorf("ParseProto() error = %v", err)
		return
	}

	excelWriter := &ExcelWriter{}
	err = excelWriter.Write(TestExcelPath+"测试配置.xlsx", schema)
	if err != nil {
		t.Errorf("saveData() err = %v", err)
		return
	}
	_ = os.Remove(TestExcelPath + "测试配置.xlsx")
}
