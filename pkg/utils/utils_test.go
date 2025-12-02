package utils

import "testing"

func TestParseJsonStr(t *testing.T) {
	var (
		sInt   = "0"
		sStr   = "\"hello\""
		sBool  = "true"
		sJson  = `{"key": "value"}`
		sFloat = "3.14"
	)

	resInt, err := ParseJsonStr[int](sInt)
	if err != nil {
		t.Errorf("ParseStr[int] failed: %v", err)
	} else {
		t.Logf("ParseStr[int] success: %d", resInt)
	}

	resStr, err := ParseJsonStr[string](sStr)
	if err != nil {
		t.Errorf("ParseStr[string] failed: %v", err)
	} else {
		t.Logf("ParseStr[string] success: %s", resStr)
	}

	resBool, err := ParseJsonStr[bool](sBool)
	if err != nil {
		t.Errorf("ParseStr[bool] failed: %v", err)
	} else {
		t.Logf("ParseStr[bool] success: %v", resBool)
	}

	resJson, err := ParseJsonStr[map[string]string](sJson)
	if err != nil {
		t.Errorf("ParseStr[map[string]string] failed: %v", err)
	} else {
		t.Logf("ParseStr[map[string]string] success: %v", resJson)
	}

	resH, err := ParseJsonStr[H](sJson)
	if err != nil {
		t.Errorf("ParseStr[H] failed: %v", err)
	} else {
		t.Logf("ParseStr[H] success: %v", resH)
	}

	resFloat, err := ParseJsonStr[float64](sFloat)
	if err != nil {
		t.Errorf("ParseStr[float64] failed: %v", err)
	} else {
		t.Logf("ParseStr[float64] success: %f", resFloat)
	}
}
