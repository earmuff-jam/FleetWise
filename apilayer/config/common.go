package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/logutils"
	"github.com/joho/godotenv"
)

// InitLogger ...
//
// initialize the logger based on the environment variable from config
func InitLogger() {

	logMode := os.Getenv("DEBUG")

	if len(logMode) <= 0 {
		logMode = "WARN"
	}

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(logMode),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
}

// Log ...
//
// function used to print details in debug mode. pass in error message to print
// error scenarios and log them. Replace it with nil to not use the error messages.
func Log(logMsg string, errMsg error, args ...interface{}) {
	if errMsg != nil {
		log.Printf("[DEBUG] "+logMsg+". error: %+v", append(args, errMsg)...)
	} else {
		log.Printf("[DEBUG] "+logMsg, args...)
	}
}

// PreloadAllTestVariables ...
//
// load environment variables
func PreloadAllTestVariables() {
	err := godotenv.Load(filepath.Join("..", "..", ".env"))
	if err != nil {
		log.Printf("No env file detected. Using os system configuration.")
	}
}
