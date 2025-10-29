package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/petrovyuri/go-envied"
)

func TestGenerateConfigIntegration(t *testing.T) {
	// Create temporary directory for tests
	tempDir := t.TempDir()

	// Create test .env files
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
		t.Fatalf("Failed to create dev.env: %v", err)
	}

	err = os.WriteFile(prodEnvFile, []byte(prodContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create prod.env: %v", err)
	}

	// Create configuration file
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
		t.Fatalf("Failed to serialize configuration: %v", err)
	}

	err = os.WriteFile(configFile, configJSON, 0644)
	if err != nil {
		t.Fatalf("Failed to create config.json: %v", err)
	}

	// Load configuration
	loadedConfig, err := envied.LoadConfigFile(configFile)
	if err != nil {
		t.Fatalf("LoadConfigFile() returned error: %v", err)
	}

	// Check loaded configuration
	if loadedConfig.PackageName != "testconfig" {
		t.Errorf("PackageName = %q, expected %q", loadedConfig.PackageName, "testconfig")
	}

	if loadedConfig.OutputDir != tempDir {
		t.Errorf("OutputDir = %q, expected %q", loadedConfig.OutputDir, tempDir)
	}

	if loadedConfig.RandomSeed != 12345 {
		t.Errorf("RandomSeed = %d, expected %d", loadedConfig.RandomSeed, 12345)
	}

	// Check environments
	if len(loadedConfig.Environments) != 2 {
		t.Errorf("Expected 2 environments, got %d", len(loadedConfig.Environments))
	}

	// Check dev environment
	devEnv, exists := loadedConfig.Environments["dev"]
	if !exists {
		t.Error("Dev environment not found")
	} else {
		if devEnv.EnvFile != devEnvFile {
			t.Errorf("Dev EnvFile = %q, expected %q", devEnv.EnvFile, devEnvFile)
		}
		if devEnv.StructName != "DevConfig" {
			t.Errorf("Dev StructName = %q, expected %q", devEnv.StructName, "DevConfig")
		}
	}

	// Check prod environment
	prodEnv, exists := loadedConfig.Environments["prod"]
	if !exists {
		t.Error("Prod environment not found")
	} else {
		if prodEnv.EnvFile != prodEnvFile {
			t.Errorf("Prod EnvFile = %q, expected %q", prodEnv.EnvFile, prodEnvFile)
		}
		if prodEnv.StructName != "ProdConfig" {
			t.Errorf("Prod StructName = %q, expected %q", prodEnv.StructName, "ProdConfig")
		}
	}
}

func TestGenerateConfigWithInvalidFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Create configuration with non-existent .env files
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
		t.Fatalf("Failed to serialize configuration: %v", err)
	}

	err = os.WriteFile(configFile, configJSON, 0644)
	if err != nil {
		t.Fatalf("Failed to create config.json: %v", err)
	}

	// Load configuration
	_, err = envied.LoadConfigFile(configFile)
	if err != nil {
		t.Errorf("LoadConfigFile() should not return error for non-existent .env file: %v", err)
	}
}

func TestGenerateConfigWithInconsistentEnvironments(t *testing.T) {
	tempDir := t.TempDir()

	// Create .env files with different variables
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
		t.Fatalf("Failed to create dev.env: %v", err)
	}

	err = os.WriteFile(prodEnvFile, []byte(prodContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create prod.env: %v", err)
	}

	// Create configuration file
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
		t.Fatalf("Failed to serialize configuration: %v", err)
	}

	err = os.WriteFile(configFile, configJSON, 0644)
	if err != nil {
		t.Fatalf("Failed to create config.json: %v", err)
	}

	// Load configuration
	_, err = envied.LoadConfigFile(configFile)
	if err != nil {
		t.Errorf("LoadConfigFile() should not return error for inconsistent environments: %v", err)
	}
}

func TestFieldTypeEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected envied.FieldType
	}{
		{
			name:     "number with leading zeros",
			input:    "007",
			expected: envied.FieldTypeInt,
		},
		{
			name:     "negative float number",
			input:    "-0.5",
			expected: envied.FieldTypeFloat,
		},
		{
			name:     "scientific notation with uppercase E",
			input:    "1.23E+02",
			expected: envied.FieldTypeFloat,
		},
		{
			name:     "scientific notation with lowercase e",
			input:    "1.23e-02",
			expected: envied.FieldTypeFloat,
		},
		{
			name:     "bool with uppercase",
			input:    "TRUE",
			expected: envied.FieldTypeBool,
		},
		{
			name:     "bool with lowercase",
			input:    "false",
			expected: envied.FieldTypeBool,
		},
		{
			name:     "bool with mixed case",
			input:    "True",
			expected: envied.FieldTypeBool,
		},
		{
			name:     "number as bool",
			input:    "1",
			expected: envied.FieldTypeBool,
		},
		{
			name:     "zero as bool",
			input:    "0",
			expected: envied.FieldTypeBool,
		},
		{
			name:     "string with numbers",
			input:    "123abc",
			expected: envied.FieldTypeString,
		},
		{
			name:     "string with symbols",
			input:    "!@#$%",
			expected: envied.FieldTypeString,
		},
		{
			name:     "empty string",
			input:    "",
			expected: envied.FieldTypeString,
		},
		{
			name:     "string with spaces",
			input:    "  hello  ",
			expected: envied.FieldTypeString,
		},
		{
			name:     "only spaces",
			input:    "   ",
			expected: envied.FieldTypeString,
		},
		{
			name:     "string with newlines",
			input:    "hello\nworld",
			expected: envied.FieldTypeString,
		},
		{
			name:     "string with tabs",
			input:    "hello\tworld",
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

func TestObfuscationConsistency(t *testing.T) {
	// Test that same values with same seed produce same results
	testString := "test value"
	seed := int64(12345)

	keys1, values1 := envied.ObfuscateString(testString, seed)
	keys2, values2 := envied.ObfuscateString(testString, seed)

	// Check that results are identical
	if len(keys1) != len(keys2) || len(values1) != len(values2) {
		t.Error("Array lengths should be the same")
	}

	for i := range keys1 {
		if keys1[i] != keys2[i] {
			t.Errorf("Keys don't match at position %d: %d != %d", i, keys1[i], keys2[i])
		}
		if values1[i] != values2[i] {
			t.Errorf("Values don't match at position %d: %d != %d", i, values1[i], values2[i])
		}
	}

	// Check that both results can be decrypted
	decrypted1 := envied.DeobfuscateString(keys1, values1)
	decrypted2 := envied.DeobfuscateString(keys2, values2)

	if decrypted1 != testString {
		t.Errorf("First decryption result %q doesn't match original %q", decrypted1, testString)
	}

	if decrypted2 != testString {
		t.Errorf("Second decryption result %q doesn't match original %q", decrypted2, testString)
	}
}

func TestDifferentSeedsProduceDifferentResults(t *testing.T) {
	// Test that different seeds produce different results
	testString := "test value"
	seed1 := int64(12345)
	seed2 := int64(54321)

	keys1, values1 := envied.ObfuscateString(testString, seed1)
	keys2, values2 := envied.ObfuscateString(testString, seed2)

	// Check that results are different
	if len(keys1) != len(keys2) || len(values1) != len(values2) {
		t.Error("Array lengths should be the same")
	}

	// Check that at least one element differs
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
		t.Error("Keys should differ for different seeds")
	}

	if !valuesDifferent {
		t.Error("Values should differ for different seeds")
	}

	// Check that both results can be decrypted
	decrypted1 := envied.DeobfuscateString(keys1, values1)
	decrypted2 := envied.DeobfuscateString(keys2, values2)

	if decrypted1 != testString {
		t.Errorf("First decryption result %q doesn't match original %q", decrypted1, testString)
	}

	if decrypted2 != testString {
		t.Errorf("Second decryption result %q doesn't match original %q", decrypted2, testString)
	}
}
