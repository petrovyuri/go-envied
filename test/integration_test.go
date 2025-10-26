package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/petrovyuri/go-envied"
)

func TestGenerateConfigIntegration(t *testing.T) {
	// Создаем временную директорию для тестов
	tempDir := t.TempDir()

	// Создаем тестовые .env файлы
	devEnvFile := filepath.Join(tempDir, "dev.env")
	prodEnvFile := filepath.Join(tempDir, "prod.env")

	devContent := `# Dev environment
TOKEN=dev_token_123
API_URL=https://dev-api.example.com
PORT=8080
DEBUG=true
TIMEOUT=30.5
EMPTY_VALUE=
`

	prodContent := `# Prod environment
TOKEN=prod_token_456
API_URL=https://api.example.com
PORT=80
DEBUG=false
TIMEOUT=60.0
EMPTY_VALUE=
`

	err := os.WriteFile(devEnvFile, []byte(devContent), 0644)
	if err != nil {
		t.Fatalf("Не удалось создать dev.env: %v", err)
	}

	err = os.WriteFile(prodEnvFile, []byte(prodContent), 0644)
	if err != nil {
		t.Fatalf("Не удалось создать prod.env: %v", err)
	}

	// Создаем конфигурационный файл
	configFile := filepath.Join(tempDir, "config.json")
	config := envied.ConfigFile{
		PackageName: "testconfig",
		OutputDir:   tempDir,
		RandomSeed:  12345,
		Environments: map[string]envied.EnvironmentConfig{
			"dev": {
				EnvFile:    devEnvFile,
				StructName: "DevConfig",
			},
			"prod": {
				EnvFile:    prodEnvFile,
				StructName: "ProdConfig",
			},
		},
	}

	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("Не удалось сериализовать конфигурацию: %v", err)
	}

	err = os.WriteFile(configFile, configJSON, 0644)
	if err != nil {
		t.Fatalf("Не удалось создать config.json: %v", err)
	}

	// Загружаем конфигурацию
	loadedConfig, err := envied.LoadConfigFile(configFile)
	if err != nil {
		t.Fatalf("LoadConfigFile() вернул ошибку: %v", err)
	}

	// Проверяем загруженную конфигурацию
	if loadedConfig.PackageName != "testconfig" {
		t.Errorf("PackageName = %q, ожидалось %q", loadedConfig.PackageName, "testconfig")
	}

	if loadedConfig.OutputDir != tempDir {
		t.Errorf("OutputDir = %q, ожидалось %q", loadedConfig.OutputDir, tempDir)
	}

	if loadedConfig.RandomSeed != 12345 {
		t.Errorf("RandomSeed = %d, ожидалось %d", loadedConfig.RandomSeed, 12345)
	}

	// Проверяем окружения
	if len(loadedConfig.Environments) != 2 {
		t.Errorf("Ожидалось 2 окружения, получено %d", len(loadedConfig.Environments))
	}

	// Проверяем dev окружение
	devEnv, exists := loadedConfig.Environments["dev"]
	if !exists {
		t.Error("Dev окружение не найдено")
	} else {
		if devEnv.EnvFile != devEnvFile {
			t.Errorf("Dev EnvFile = %q, ожидалось %q", devEnv.EnvFile, devEnvFile)
		}
		if devEnv.StructName != "DevConfig" {
			t.Errorf("Dev StructName = %q, ожидалось %q", devEnv.StructName, "DevConfig")
		}
	}

	// Проверяем prod окружение
	prodEnv, exists := loadedConfig.Environments["prod"]
	if !exists {
		t.Error("Prod окружение не найдено")
	} else {
		if prodEnv.EnvFile != prodEnvFile {
			t.Errorf("Prod EnvFile = %q, ожидалось %q", prodEnv.EnvFile, prodEnvFile)
		}
		if prodEnv.StructName != "ProdConfig" {
			t.Errorf("Prod StructName = %q, ожидалось %q", prodEnv.StructName, "ProdConfig")
		}
	}
}

func TestGenerateConfigWithInvalidFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Создаем конфигурацию с несуществующими .env файлами
	configFile := filepath.Join(tempDir, "config.json")
	config := envied.ConfigFile{
		PackageName: "testconfig",
		OutputDir:   tempDir,
		RandomSeed:  12345,
		Environments: map[string]envied.EnvironmentConfig{
			"dev": {
				EnvFile:    "nonexistent.env",
				StructName: "DevConfig",
			},
		},
	}

	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("Не удалось сериализовать конфигурацию: %v", err)
	}

	err = os.WriteFile(configFile, configJSON, 0644)
	if err != nil {
		t.Fatalf("Не удалось создать config.json: %v", err)
	}

	// Загружаем конфигурацию
	_, err = envied.LoadConfigFile(configFile)
	if err != nil {
		t.Errorf("LoadConfigFile() не должен вернуть ошибку для несуществующего .env файла: %v", err)
	}
}

func TestGenerateConfigWithInconsistentEnvironments(t *testing.T) {
	tempDir := t.TempDir()

	// Создаем .env файлы с разными переменными
	devEnvFile := filepath.Join(tempDir, "dev.env")
	prodEnvFile := filepath.Join(tempDir, "prod.env")

	devContent := `TOKEN=dev_token
API_URL=https://dev-api.example.com
PORT=8080
DEBUG=true
`

	prodContent := `TOKEN=prod_token
API_URL=https://api.example.com
PORT=80
DEBUG=false
EXTRA_VAR=extra_value
`

	err := os.WriteFile(devEnvFile, []byte(devContent), 0644)
	if err != nil {
		t.Fatalf("Не удалось создать dev.env: %v", err)
	}

	err = os.WriteFile(prodEnvFile, []byte(prodContent), 0644)
	if err != nil {
		t.Fatalf("Не удалось создать prod.env: %v", err)
	}

	// Создаем конфигурационный файл
	configFile := filepath.Join(tempDir, "config.json")
	config := envied.ConfigFile{
		PackageName: "testconfig",
		OutputDir:   tempDir,
		RandomSeed:  12345,
		Environments: map[string]envied.EnvironmentConfig{
			"dev": {
				EnvFile:    devEnvFile,
				StructName: "DevConfig",
			},
			"prod": {
				EnvFile:    prodEnvFile,
				StructName: "ProdConfig",
			},
		},
	}

	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("Не удалось сериализовать конфигурацию: %v", err)
	}

	err = os.WriteFile(configFile, configJSON, 0644)
	if err != nil {
		t.Fatalf("Не удалось создать config.json: %v", err)
	}

	// Загружаем конфигурацию
	_, err = envied.LoadConfigFile(configFile)
	if err != nil {
		t.Errorf("LoadConfigFile() не должен вернуть ошибку для несогласованных окружений: %v", err)
	}
}

func TestFieldTypeEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected envied.FieldType
	}{
		{
			name:     "число с ведущими нулями",
			input:    "007",
			expected: envied.FieldTypeInt,
		},
		{
			name:     "отрицательное число с плавающей точкой",
			input:    "-0.5",
			expected: envied.FieldTypeFloat,
		},
		{
			name:     "научная нотация с заглавной E",
			input:    "1.23E+02",
			expected: envied.FieldTypeFloat,
		},
		{
			name:     "научная нотация с маленькой e",
			input:    "1.23e-02",
			expected: envied.FieldTypeFloat,
		},
		{
			name:     "bool с заглавными буквами",
			input:    "TRUE",
			expected: envied.FieldTypeBool,
		},
		{
			name:     "bool с маленькими буквами",
			input:    "false",
			expected: envied.FieldTypeBool,
		},
		{
			name:     "bool с смешанным регистром",
			input:    "True",
			expected: envied.FieldTypeBool,
		},
		{
			name:     "число как bool",
			input:    "1",
			expected: envied.FieldTypeBool,
		},
		{
			name:     "ноль как bool",
			input:    "0",
			expected: envied.FieldTypeBool,
		},
		{
			name:     "строка с числами",
			input:    "123abc",
			expected: envied.FieldTypeString,
		},
		{
			name:     "строка с символами",
			input:    "!@#$%",
			expected: envied.FieldTypeString,
		},
		{
			name:     "пустая строка",
			input:    "",
			expected: envied.FieldTypeString,
		},
		{
			name:     "строка с пробелами",
			input:    "  hello  ",
			expected: envied.FieldTypeString,
		},
		{
			name:     "только пробелы",
			input:    "   ",
			expected: envied.FieldTypeString,
		},
		{
			name:     "строка с переносами строк",
			input:    "hello\nworld",
			expected: envied.FieldTypeString,
		},
		{
			name:     "строка с табуляцией",
			input:    "hello\tworld",
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

func TestObfuscationConsistency(t *testing.T) {
	// Тестируем, что одинаковые значения с одинаковым seed дают одинаковый результат
	testString := "test value"
	seed := int64(12345)

	keys1, values1 := envied.ObfuscateString(testString, seed)
	keys2, values2 := envied.ObfuscateString(testString, seed)

	// Проверяем, что результаты идентичны
	if len(keys1) != len(keys2) || len(values1) != len(values2) {
		t.Error("Длины массивов должны быть одинаковыми")
	}

	for i := range keys1 {
		if keys1[i] != keys2[i] {
			t.Errorf("Ключи не совпадают на позиции %d: %d != %d", i, keys1[i], keys2[i])
		}
		if values1[i] != values2[i] {
			t.Errorf("Значения не совпадают на позиции %d: %d != %d", i, values1[i], values2[i])
		}
	}

	// Проверяем, что можно расшифровать оба результата
	decrypted1 := envied.DeobfuscateString(keys1, values1)
	decrypted2 := envied.DeobfuscateString(keys2, values2)

	if decrypted1 != testString {
		t.Errorf("Первый результат расшифровки %q не совпадает с исходным %q", decrypted1, testString)
	}

	if decrypted2 != testString {
		t.Errorf("Второй результат расшифровки %q не совпадает с исходным %q", decrypted2, testString)
	}
}

func TestDifferentSeedsProduceDifferentResults(t *testing.T) {
	// Тестируем, что разные seed дают разные результаты
	testString := "test value"
	seed1 := int64(12345)
	seed2 := int64(54321)

	keys1, values1 := envied.ObfuscateString(testString, seed1)
	keys2, values2 := envied.ObfuscateString(testString, seed2)

	// Проверяем, что результаты разные
	if len(keys1) != len(keys2) || len(values1) != len(values2) {
		t.Error("Длины массивов должны быть одинаковыми")
	}

	// Проверяем, что хотя бы один элемент отличается
	keysDifferent := false
	valuesDifferent := false

	for i := range keys1 {
		if keys1[i] != keys2[i] {
			keysDifferent = true
		}
		if values1[i] != values2[i] {
			valuesDifferent = true
		}
	}

	if !keysDifferent {
		t.Error("Ключи должны отличаться для разных seed")
	}

	if !valuesDifferent {
		t.Error("Значения должны отличаться для разных seed")
	}

	// Проверяем, что оба результата можно расшифровать
	decrypted1 := envied.DeobfuscateString(keys1, values1)
	decrypted2 := envied.DeobfuscateString(keys2, values2)

	if decrypted1 != testString {
		t.Errorf("Первый результат расшифровки %q не совпадает с исходным %q", decrypted1, testString)
	}

	if decrypted2 != testString {
		t.Errorf("Второй результат расшифровки %q не совпадает с исходным %q", decrypted2, testString)
	}
}
