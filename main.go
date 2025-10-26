// It reads environment variables from .env files, encrypts sensitive data, and generates
// type-safe Go configuration files.
package envied

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// FieldType represents the type of a configuration field
type FieldType string

const (
	FieldTypeString FieldType = "string"
	FieldTypeInt    FieldType = "int"
	FieldTypeBool   FieldType = "bool"
	FieldTypeFloat  FieldType = "float64"
)

// Field represents a configuration field
type Field struct {
	EnvName      string    // Environment variable name (used as field name)
	Type         FieldType // Field type
	Value        string    // Field value
	DefaultValue string    // Default value if env var is not set
	Optional     bool      // Whether the field is optional
}

// ObfuscationResult contains the obfuscated field data
type ObfuscationResult struct {
	KeyName   string
	ValueName string
	Key       interface{}
	Value     interface{}
}

// Config represents the configuration generation settings
type Config struct {
	PackageName string  // Go package name
	Environment string  // Environment name (dev, prod, etc.)
	Fields      []Field // Configuration fields
	OutputDir   string  // Output directory for generated files
}

// Generator handles configuration file generation
type Generator struct {
	config *Config
}

// ConfigFile structure for configuration file
type ConfigFile struct {
	PackageName  string                       `json:"package_name"`
	OutputDir    string                       `json:"output_dir"`
	RandomSeed   int                          `json:"random_seed,omitempty"`
	Environments map[string]EnvironmentConfig `json:"environments"`
}

type EnvironmentConfig struct {
	EnvFile    string `json:"env_file"`
	StructName string `json:"struct_name"`
}

// ObfuscateString obfuscates a string value using XOR with random keys for each character
func ObfuscateString(value string, seed int64) ([]int, []int) {
	var r *rand.Rand
	if seed == 0 {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	} else {
		r = rand.New(rand.NewSource(seed))
	}

	runes := []rune(value)
	keys := make([]int, len(runes))
	encryptedValues := make([]int, len(runes))

	for i, char := range runes {
		key := r.Intn(1 << 32)
		keys[i] = key
		encryptedValues[i] = int(char) ^ key
	}

	return keys, encryptedValues
}

// DeobfuscateString deobfuscates a string value using XOR with the keys
func DeobfuscateString(keys, encryptedValues []int) string {
	if len(keys) != len(encryptedValues) {
		return ""
	}

	runes := make([]rune, len(keys))
	for i := range keys {
		runes[i] = rune(keys[i] ^ encryptedValues[i])
	}

	return string(runes)
}

// ParseInt converts a string to int
func ParseInt(value string) int {
	result, _ := strconv.Atoi(value)
	return result
}

// ParseBool converts a string to bool
func ParseBool(value string) bool {
	result, _ := strconv.ParseBool(value)
	return result
}

// ParseFloat converts a string to float64
func ParseFloat(value string) float64 {
	result, _ := strconv.ParseFloat(value, 64)
	return result
}

// Deobfuscate deobfuscates a value using simple XOR obfuscation
// Similar to the original envied package for Dart/Flutter
func Deobfuscate(obfuscatedValue string, key string) string {
	if obfuscatedValue == "" {
		return ""
	}

	// Decode base64
	data, err := base64.StdEncoding.DecodeString(obfuscatedValue)
	if err != nil {
		fmt.Printf("Error decoding base64: %v\n", err)
		return ""
	}

	// Simple XOR deobfuscation with provided key
	keyBytes := []byte(key)
	result := make([]byte, len(data))

	for i := 0; i < len(data); i++ {
		result[i] = data[i] ^ keyBytes[i%len(keyBytes)]
	}

	return string(result)
}

// DeobfuscateWithDefaultKey deobfuscates a value using default key
// For backward compatibility
func DeobfuscateWithDefaultKey(obfuscatedValue string) string {
	return Deobfuscate(obfuscatedValue, "go-envied-obfuscation")
}

// Obfuscate obfuscates a value using simple XOR obfuscation
// Similar to the original envied package for Dart/Flutter
func Obfuscate(value string, key string) string {
	if value == "" {
		return ""
	}

	// Simple XOR obfuscation with provided key
	keyBytes := []byte(key)
	data := []byte(value)
	result := make([]byte, len(data))

	for i := 0; i < len(data); i++ {
		result[i] = data[i] ^ keyBytes[i%len(keyBytes)]
	}

	// Encode as base64
	return base64.StdEncoding.EncodeToString(result)
}

// generateObfuscatedField generates obfuscated field data based on type and value
func generateObfuscatedField(fieldName string, fieldType FieldType, value string, seed int64) (*ObfuscationResult, error) {
	switch fieldType {
	case FieldTypeString:
		keys, encryptedValues := ObfuscateString(value, seed)
		return &ObfuscationResult{
			KeyName:   fmt.Sprintf("_enviedkey%s", fieldName),
			ValueName: fmt.Sprintf("_envieddata%s", fieldName),
			Key:       keys,
			Value:     encryptedValues,
		}, nil

	case FieldTypeFloat:
		// For float64, we'll treat it as string for now
		keys, encryptedValues := ObfuscateString(value, seed)
		return &ObfuscationResult{
			KeyName:   fmt.Sprintf("_enviedkey%s", fieldName),
			ValueName: fmt.Sprintf("_envieddata%s", fieldName),
			Key:       keys,
			Value:     encryptedValues,
		}, nil

	default:
		// For int and bool, no obfuscation needed
		return nil, nil
	}
}

// DetectFieldType automatically detects the type of a field based on its value
func DetectFieldType(value string) FieldType {
	// Try to parse as bool first (since "1" and "0" are valid bools)
	if _, err := strconv.ParseBool(value); err == nil {
		return FieldTypeBool
	}

	// Try to parse as int
	if _, err := strconv.Atoi(value); err == nil {
		return FieldTypeInt
	}

	// Try to parse as float
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return FieldTypeFloat
	}

	// Default to string
	return FieldTypeString
}

// extractFieldsFromEnvVars extracts fields from environment variables
func extractFieldsFromEnvVars(envVars map[string]string) []Field {
	var fields []Field

	for envName, value := range envVars {
		var fieldType FieldType
		if value == "" {
			fieldType = FieldTypeString // Empty values are treated as strings
		} else {
			fieldType = DetectFieldType(value)
		}

		fields = append(fields, Field{
			EnvName: envName,
			Type:    fieldType,
			Value:   value,
		})
	}

	return fields
}

// checkEnvironmentConsistency checks if all environments have the same variables
func checkEnvironmentConsistency(allEnvVars map[string]map[string]string) error {
	if len(allEnvVars) < 2 {
		return nil // No need to check consistency with only one environment
	}

	// Get all variable names from all environments
	allVars := make(map[string]bool)
	for _, envVars := range allEnvVars {
		for varName := range envVars {
			allVars[varName] = true
		}
	}

	// Check that each environment has all variables
	for envName, envVars := range allEnvVars {
		for varName := range allVars {
			if _, exists := envVars[varName]; !exists {
				return fmt.Errorf("❌ ERROR: variable '%s' is missing in environment '%s'", varName, envName)
			}
		}
	}

	fmt.Println("✅ Environment consistency check passed - all environments have the same variables")
	return nil
}

// LoadEnvFile loads environment variables from a .env file and returns Field slice
func LoadEnvFile(filePath string) ([]Field, error) {
	envVars, err := ReadEnvFile(filePath)
	if err != nil {
		return nil, err
	}

	return extractFieldsFromEnvVars(envVars), nil
}

// ReadEnvFile reads environment variables from a file
func ReadEnvFile(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	envVars := make(map[string]string)

	// Simple line-by-line reading
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			envVars[parts[0]] = parts[1]
		}
	}

	return envVars, nil
}

func NewGenerator(config *Config) *Generator {
	return &Generator{
		config: config,
	}
}

// GenerateFromEnvFile reads environment variables from a .env file and generates configuration
func (g *Generator) GenerateFromEnvFile(envFilePath string) error {
	envVars, err := ReadEnvFile(envFilePath)
	if err != nil {
		return fmt.Errorf("failed to read env file %s: %w", envFilePath, err)
	}

	// Extract fields from environment variables
	g.config.Fields = extractFieldsFromEnvVars(envVars)

	return g.generateConfigFile()
}

// LoadConfigFile loads configuration from JSON file
func LoadConfigFile(configFilePath string) (*ConfigFile, error) {
	// Read configuration file
	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configFilePath, err)
	}

	var configFile ConfigFile
	err = json.Unmarshal(configData, &configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configFilePath, err)
	}

	return &configFile, nil
}

// GenerateFromConfigFile generates configurations from JSON file
func GenerateFromConfigFile(configFilePath string) error {
	configFile, err := LoadConfigFile(configFilePath)
	if err != nil {
		return err
	}

	// Collect all environment variables from all environments for consistency check
	allEnvVars := make(map[string]map[string]string)
	for envName, envConfig := range configFile.Environments {
		envVars, err := ReadEnvFile(envConfig.EnvFile)
		if err != nil {
			return fmt.Errorf("failed to read env file %s: %w", envConfig.EnvFile, err)
		}
		allEnvVars[envName] = envVars
	}

	// Check consistency between environments
	if err := checkEnvironmentConsistency(allEnvVars); err != nil {
		return fmt.Errorf("environment consistency check failed: %w", err)
	}

	// Generate single merged configuration file
	fmt.Println("🔄 Generating merged configuration file...")

	// Prepare data for merged template
	mergedData := struct {
		PackageName  string
		RandomSeed   int64
		Environments map[string]struct {
			StructName string
			Fields     []Field
			Obfuscated map[string]*ObfuscationResult
		}
		AllFields []Field
	}{
		PackageName: configFile.PackageName,
		RandomSeed:  int64(configFile.RandomSeed),
		Environments: make(map[string]struct {
			StructName string
			Fields     []Field
			Obfuscated map[string]*ObfuscationResult
		}),
		AllFields: extractFieldsFromEnvVars(allEnvVars["dev"]), // Use dev as reference for interface
	}

	// Prepare fields for each environment
	for envName, envConfig := range configFile.Environments {
		envVars := allEnvVars[envName]
		fields := extractFieldsFromEnvVars(envVars)
		obfuscated := make(map[string]*ObfuscationResult)

		// Generate obfuscated data for each field
		for _, field := range fields {
			if field.Value != "" {
				result, err := generateObfuscatedField(field.EnvName, field.Type, field.Value, mergedData.RandomSeed)
				if err != nil {
					return fmt.Errorf("failed to obfuscate field %s: %w", field.EnvName, err)
				}
				obfuscated[field.EnvName] = result
			}
		}

		mergedData.Environments[envName] = struct {
			StructName string
			Fields     []Field
			Obfuscated map[string]*ObfuscationResult
		}{
			StructName: envConfig.StructName,
			Fields:     fields,
			Obfuscated: obfuscated,
		}
	}

	// Generate merged file
	outputFile := filepath.Join(configFile.OutputDir, "config_env.gen.go")
	err = generateMergedFile(outputFile, mergedData)
	if err != nil {
		return fmt.Errorf("failed to generate merged configuration: %w", err)
	}
	fmt.Println("✅ Merged configuration file generated successfully!")

	fmt.Println("\n🎉 All configurations generated!")
	fmt.Printf("📁 Files are located in %s\n", configFile.OutputDir)
	fmt.Println("🔧 You can now use the generated configurations directly")

	return nil
}

// AutoGenerate automatically generates configurations
// Searches for configuration file in current directory and parent directories
func AutoGenerate() error {
	configFile := findConfigFile()
	if configFile == "" {
		return fmt.Errorf("configuration file go-envied-config.json not found")
	}

	fmt.Printf("🔧 Automatic configuration generation from file: %s\n", configFile)
	return GenerateFromConfigFile(configFile)
}

// findConfigFile searches for configuration file in current directory and parent directories
func findConfigFile() string {
	configFileName := "go-envied-config.json"

	// Check current directory
	if _, err := os.Stat(configFileName); err == nil {
		return configFileName
	}

	// Check parent directories (maximum 3 levels up)
	currentDir, _ := os.Getwd()
	for i := 0; i < 3; i++ {
		parentPath := filepath.Join(currentDir, strings.Repeat("../", i+1), configFileName)
		if _, err := os.Stat(parentPath); err == nil {
			return parentPath
		}
	}

	return ""
}

// Init automatically generates configurations when package is imported
func Init() {
	err := AutoGenerate()
	if err != nil {
		fmt.Printf("⚠️ Warning: failed to generate configurations: %v\n", err)
		fmt.Println("💡 Make sure go-envied-config.json file exists in the project root")
	}
}

// GenerateFromEnvVars generates configuration from environment variables with strict validation
func (g *Generator) GenerateFromEnvVars() error {
	for i, field := range g.config.Fields {
		if value := os.Getenv(field.EnvName); value != "" {
			g.config.Fields[i].Value = value
		} else if os.Getenv(field.EnvName) == "" {
			// Check if variable exists but is empty
			if _, exists := os.LookupEnv(field.EnvName); exists {
				return fmt.Errorf("❌ ERROR: environment variable '%s' is empty", field.EnvName)
			}
		} else if field.DefaultValue != "" {
			// Only use default value if explicitly provided
			g.config.Fields[i].Value = field.DefaultValue
		} else if !field.Optional {
			return fmt.Errorf("❌ ERROR: required environment variable '%s' not found", field.EnvName)
		}
	}

	return g.generateConfigFile()
}

// generateConfigFile generates the Go configuration file
func (g *Generator) generateConfigFile() error {
	// Extract environment name from Environment (e.g., "DevConfig" -> "dev")
	envName := strings.ToLower(g.config.Environment)
	envName = strings.TrimSuffix(envName, "config")
	outputFile := filepath.Join(g.config.OutputDir, fmt.Sprintf("config_%s.go", envName))

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(g.config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Obfuscate all string fields before generating the file
	for i, field := range g.config.Fields {
		if field.Type == FieldTypeString && field.Value != "" {
			obfuscatedValue := Obfuscate(field.Value, "go-envied-obfuscation")
			g.config.Fields[i].Value = obfuscatedValue
		}
	}

	// Generate configuration file
	return g.generateFile(outputFile, configTemplate)
}

// generateFile generates a file from template
func (g *Generator) generateFile(outputFile string, templateStr string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	tmpl, err := template.New("config").Parse(templateStr)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	return tmpl.Execute(file, g.config)
}

// generateMergedFile generates a single merged configuration file
func generateMergedFile(outputFile string, data interface{}) error {
	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Generate code directly instead of using template
	return generateCodeDirectly(file, data)
}

// generateCodeDirectly generates the Go code directly
func generateCodeDirectly(file *os.File, data interface{}) error {
	// Type assertion to get the data
	mergedData, ok := data.(struct {
		PackageName  string
		RandomSeed   int64
		Environments map[string]struct {
			StructName string
			Fields     []Field
			Obfuscated map[string]*ObfuscationResult
		}
		AllFields []Field
	})
	if !ok {
		return fmt.Errorf("invalid data type for code generation")
	}

	// Write package header
	fmt.Fprintf(file, "// Code generated by go-envied. DO NOT EDIT.\n")
	fmt.Fprintf(file, "// Generated merged configuration file for all environments\n\n")
	fmt.Fprintf(file, "package %s\n\n", mergedData.PackageName)
	fmt.Fprintf(file, "import \"github.com/petrovyuri/go-envied\"\n\n")

	// Write interface
	fmt.Fprintf(file, "// ConfigInterface defines the interface for all generated configurations\n")
	fmt.Fprintf(file, "type ConfigInterface interface {\n")
	for _, field := range mergedData.AllFields {
		fmt.Fprintf(file, "\tGet%s() %s\n", field.EnvName, field.Type)
	}
	fmt.Fprintf(file, "}\n\n")

	// Write each environment
	for envName, envData := range mergedData.Environments {
		// Write static constants for keys and values with environment prefix
		for fieldName, obfuscated := range envData.Obfuscated {
			if obfuscated == nil {
				continue // Skip fields that don't need obfuscation
			}
			// Write key constant with environment prefix
			keyConstName := fmt.Sprintf("%s_%s", strings.ToUpper(envName), obfuscated.KeyName)
			fmt.Fprintf(file, "// Static key for %s in %s environment\n", fieldName, envName)
			fmt.Fprintf(file, "var %s = ", keyConstName)

			switch key := obfuscated.Key.(type) {
			case []int:
				fmt.Fprintf(file, "[]int{")
				for i, v := range key {
					if i > 0 {
						fmt.Fprintf(file, ", ")
					}
					fmt.Fprintf(file, "%d", v)
				}
				fmt.Fprintf(file, "}\n\n")
			case bool:
				fmt.Fprintf(file, "%t\n\n", key)
			case int:
				fmt.Fprintf(file, "%d\n\n", key)
			default:
				fmt.Fprintf(file, "%v\n\n", key)
			}

			// Write value constant if different from field name
			if obfuscated.ValueName != fieldName {
				valueConstName := fmt.Sprintf("%s_%s", strings.ToUpper(envName), obfuscated.ValueName)
				fmt.Fprintf(file, "// Static encrypted data for %s in %s environment\n", fieldName, envName)
				fmt.Fprintf(file, "var %s = []int{", valueConstName)

				switch value := obfuscated.Value.(type) {
				case []int:
					for i, v := range value {
						if i > 0 {
							fmt.Fprintf(file, ", ")
						}
						fmt.Fprintf(file, "%d", v)
					}
				default:
					fmt.Fprintf(file, "%v", value)
				}
				fmt.Fprintf(file, "}\n\n")
			}
		}

		// Write struct
		fmt.Fprintf(file, "// %sConfig - generated configuration for %s environment\n", envData.StructName, envName)
		fmt.Fprintf(file, "type %sConfig struct {\n", envData.StructName)
		for _, field := range envData.Fields {
			fmt.Fprintf(file, "\t%s %s\n", field.EnvName, field.Type)
		}
		fmt.Fprintf(file, "}\n\n")

		// Write constructor
		fmt.Fprintf(file, "// New%sConfig creates a new configuration for %s environment\n", envData.StructName, envName)
		fmt.Fprintf(file, "func New%sConfig() *%sConfig {\n", envData.StructName, envData.StructName)
		fmt.Fprintf(file, "\treturn &%sConfig{\n", envData.StructName)

		for _, field := range envData.Fields {
			if obfuscated, exists := envData.Obfuscated[field.EnvName]; exists {
				envPrefix := strings.ToUpper(envName)
				switch field.Type {
				case FieldTypeString:
					keyConstName := fmt.Sprintf("%s_%s", envPrefix, obfuscated.KeyName)
					valueConstName := fmt.Sprintf("%s_%s", envPrefix, obfuscated.ValueName)
					fmt.Fprintf(file, "\t\t%s: envied.DeobfuscateString(%s, %s),\n", field.EnvName, keyConstName, valueConstName)
				case FieldTypeFloat:
					keyConstName := fmt.Sprintf("%s_%s", envPrefix, obfuscated.KeyName)
					valueConstName := fmt.Sprintf("%s_%s", envPrefix, obfuscated.ValueName)
					fmt.Fprintf(file, "\t\t%s: envied.ParseFloat(envied.DeobfuscateString(%s, %s)),\n", field.EnvName, keyConstName, valueConstName)
				case FieldTypeInt:
					fmt.Fprintf(file, "\t\t%s: envied.ParseInt(\"%s\"),\n", field.EnvName, field.Value)
				case FieldTypeBool:
					fmt.Fprintf(file, "\t\t%s: envied.ParseBool(\"%s\"),\n", field.EnvName, field.Value)
				default:
					fmt.Fprintf(file, "\t\t%s: \"%s\",\n", field.EnvName, field.Value)
				}
			} else {
				// For int and bool, use simple parsing functions
				switch field.Type {
				case FieldTypeInt:
					fmt.Fprintf(file, "\t\t%s: envied.ParseInt(\"%s\"),\n", field.EnvName, field.Value)
				case FieldTypeBool:
					fmt.Fprintf(file, "\t\t%s: envied.ParseBool(\"%s\"),\n", field.EnvName, field.Value)
				default:
					fmt.Fprintf(file, "\t\t%s: \"%s\",\n", field.EnvName, field.Value)
				}
			}
		}
		fmt.Fprintf(file, "\t}\n")
		fmt.Fprintf(file, "}\n\n")

		// Write getter methods
		fmt.Fprintf(file, "// Getter methods for %sConfig\n", envData.StructName)
		for _, field := range envData.Fields {
			fmt.Fprintf(file, "func (c *%sConfig) Get%s() %s {\n", envData.StructName, field.EnvName, field.Type)
			fmt.Fprintf(file, "\treturn c.%s\n", field.EnvName)
			fmt.Fprintf(file, "}\n\n")
		}
	}

	return nil
}

// Template for generated configuration file
const configTemplate = `// Code generated by go-envied. DO NOT EDIT.
// Generated for {{.Environment}} environment

package {{.PackageName}}

import "github.com/petrovyuri/go-envied"

// {{.Environment}}Config - generated configuration for {{.Environment}} environment
type {{.Environment}}Config struct {
{{range .Fields}}	{{.EnvName}} {{.Type}}
{{end}}}

// New{{.Environment}}Config creates a new configuration for {{.Environment}} environment
func New{{.Environment}}Config() *{{.Environment}}Config {
	return &{{.Environment}}Config{
{{range .Fields}}{{if eq .Type "string"}}		{{.EnvName}}: envied.Deobfuscate("{{.Value}}"),
{{else if eq .Type "int"}}		{{.EnvName}}: envied.ParseInt("{{.Value}}"),
{{else if eq .Type "bool"}}		{{.EnvName}}: envied.ParseBool("{{.Value}}"),
{{else if eq .Type "float64"}}		{{.EnvName}}: envied.ParseFloat("{{.Value}}"),
{{else}}		{{.EnvName}}: "{{.Value}}",
{{end}}{{end}}	}
}

// Getter methods
{{range .Fields}}func (c *{{$.Environment}}Config) Get{{.EnvName}}() {{.Type}} {
	return c.{{.EnvName}}
}

{{end}}`
