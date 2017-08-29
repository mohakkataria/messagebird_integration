package util

import (
	"regexp"
	"testing"
)

func TestConvertIntToHex(t *testing.T) {
	expected := "01"
	actual := ConvertIntToHex(1)
	if actual != expected {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", actual, expected)
	}
}

func TestIsFloat64(t *testing.T) {
	expected := true
	actual := IsFloat64(123.12)
	if actual != expected {
		t.Errorf("Test failed, expected: '%t', got:  '%t'", actual, expected)
	}

	expected = false
	actual = IsFloat64(123)
	if actual != expected {
		t.Errorf("Test failed, expected: '%t', got:  '%t'", actual, expected)
	}
}

func TestIsString(t *testing.T) {
	expected := true
	actual := IsString("1")
	if actual != expected {
		t.Errorf("Test failed, expected: '%t', got:  '%t'", actual, expected)
	}

	expected = false
	actual = IsString(123)
	if actual != expected {
		t.Errorf("Test failed, expected: '%t', got:  '%t'", actual, expected)
	}
}

func TestIsUnicode(t *testing.T) {
	expected := true
	actual := IsUnicode("日本語")
	if actual != expected {
		t.Errorf("Test failed, expected: '%t', got:  '%t'", actual, expected)
	}
}

func TestGetRandom2DigitHex(t *testing.T) {
	actual := GetRandom2DigitHex()
	if !regexp.MustCompile(`^[A-F0-9][A-F0-9]$`).MatchString(actual) {
		t.Errorf("Test failed. Not a valid random hex %s", actual)
	}
}
