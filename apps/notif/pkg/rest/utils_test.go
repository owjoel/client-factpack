package rest

import "testing"

func TestParseIntWithDefault_ValidPositiveNumber(t *testing.T) {
	result := parseIntWithDefault("5", 10)
	if result != 5 {
		t.Errorf("Expected 5, got %d", result)
	}
}

func TestParseIntWithDefault_ZeroInput(t *testing.T) {
	result := parseIntWithDefault("0", 10)
	if result != 10 {
		t.Errorf("Expected default 10 for input 0, got %d", result)
	}
}

func TestParseIntWithDefault_NegativeInput(t *testing.T) {
	result := parseIntWithDefault("-3", 7)
	if result != 7 {
		t.Errorf("Expected default 7 for negative input, got %d", result)
	}
}

func TestParseIntWithDefault_InvalidString(t *testing.T) {
	result := parseIntWithDefault("abc", 42)
	if result != 42 {
		t.Errorf("Expected default 42 for invalid input, got %d", result)
	}
}

func TestParseIntWithDefault_EmptyString(t *testing.T) {
	result := parseIntWithDefault("", 3)
	if result != 3 {
		t.Errorf("Expected default 3 for empty string, got %d", result)
	}
}
