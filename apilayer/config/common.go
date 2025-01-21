package config

import (
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
)

const CTO_USER = "community_test"
const CEO_USER = "ceo_test"

const DefaultTokenValidityTime = "15"

// PreloadAllTestVariables ...
//
// load environment variables
func PreloadAllTestVariables() {
	err := godotenv.Load(filepath.Join("..", "..", ".env"))
	if err != nil {
		log.Printf("No env file detected. Using os system configuration.")
	}
}
