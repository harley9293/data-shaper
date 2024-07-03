package writer

import (
	"github.com/harley9293/data-shaper/internal/pbz/core"
	"github.com/harley9293/data-shaper/internal/pbz/parser"
	"os"
	"testing"
)

const TestExcelPath = "../test/"

func Test_saveData(t *testing.T) {
	schema := core.NewProtoExcelSchema(TestExcelPath, &parser.ProtoParser{}, &ExcelWriter{})
	err := schema.ParseProto(TestExcelPath + "test_all.proto")
	if err != nil {
		t.Errorf("ParseProto() error = %v", err)
		return
	}

	err = schema.SaveData()
	if err != nil {
		t.Errorf("saveData() err = %v", err)
		return
	}
	_ = os.Remove(TestExcelPath + "测试配置.xlsx")
}
