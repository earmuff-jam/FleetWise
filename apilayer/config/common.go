package config

import (
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
)

// the auth key required to decode any request to mashed api
var BASE_LICENSE_KEY = []byte("2530d6a4-5d42-4758-b331-2fbbfed27bf9")

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
