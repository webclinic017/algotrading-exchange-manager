package srv

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() bool {

	// Load .env file, if not in production

	println("PRODUCTION - ", os.Getenv("PRODUCTION"))
	if os.Getenv("PRODUCTION") != "true" {
		WarningLogger.Println("DEVELOPMENT ENV")
		InfoLogger.Println("Ensure ENV variables are set in ENV_settings.env")
		FileCopyIfMissing("app/templates/ENV_Settings.env", "app/config/ENV_Settings.env")
		_ = godotenv.Load("app/config/ENV_Settings.env")
	} else {
		InfoLogger.Println("PRODUCTION ENV- Ensure ENV variables are set")
	}

	InfoLogger.Println("user ID", os.Getenv("USER_ID"))

	if 0 >= len(os.Getenv("LIVE_TRADING_MODE")) {
		ErrorLogger.Println("LIVE_TRADING_MODE not set")
		return false
	}

	if 0 >= len(os.Getenv("USER_ID")) {
		ErrorLogger.Println("USER_ID not set")
		return false
	}
	if 0 >= len(os.Getenv("TFA_AUTH")) {
		ErrorLogger.Println("TFA_AUTH not set")
		return false
	}
	if 0 >= len(os.Getenv("PASSWORD")) {

		ErrorLogger.Println("PASSWORD not set")
		return false
	}
	if 0 >= len(os.Getenv("API_KEY")) {
		ErrorLogger.Println("API_KEY not set")
		return false
	}
	if 0 >= len(os.Getenv("API_SECRET")) {
		ErrorLogger.Println("API_SECRET not set")
		return false
	}
	if 0 >= len(os.Getenv("TIMESCALEDB_ADDRESS")) {
		ErrorLogger.Println("TIMESCALEDB_ADDRESS not set")
		return false
	}
	if 0 >= len(os.Getenv("TIMESCALEDB_USERNAME")) {
		ErrorLogger.Println("TIMESCALEDB_USERNAME not set")
		return false
	}
	if 0 >= len(os.Getenv("TIMESCALEDB_PASSWORD")) {
		ErrorLogger.Println("TIMESCALEDB_PASSWORD not set")
		return false
	}
	if 0 >= len(os.Getenv("TIMESCALEDB_PORT")) {
		ErrorLogger.Println("TIMESCALEDB_PORT not set")
		return false
	}
	return true
}
