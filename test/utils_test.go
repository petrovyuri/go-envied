package test

import (
	"testing"

	"github.com/petrovyuri/go-envied"
)

func TestParseInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "positive number",
			input:    "123",
			expected: 123,
		},
		{
			name:     "negative number",
			input:    "-456",
			expected: -456,
		},
		{
			name:     "zero",
			input:    "0",
			expected: 0,
		},
		{
			name:     "large number",
			input:    "2147483647",
			expected: 2147483647,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "invalid string",
			input:    "abc",
			expected: 0,
		},
		{
			name:     "string with spaces",
			input:    "  123  ",
			expected: 0, // strconv.Atoi doesn't handle spaces
		},
		{
			name:     "float number",
			input:    "123.45",
			expected: 0, // strconv.Atoi doesn't handle float
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := envied.ParseInt(tt.input)
			if result != tt.expected {
				t.Errorf("ParseInt(%q) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "true",
			input:    "true",
			expected: true,
		},
		{
			name:     "false",
			input:    "false",
			expected: false,
		},
		{
			name:     "TRUE",
			input:    "TRUE",
			expected: true,
		},
		{
			name:     "FALSE",
			input:    "FALSE",
			expected: false,
		},
		{
			name:     "True",
			input:    "True",
			expected: true,
		},
		{
			name:     "False",
			input:    "False",
			expected: false,
		},
		{
			name:     "1",
			input:    "1",
			expected: true,
		},
		{
			name:     "0",
			input:    "0",
			expected: false,
		},
		{
			name:     "t",
			input:    "t",
			expected: true,
		},
		{
			name:     "f",
			input:    "f",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "invalid string",
			input:    "maybe",
			expected: false,
		},
		{
			name:     "string with spaces",
			input:    " true ",
			expected: false, // strconv.ParseBool doesn't handle spaces
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := envied.ParseBool(tt.input)
			if result != tt.expected {
				t.Errorf("ParseBool(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{
			name:     "positive number",
			input:    "123.45",
			expected: 123.45,
		},
		{
			name:     "negative number",
			input:    "-456.78",
			expected: -456.78,
		},
		{
			name:     "integer number",
			input:    "123",
			expected: 123.0,
		},
		{
			name:     "zero",
			input:    "0",
			expected: 0.0,
		},
		{
			name:     "zero with dot",
			input:    "0.0",
			expected: 0.0,
		},
		{
			name:     "scientific notation",
			input:    "1.23e+02",
			expected: 123.0,
		},
		{
			name:     "negative scientific notation",
			input:    "-1.23e-02",
			expected: -0.0123,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0.0,
		},
		{
			name:     "invalid string",
			input:    "abc",
			expected: 0.0,
		},
		{
			name:     "string with spaces",
			input:    "  123.45  ",
			expected: 0.0, // strconv.ParseFloat doesn't handle spaces
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := envied.ParseFloat(tt.input)
			if result != tt.expected {
				t.Errorf("ParseFloat(%q) = %f, expected %f", tt.input, result, tt.expected)
			}
		})
	}
}

func TestObfuscateString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		seed  int64
	}{
		{
			name:  "regular string",
			input: "hello world",
			seed:  12345,
		},
		{
			name:  "empty string",
			input: "",
			seed:  12345,
		},
		{
			name:  "string with symbols",
			input: "!@#$%^&*()",
			seed:  12345,
		},
		{
			name:  "string with cyrillic",
			input: "привет мир",
			seed:  12345,
		},
		{
			name:  "long string",
			input: "это очень длинная строка для тестирования обфускации",
			seed:  12345,
		},
		{
			name:  "string with numbers",
			input: "1234567890",
			seed:  12345,
		},
		{
			name:  "random seed",
			input: "test string",
			seed:  0, // will use random seed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, encryptedValues := envied.ObfuscateString(tt.input, tt.seed)

			// Check that array lengths match
			if len(keys) != len(encryptedValues) {
				t.Errorf("Key lengths (%d) and encrypted values (%d) don't match",
					len(keys), len(encryptedValues))
			}

			// Check that length matches input string length
			if len(keys) != len([]rune(tt.input)) {
				t.Errorf("Key length (%d) doesn't match input string length (%d)",
					len(keys), len([]rune(tt.input)))
			}

			// Check that we can decrypt back
			decrypted := envied.DeobfuscateString(keys, encryptedValues)
			if decrypted != tt.input {
				t.Errorf("Decrypted string (%q) doesn't match original (%q)",
					decrypted, tt.input)
			}
		})
	}
}

func TestDeobfuscateString(t *testing.T) {
	tests := []struct {
		name            string
		keys            []int
		encryptedValues []int
		expected        string
	}{
		{
			name:            "empty arrays",
			keys:            []int{},
			encryptedValues: []int{},
			expected:        "",
		},
		{
			name:            "different array lengths",
			keys:            []int{1, 2},
			encryptedValues: []int{3},
			expected:        "",
		},
		{
			name:            "single character",
			keys:            []int{100},
			encryptedValues: []int{100 ^ int('a')}, // 'a' = 97
			expected:        "a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := envied.DeobfuscateString(tt.keys, tt.encryptedValues)
			if result != tt.expected {
				t.Errorf("DeobfuscateString() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestObfuscateDeobfuscateRoundTrip(t *testing.T) {
	testStrings := []string{
		"",
		"a",
		"hello",
		"привет",
		"123456",
		"!@#$%^&*()",
		"многострочная\nстрока\tс\tтабуляцией",
		"строка с пробелами и символами: !@#$%^&*()_+-=[]{}|;':\",./<>?",
	}

	for _, testString := range testStrings {
		t.Run(testString, func(t *testing.T) {
			keys, encryptedValues := envied.ObfuscateString(testString, 12345)
			decrypted := envied.DeobfuscateString(keys, encryptedValues)

			if decrypted != testString {
				t.Errorf("Error in encrypt-decrypt round trip: original string %q, result %q",
					testString, decrypted)
			}
		})
	}
}
