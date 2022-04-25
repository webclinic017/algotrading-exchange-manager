package srv

import (
	"algo-ex-mgr/app/appdata"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables(path string, live bool) bool {

	// Load .env file, if not in production
	var parseEnv = true

	err := godotenv.Load(path)
	if err != nil {
		WarningLogger.Println("Error loading .env file", err)
	}
	// Load and check values
	for _, value := range appdata.UserSettings {
		appdata.Env[value] = os.Getenv(value)
		if os.Getenv(value) == "" && value != "DB_TEST_PREFIX" {
			println(value, " is not set")
			parseEnv = false
		}

	}

	if live {
		appdata.Env["DB_TEST_PREFIX"] = ""
	}

	appdata.Env["DB_TBL_USER_SETTING"] = appdata.Env["DB_TBL_PREFIX_USER_ID"] + appdata.Env["DB_TBL_USER_SETTING"] + appdata.Env["DB_TEST_PREFIX"]
	appdata.Env["DB_TBL_USER_SYMBOLS"] = appdata.Env["DB_TBL_PREFIX_USER_ID"] + appdata.Env["DB_TBL_USER_SYMBOLS"] + appdata.Env["DB_TEST_PREFIX"]
	appdata.Env["DB_TBL_USER_STRATEGIES"] = appdata.Env["DB_TBL_PREFIX_USER_ID"] + appdata.Env["DB_TBL_USER_STRATEGIES"] + appdata.Env["DB_TEST_PREFIX"]
	appdata.Env["DB_TBL_ORDER_BOOK"] = appdata.Env["DB_TBL_PREFIX_USER_ID"] + appdata.Env["DB_TBL_ORDER_BOOK"] + appdata.Env["DB_TEST_PREFIX"]
	appdata.Env["DB_TBL_TICK_NSEFUT"] = appdata.Env["DB_TBL_TICK_NSEFUT"] + appdata.Env["DB_TEST_PREFIX"]
	appdata.Env["DB_TBL_TICK_NSESTK"] = appdata.Env["DB_TBL_TICK_NSESTK"] + appdata.Env["DB_TEST_PREFIX"]

	appdata.Env["DB_URL"] = "postgres://" + appdata.Env["TIMESCALEDB_USERNAME"] + ":" + appdata.Env["TIMESCALEDB_PASSWORD"] + "@" + appdata.Env["TIMESCALEDB_ADDRESS"] + ":" + appdata.Env["TIMESCALEDB_PORT"]

	return parseEnv
}
