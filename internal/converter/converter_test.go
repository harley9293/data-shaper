package converter

import (
	"os"
	"testing"
)

func TestProtoToExcel(t *testing.T) {
	err := ProtoToExcel("../../test/test.proto", "")
	if err != nil {
		t.Errorf("ProtoToExcel() error = %v", err)
	}

	_ = os.Remove("测试配置.xlsx")
}
