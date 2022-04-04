package srv

import (
	"algo-ex-mgr/app/appdata"
	"os"

	"github.com/joho/godotenv"
)

var (
	Env = make(map[string]string)
)

func LoadEnvVariables(path string) bool {

	// Load .env file, if not in production
	var parseEnv = true

	err := godotenv.Load(path)
	if err != nil {
		WarningLogger.Println("Error loading .env file", err)
	}
	// Load and check values
	for _, value := range appdata.UserSettings {
		Env[value] = os.Getenv(value)
		if os.Getenv(value) == "" {
			println(value, " is not set")
			parseEnv = false
		}

	}
	return parseEnv
}
