package main

import (
	"flag"
	"fmt"
	"log"

	// Import generated configurations
	configPkg "example/internal/config"
)

func main() {
	fmt.Println("🎯 go-envied Usage Example")
	fmt.Println("=====================================")

	// Parse command line flags
	env, err := getEnv()
	if err != nil {
		log.Fatalf("Failed to get environment: %v", err)
	}

	// Create configuration for the selected environment
	cfg, err := configPkg.NewConfig(env)
	if err != nil {
		log.Fatalf("Failed to create config: %v", err)
	}

	fmt.Printf("\n🔧 Using %s Environment Configuration:\n", env)
	fmt.Printf("  DATABASE_URL: %s\n", cfg.DATABASE_URL)
	fmt.Printf("  DEBUG_MODE: %t\n", cfg.DEBUG_MODE)
	fmt.Printf("  PORT: %d\n", cfg.PORT)
	fmt.Printf("  TEMPERATURE: %f\n", cfg.TEMPERATURE)
	fmt.Printf("  MAX_TOKENS: %s\n", cfg.MAX_TOKENS)
}

func getEnv() (string, error) {
	flag.Parse()

	// Define environment
	var env string

	// Check positional arguments (for compatibility with go run main.go dev)
	if len(flag.Args()) > 0 {
		switch flag.Args()[0] {
		case configPkg.EnvDev:
			env = configPkg.EnvDev
		default:
			env = configPkg.EnvProd
		}
	} else {
		// Default to production environment
		env = configPkg.EnvProd
	}
	fmt.Printf("Using environment: %s\n", env)
	return env, nil
}
