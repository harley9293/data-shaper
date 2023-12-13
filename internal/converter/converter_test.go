package converter

import "testing"

func TestProtoToExcel(t *testing.T) {
	err := ProtoToExcel("../../test/test.proto", "")
	if err != nil {
		t.Errorf("ProtoToExcel() error = %v", err)
	}
}
