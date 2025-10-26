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
  "output_dir": "./internal/config",
  "environments": {
    "dev": {
      "env_file": "./env/dev.env",
      "struct_name": "DevConfig"
    },
    "prod": {
      "env_file": "./env/prod.env",
      "struct_name": "ProdConfig"
    }
  }
}
```

#### 2. Create Environment Files

**env/dev.env:**
```env
ENV1=your-dev-value-1
ENV2=123
ENV3=true
ENV4=https://example.com/api
ENV5=30
```

**env/prod.env:**
```env
ENV1=your-prod-value-1
ENV2=456
ENV3=false
ENV4=https://api.example.com
ENV5=7200
```

#### 3. Run Generation

```bash
# Create directory cmd/gen_env
mkdir -p cmd/gen_env

# Create main.go file in this directory
touch cmd/gen_env/main.go
```

In this file add the following code:
```go
package main

import "github.com/petrovyuri/go-envied"

func main() {
	envied.AutoGenerate()
}
```

```bash
# Run generation
go run cmd/gen_env/main.go
```

In the `internal/config` directory the file `config_env.gen.go` will be generated

#### 4. Use Generated Configurations

```go
package main

import (
    "fmt"
    "your-project/internal/config"
)

func main() {
    // Dev configuration
    devConfig := config.NewDevConfigConfig()
    fmt.Printf("Dev ENV1: %s\n", devConfig.GetENV1())
    fmt.Printf("Dev ENV2: %d\n", devConfig.GetENV2())
    
    // Prod configuration
    prodConfig := config.NewProdConfigConfig()
    fmt.Printf("Prod ENV1: %s\n", prodConfig.GetENV1())
    fmt.Printf("Prod ENV2: %d\n", prodConfig.GetENV2())
    
    // Polymorphism through interface
    configs := []config.ConfigInterface{
        devConfig,
        prodConfig,
    }
    
    for i, cfg := range configs {
        fmt.Printf("Config %d - ENV1: %s\n", i+1, cfg.GetENV1())
    }
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

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

Inspired by the [envied](https://pub.dev/packages/envied) package for Flutter.