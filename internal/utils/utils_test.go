package utils

import (
	"testing"
)

func TestGenerateShortCode(t *testing.T) {
	code, err := GenerateShortCode()
	if err != nil {
		t.Errorf("GenerateShortCode() error = %v", err)
		return
	}

	if len(code) != 6 {
		t.Errorf("GenerateShortCode() length = %d; want 6", len(code))
	}

	// Check for unique results (mostly)
	code2, _ := GenerateShortCode()
	if code == code2 {
		t.Errorf("GenerateShortCode() produced identical codes: %s", code)
	}
}
