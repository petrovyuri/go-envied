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
			name:     "положительное число",
			input:    "123",
			expected: 123,
		},
		{
			name:     "отрицательное число",
			input:    "-456",
			expected: -456,
		},
		{
			name:     "ноль",
			input:    "0",
			expected: 0,
		},
		{
			name:     "большое число",
			input:    "2147483647",
			expected: 2147483647,
		},
		{
			name:     "пустая строка",
			input:    "",
			expected: 0,
		},
		{
			name:     "невалидная строка",
			input:    "abc",
			expected: 0,
		},
		{
			name:     "строка с пробелами",
			input:    "  123  ",
			expected: 0, // strconv.Atoi не обрабатывает пробелы
		},
		{
			name:     "число с плавающей точкой",
			input:    "123.45",
			expected: 0, // strconv.Atoi не обрабатывает float
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := envied.ParseInt(tt.input)
			if result != tt.expected {
				t.Errorf("ParseInt(%q) = %d, ожидалось %d", tt.input, result, tt.expected)
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
			name:     "пустая строка",
			input:    "",
			expected: false,
		},
		{
			name:     "невалидная строка",
			input:    "maybe",
			expected: false,
		},
		{
			name:     "строка с пробелами",
			input:    " true ",
			expected: false, // strconv.ParseBool не обрабатывает пробелы
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := envied.ParseBool(tt.input)
			if result != tt.expected {
				t.Errorf("ParseBool(%q) = %v, ожидалось %v", tt.input, result, tt.expected)
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
			name:     "положительное число",
			input:    "123.45",
			expected: 123.45,
		},
		{
			name:     "отрицательное число",
			input:    "-456.78",
			expected: -456.78,
		},
		{
			name:     "целое число",
			input:    "123",
			expected: 123.0,
		},
		{
			name:     "ноль",
			input:    "0",
			expected: 0.0,
		},
		{
			name:     "ноль с точкой",
			input:    "0.0",
			expected: 0.0,
		},
		{
			name:     "научная нотация",
			input:    "1.23e+02",
			expected: 123.0,
		},
		{
			name:     "отрицательная научная нотация",
			input:    "-1.23e-02",
			expected: -0.0123,
		},
		{
			name:     "пустая строка",
			input:    "",
			expected: 0.0,
		},
		{
			name:     "невалидная строка",
			input:    "abc",
			expected: 0.0,
		},
		{
			name:     "строка с пробелами",
			input:    "  123.45  ",
			expected: 0.0, // strconv.ParseFloat не обрабатывает пробелы
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := envied.ParseFloat(tt.input)
			if result != tt.expected {
				t.Errorf("ParseFloat(%q) = %f, ожидалось %f", tt.input, result, tt.expected)
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
			name:  "обычная строка",
			input: "hello world",
			seed:  12345,
		},
		{
			name:  "пустая строка",
			input: "",
			seed:  12345,
		},
		{
			name:  "строка с символами",
			input: "!@#$%^&*()",
			seed:  12345,
		},
		{
			name:  "строка с кириллицей",
			input: "привет мир",
			seed:  12345,
		},
		{
			name:  "длинная строка",
			input: "это очень длинная строка для тестирования обфускации",
			seed:  12345,
		},
		{
			name:  "строка с числами",
			input: "1234567890",
			seed:  12345,
		},
		{
			name:  "случайный seed",
			input: "test string",
			seed:  0, // будет использован случайный seed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, encryptedValues := envied.ObfuscateString(tt.input, tt.seed)

			// Проверяем, что длины массивов совпадают
			if len(keys) != len(encryptedValues) {
				t.Errorf("Длины ключей (%d) и зашифрованных значений (%d) не совпадают",
					len(keys), len(encryptedValues))
			}

			// Проверяем, что длина соответствует длине входной строки
			if len(keys) != len([]rune(tt.input)) {
				t.Errorf("Длина ключей (%d) не соответствует длине входной строки (%d)",
					len(keys), len([]rune(tt.input)))
			}

			// Проверяем, что можно расшифровать обратно
			decrypted := envied.DeobfuscateString(keys, encryptedValues)
			if decrypted != tt.input {
				t.Errorf("Расшифрованная строка (%q) не совпадает с исходной (%q)",
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
			name:            "пустые массивы",
			keys:            []int{},
			encryptedValues: []int{},
			expected:        "",
		},
		{
			name:            "разные длины массивов",
			keys:            []int{1, 2},
			encryptedValues: []int{3},
			expected:        "",
		},
		{
			name:            "один символ",
			keys:            []int{100},
			encryptedValues: []int{100 ^ int('a')}, // 'a' = 97
			expected:        "a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := envied.DeobfuscateString(tt.keys, tt.encryptedValues)
			if result != tt.expected {
				t.Errorf("DeobfuscateString() = %q, ожидалось %q", result, tt.expected)
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
				t.Errorf("Ошибка в цикле шифрование-расшифровка: исходная строка %q, результат %q",
					testString, decrypted)
			}
		})
	}
}
