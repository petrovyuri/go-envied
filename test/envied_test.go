package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/petrovyuri/go-envied"
)

func TestDetectFieldType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected envied.FieldType
	}{
		{
			name:     "integer number",
			input:    "123",
			expected: envied.FieldTypeInt,
		},
		{
			name:     "negative number",
			input:    "-456",
			expected: envied.FieldTypeInt,
		},
		{
			name:     "zero",
			input:    "0",
			expected: envied.FieldTypeBool,
		},
		{
			name:     "float number",
			input:    "123.45",
			expected: envied.FieldTypeFloat,
		},
		{
			name:     "negative float number",
			input:    "-456.78",
			expected: envied.FieldTypeFloat,
		},
		{
			name:     "scientific notation",
			input:    "1.23e+02",
			expected: envied.FieldTypeFloat,
		},
		{
			name:     "true",
			input:    "true",
			expected: envied.FieldTypeBool,
		},
		{
			name:     "false",
			input:    "false",
			expected: envied.FieldTypeBool,
		},
		{
			name:     "TRUE",
			input:    "TRUE",
			expected: envied.FieldTypeBool,
		},
		{
			name:     "FALSE",
			input:    "FALSE",
			expected: envied.FieldTypeBool,
		},
		{
			name:     "1",
			input:    "1",
			expected: envied.FieldTypeBool,
		},
		{
			name:     "0",
			input:    "0",
			expected: envied.FieldTypeBool,
		},
		{
			name:     "t",
			input:    "t",
			expected: envied.FieldTypeBool,
		},
		{
			name:     "f",
			input:    "f",
			expected: envied.FieldTypeBool,
		},
		{
			name:     "regular string",
			input:    "hello world",
			expected: envied.FieldTypeString,
		},
		{
			name:     "empty string",
			input:    "",
			expected: envied.FieldTypeString,
		},
		{
			name:     "string with symbols",
			input:    "!@#$%^&*()",
			expected: envied.FieldTypeString,
		},
		{
			name:     "URL",
			input:    "https://example.com",
			expected: envied.FieldTypeString,
		},
		{
			name:     "invalid number",
			input:    "123abc",
			expected: envied.FieldTypeString,
		},
		{
			name:     "string with spaces",
			input:    "  hello  ",
			expected: envied.FieldTypeString,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := envied.DetectFieldType(tt.input)
			if result != tt.expected {
				t.Errorf("DetectFieldType(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLoadEnvFile(t *testing.T) {
	// Create temporary .env file for testing
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, "test.env")

	envContent := `# Test .env file
STRING_VALUE=hello world
INT_VALUE=123
FLOAT_VALUE=123.45
BOOL_VALUE=true
EMPTY_VALUE=
QUOTED_VALUE="quoted string"
SPACES_VALUE=  value with spaces  
COMMENT_VALUE=value#with#hash
MULTILINE_VALUE=line1
MULTILINE_VALUE2=line2
`

	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}

	// Test file loading
	fields, err := envied.LoadEnvFile(envFile)
	if err != nil {
		t.Fatalf("LoadEnvFile() returned error: %v", err)
	}

	// Check field count
	expectedCount := 10 // all fields including EMPTY_VALUE
	if len(fields) != expectedCount {
		t.Errorf("Expected %d fields, got %d", expectedCount, len(fields))
	}

	// Check specific fields
	fieldMap := make(map[string]envied.Field)
	for _, field := range fields {
		fieldMap[field.EnvName] = field
	}

	// Check STRING_VALUE
	if field, exists := fieldMap["STRING_VALUE"]; exists {
		if field.Value != "hello world" {
			t.Errorf("STRING_VALUE = %q, expected %q", field.Value, "hello world")
		}
		if field.Type != envied.FieldTypeString {
			t.Errorf("STRING_VALUE type = %v, expected %v", field.Type, envied.FieldTypeString)
		}
	} else {
		t.Error("STRING_VALUE not found")
	}

	// Check INT_VALUE
	if field, exists := fieldMap["INT_VALUE"]; exists {
		if field.Value != "123" {
			t.Errorf("INT_VALUE = %q, expected %q", field.Value, "123")
		}
		if field.Type != envied.FieldTypeInt {
			t.Errorf("INT_VALUE type = %v, expected %v", field.Type, envied.FieldTypeInt)
		}
	} else {
		t.Error("INT_VALUE not found")
	}

	// Check FLOAT_VALUE
	if field, exists := fieldMap["FLOAT_VALUE"]; exists {
		if field.Value != "123.45" {
			t.Errorf("FLOAT_VALUE = %q, expected %q", field.Value, "123.45")
		}
		if field.Type != envied.FieldTypeFloat {
			t.Errorf("FLOAT_VALUE type = %v, expected %v", field.Type, envied.FieldTypeFloat)
		}
	} else {
		t.Error("FLOAT_VALUE not found")
	}

	// Check BOOL_VALUE
	if field, exists := fieldMap["BOOL_VALUE"]; exists {
		if field.Value != "true" {
			t.Errorf("BOOL_VALUE = %q, expected %q", field.Value, "true")
		}
		if field.Type != envied.FieldTypeBool {
			t.Errorf("BOOL_VALUE type = %v, expected %v", field.Type, envied.FieldTypeBool)
		}
	} else {
		t.Error("BOOL_VALUE not found")
	}

	// Check EMPTY_VALUE
	if field, exists := fieldMap["EMPTY_VALUE"]; exists {
		if field.Value != "" {
			t.Errorf("EMPTY_VALUE = %q, expected empty string", field.Value)
		}
		if field.Type != envied.FieldTypeString {
			t.Errorf("EMPTY_VALUE type = %v, expected %v", field.Type, envied.FieldTypeString)
		}
	} else {
		t.Error("EMPTY_VALUE not found")
	}
}

func TestLoadEnvFileNotFound(t *testing.T) {
	// Test loading non-existent file
	_, err := envied.LoadEnvFile("nonexistent.env")
	if err == nil {
		t.Error("LoadEnvFile() should return error for non-existent file")
	}
}

func TestLoadEnvFileEmpty(t *testing.T) {
	// Create empty .env file
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, "empty.env")

	err := os.WriteFile(envFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create empty .env file: %v", err)
	}

	fields, err := envied.LoadEnvFile(envFile)
	if err != nil {
		t.Fatalf("LoadEnvFile() returned error: %v", err)
	}

	if len(fields) != 0 {
		t.Errorf("Expected 0 fields for empty file, got %d", len(fields))
	}
}

func TestLoadEnvFileOnlyComments(t *testing.T) {
	// Create .env file with only comments
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, "comments.env")

	envContent := `# This is a comment
# Another comment
# And one more
`

	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file with comments: %v", err)
	}

	fields, err := envied.LoadEnvFile(envFile)
	if err != nil {
		t.Fatalf("LoadEnvFile() returned error: %v", err)
	}

	if len(fields) != 0 {
		t.Errorf("Expected 0 fields for file with only comments, got %d", len(fields))
	}
}
