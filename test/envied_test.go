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
			name:     "целое число",
			input:    "123",
			expected: envied.FieldTypeInt,
		},
		{
			name:     "отрицательное число",
			input:    "-456",
			expected: envied.FieldTypeInt,
		},
		{
			name:     "ноль",
			input:    "0",
			expected: envied.FieldTypeBool,
		},
		{
			name:     "число с плавающей точкой",
			input:    "123.45",
			expected: envied.FieldTypeFloat,
		},
		{
			name:     "отрицательное число с плавающей точкой",
			input:    "-456.78",
			expected: envied.FieldTypeFloat,
		},
		{
			name:     "научная нотация",
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
			name:     "обычная строка",
			input:    "hello world",
			expected: envied.FieldTypeString,
		},
		{
			name:     "пустая строка",
			input:    "",
			expected: envied.FieldTypeString,
		},
		{
			name:     "строка с символами",
			input:    "!@#$%^&*()",
			expected: envied.FieldTypeString,
		},
		{
			name:     "URL",
			input:    "https://example.com",
			expected: envied.FieldTypeString,
		},
		{
			name:     "невалидное число",
			input:    "123abc",
			expected: envied.FieldTypeString,
		},
		{
			name:     "строка с пробелами",
			input:    "  hello  ",
			expected: envied.FieldTypeString,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := envied.DetectFieldType(tt.input)
			if result != tt.expected {
				t.Errorf("DetectFieldType(%q) = %v, ожидалось %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLoadEnvFile(t *testing.T) {
	// Создаем временный .env файл для тестирования
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, "test.env")

	envContent := `# Тестовый .env файл
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
		t.Fatalf("Не удалось создать тестовый .env файл: %v", err)
	}

	// Тестируем загрузку файла
	fields, err := envied.LoadEnvFile(envFile)
	if err != nil {
		t.Fatalf("LoadEnvFile() вернул ошибку: %v", err)
	}

	// Проверяем количество полей
	expectedCount := 10 // все поля включая EMPTY_VALUE
	if len(fields) != expectedCount {
		t.Errorf("Ожидалось %d полей, получено %d", expectedCount, len(fields))
	}

	// Проверяем конкретные поля
	fieldMap := make(map[string]envied.Field)
	for _, field := range fields {
		fieldMap[field.EnvName] = field
	}

	// Проверяем STRING_VALUE
	if field, exists := fieldMap["STRING_VALUE"]; exists {
		if field.Value != "hello world" {
			t.Errorf("STRING_VALUE = %q, ожидалось %q", field.Value, "hello world")
		}
		if field.Type != envied.FieldTypeString {
			t.Errorf("STRING_VALUE тип = %v, ожидалось %v", field.Type, envied.FieldTypeString)
		}
	} else {
		t.Error("STRING_VALUE не найден")
	}

	// Проверяем INT_VALUE
	if field, exists := fieldMap["INT_VALUE"]; exists {
		if field.Value != "123" {
			t.Errorf("INT_VALUE = %q, ожидалось %q", field.Value, "123")
		}
		if field.Type != envied.FieldTypeInt {
			t.Errorf("INT_VALUE тип = %v, ожидалось %v", field.Type, envied.FieldTypeInt)
		}
	} else {
		t.Error("INT_VALUE не найден")
	}

	// Проверяем FLOAT_VALUE
	if field, exists := fieldMap["FLOAT_VALUE"]; exists {
		if field.Value != "123.45" {
			t.Errorf("FLOAT_VALUE = %q, ожидалось %q", field.Value, "123.45")
		}
		if field.Type != envied.FieldTypeFloat {
			t.Errorf("FLOAT_VALUE тип = %v, ожидалось %v", field.Type, envied.FieldTypeFloat)
		}
	} else {
		t.Error("FLOAT_VALUE не найден")
	}

	// Проверяем BOOL_VALUE
	if field, exists := fieldMap["BOOL_VALUE"]; exists {
		if field.Value != "true" {
			t.Errorf("BOOL_VALUE = %q, ожидалось %q", field.Value, "true")
		}
		if field.Type != envied.FieldTypeBool {
			t.Errorf("BOOL_VALUE тип = %v, ожидалось %v", field.Type, envied.FieldTypeBool)
		}
	} else {
		t.Error("BOOL_VALUE не найден")
	}

	// Проверяем EMPTY_VALUE
	if field, exists := fieldMap["EMPTY_VALUE"]; exists {
		if field.Value != "" {
			t.Errorf("EMPTY_VALUE = %q, ожидалось пустую строку", field.Value)
		}
		if field.Type != envied.FieldTypeString {
			t.Errorf("EMPTY_VALUE тип = %v, ожидалось %v", field.Type, envied.FieldTypeString)
		}
	} else {
		t.Error("EMPTY_VALUE не найден")
	}
}

func TestLoadEnvFileNotFound(t *testing.T) {
	// Тестируем загрузку несуществующего файла
	_, err := envied.LoadEnvFile("nonexistent.env")
	if err == nil {
		t.Error("LoadEnvFile() должен вернуть ошибку для несуществующего файла")
	}
}

func TestLoadEnvFileEmpty(t *testing.T) {
	// Создаем пустой .env файл
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, "empty.env")

	err := os.WriteFile(envFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Не удалось создать пустой .env файл: %v", err)
	}

	fields, err := envied.LoadEnvFile(envFile)
	if err != nil {
		t.Fatalf("LoadEnvFile() вернул ошибку: %v", err)
	}

	if len(fields) != 0 {
		t.Errorf("Ожидалось 0 полей для пустого файла, получено %d", len(fields))
	}
}

func TestLoadEnvFileOnlyComments(t *testing.T) {
	// Создаем .env файл только с комментариями
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, "comments.env")

	envContent := `# Это комментарий
# Еще один комментарий
# И еще один
`

	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Не удалось создать .env файл с комментариями: %v", err)
	}

	fields, err := envied.LoadEnvFile(envFile)
	if err != nil {
		t.Fatalf("LoadEnvFile() вернул ошибку: %v", err)
	}

	if len(fields) != 0 {
		t.Errorf("Ожидалось 0 полей для файла только с комментариями, получено %d", len(fields))
	}
}
