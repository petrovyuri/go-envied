package config

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
