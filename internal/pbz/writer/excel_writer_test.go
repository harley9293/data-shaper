package writer

import (
	"github.com/harley9293/data-shaper/internal/pbz/core"
	"github.com/harley9293/data-shaper/internal/pbz/parser"
	"os"
	"testing"
)

const TestExcelPath = "../test/"

func Test_saveData(t *testing.T) {
	schema := &core.ProtoExcelSchema{
		FilePath: TestExcelPath,
	}

	protoParser := &parser.ProtoParser{}
	err := protoParser.Parse(TestExcelPath+"test_all.proto", schema)
	if err != nil {
		t.Errorf("ParseProto() error = %v", err)
		return
	}

	excelWriter := &ExcelWriter{}
	err = excelWriter.Write(TestExcelPath+"测试配置读取.xlsx", schema)
	if err != nil {
		t.Errorf("saveData() err = %v", err)
		return
	}
	_ = os.Remove(TestExcelPath + "测试配置.xlsx")
}
