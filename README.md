# go-envied

Type-safe configuration generator for Go, inspired by the Flutter [envied](https://github.com/petercinibulk/envied) package.

## 🚀 Features

- 🌍 **Multiple Environments**: Support for dev, prod, staging and other environments
- 📁 **Automatic Scanning**: All variables from .env files are automatically included
- ⚡ **Speed**: No file reading during runtime - everything is compiled
- 🛡️ **Strict Validation**: Generation stops with error if variable is not found or empty
- 🔧 **JSON Configuration**: Support for JSON configuration files
- 🚀 **CLI Utility**: Convenient command line for generating configurations
- 📊 **Parsing Utilities**: Built-in functions for parsing various data types
- 🎨 **Single File**: All configurations and interface in one file `config_env.gen.go`
- 🔗 **Universal Interface**: Automatically generated `ConfigInterface` for all configurations
- ✅ **Consistency Check**: All environments must have the same variables
- 🛠️ **Cross-platform**: Support for Linux, macOS and Windows

## 📦 Installation

### Using as Library

```bash
go get github.com/petrovyuri/go-envied
```

## 🚀 Quick Start

### JSON Configuration

#### 1. Create JSON Configuration File (`go-envied-config.json`)

```json
{
  "package_name": "config",
  "output_dir": "path/to/your/package",
  "environments": [
    {
      "env_file": "path/to/your/env/dev.env",
      "struct_name": "DevConfig"
    },
    {
      "env_file": "path/to/your/env/prod.env",
      "struct_name": "ProdConfig"
    }
  ]
}
```

#### 2. Create Environment Files

**env/dev.env:**
```env
# Development environment configuration

# Example of String parsing and obfuscation
DATABASE_URL=dev-database-url
# Example of Bool parsing
DEBUG_MODE=true
# Example of integer parsing
PORT=10000
# Example of float parsing
TEMPERATURE=0.1
# Although this is a number, it can be wrapped in quotes,
# then it will be treated as a string
# and it will be obfuscated
MAX_TOKENS="10" 
```

**env/prod.env:**
```env
# Production environment configuration
# Example of String parsing and obfuscation
DATABASE_URL=prod-database-url
# Example of Bool parsing
DEBUG_MODE=false
# Example of integer parsing
PORT=80
# Example of float parsing
TEMPERATURE=0.8
# Although this is a number, it can be wrapped in quotes,
# then it will be treated as a string
# and it will be obfuscated
MAX_TOKENS="1000" 
```

#### 3. Run Generation

```bash
# Create directory cmd/generate
mkdir -p cmd/generate

# Create main.go file in this directory
touch cmd/generate/main.go
```

In this file add the following code:
```go
package main

// This file is used to generate the configurations

import (
	"log"

	"github.com/petrovyuri/go-envied"
)

func main() {
	log.Printf("🚀 Generating configurations with go-envied...")

	// Automatic generation from JSON configuration
	err := envied.AutoGenerate()
	if err != nil {
		log.Fatalf("❌ Configuration generation error: %v", err)
	}

	log.Printf("✅ Configurations generated successfully!")
	log.Printf("📁 Files are located in ./generated directory")
}

```

```bash
# Run generation
go run cmd/generate/main.go
``` 

In the `path/to/your/package` directory the file `config_env.gen.go` will be generated

#### 4. Use Generated Configurations

```go
package path/to/your/package

import (
	"fmt"
)

const (
	EnvDev  = "dev"
	EnvProd = "prod"
)

type Config struct {
	DATABASE_URL string
	DEBUG_MODE   bool
	PORT         int
	TEMPERATURE  float64
	MAX_TOKENS   string
}

func NewConfig(env string) (*Config, error) {
	// Create configurations for different environments
	var currentConfig ConfigInterface
	switch env {
	case EnvDev:
		currentConfig = NewDevConfigConfig()
		fmt.Println("  Using development configuration")
	default:
		currentConfig = NewProdConfigConfig()
		fmt.Printf("  Unknown environment '%s', using development configuration\n", env)
	}

	return &Config{
		DATABASE_URL: currentConfig.GetDATABASE_URL(),
		DEBUG_MODE:   currentConfig.GetDEBUG_MODE(),
		PORT:         currentConfig.GetPORT(),
		TEMPERATURE:  currentConfig.GetTEMPERATURE(),
		MAX_TOKENS:   currentConfig.GetMAX_TOKENS(),
	}, nil
}
```

## 📊 Field Types

- `string` - string values
- `int` - integers
- `bool` - boolean values (true/false)
- `float64` - floating point numbers

## ⚙️ Field Options

- **Automatic Type Detection**: System automatically detects type based on value
- **Strict Validation**: All fields are required and cannot be empty
- **Consistency Check**: All environments must have the same variables

## 🎯 go-envied Advantages

### Compared to Regular Environment Variables:

- 🚀 **Zero Configuration** - just add variable to `.env`
- 🎨 **Single File** - all configurations in one place
- 🔄 **Automatic Synchronization** - interface is always up to date
- 🛡️ **Fail Fast** - problems are visible immediately during generation

### Compared to Other Solutions:

- 📊 **Automatic Type Detection** - no need to specify types manually
- 🔗 **Universal Interface** - polymorphism out of the box
- ✅ **Consistency Check** - all environments are synchronized
- 🚀 **CLI Utility** - convenient command line

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes
4. Add tests
5. Submit pull request

## 📁 Example

A complete working example of go-envied usage can be found in the [`example/`](example/) folder.

The example includes:
- Configuration files for different environments (`dev.env`, `prod.env`)
- JSON configuration (`go-envied-config.json`)
- Configuration generator (`cmd/generate/main.go`)
- Main application using generated configurations (`main.go`)

To run the example:
```bash
cd example
go run main.go dev  
```

or

```bash
cd example
go run main.go prod # or any other environment
```
## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

Inspired by the [envied](https://pub.dev/packages/envied) package for Flutter.