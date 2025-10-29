# go-envied Usage Example

This example demonstrates how to use the `go-envied` package to generate type-safe configurations from environment variables.

## Project Structure

```
example/
├── go-envied-config.json    # JSON configuration
├── env/
│   ├── dev.env             # Development environment variables
│   └── prod.env            # Production environment variables
├── cmd/
│   └── generate/
│       └── main.go         # Configuration generator
├── generated/              # Generated files (created automatically)
├── main.go                 # Usage example
└── go.mod                  # Go module
```

## How to Run the Example

### 1. Generate Configurations

```bash
cd example
go run cmd/generate/main.go
```

This command will create the `generated/config_env.gen.go` file with type-safe configurations.

### 2. Run the Example

```bash
go run main.go
```

### 3. Run with Different Environments

```bash
# For development (default)
APP_ENV=dev go run main.go

# For production
APP_ENV=prod go run main.go
```

## What Happens

1. **Generation**: The package reads JSON configuration and `.env` files
2. **Validation**: Consistency between environments is checked
3. **Encryption**: String values are encrypted for security
4. **Code Generation**: Go code with type-safe structures is created
5. **Usage**: Ready-to-use configurations can be used in code

## Benefits

- ✅ **Type Safety**: All variables have correct types
- ✅ **Security**: String values are encrypted
- ✅ **Polymorphism**: Single interface for all environments
- ✅ **Validation**: Errors are detected during generation
- ✅ **Performance**: No file reading during runtime
