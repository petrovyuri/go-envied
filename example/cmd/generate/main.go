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
