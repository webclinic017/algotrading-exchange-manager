package srv

import (
	"algo-ex-mgr/app/appdata"
	"fmt"
	"os"

	"testing"
)

const ()

type LoadEnvVariablesTesting struct {
	path     string
	expected bool
}

// ** This is live testcase - update dates are per current symbols dates and levels.
// ** Result needs to be verified manually!!!
var LoadEnvVariablesTests = []LoadEnvVariablesTesting{
	{"/home/parag/devArea/algotrading-exchange-manager/app/zfiles/Unittest-Support-Files/envs/ userSettings-key-missing.env", false},
	{"/home/parag/devArea/algotrading-exchange-manager/app/zfiles/Unittest-Support-Files/envs/ userSettings-1-missing.env", false},
	{"/home/parag/devArea/algotrading-exchange-manager/app/zfiles/Unittest-Support-Files/envs/ userSettings-5-missing.env", false},
	{"/home/parag/devArea/algotrading-exchange-manager/app/zfiles/Unittest-Support-Files/envs/ userSettings-all-missing.env", false},
	{"/home/parag/devArea/algotrading-exchange-manager/app/zfiles/Unittest-Support-Files/envs/ userSettings-0-missing.env", true},
	{"", false},
}

func TestLoadEnvVariables(t *testing.T) {
	t.Parallel()

	InitLogger()

	for _, test := range LoadEnvVariablesTests {

		actual := LoadEnvVariables(test.path, false)
		if actual != test.expected {
			fmt.Println(" FAIL -  ", test.path, test.expected)
		}
		clearEnvVariables()

	}
	fmt.Println()
}

func clearEnvVariables() {
	for _, value := range appdata.UserSettings {
		os.Unsetenv(value)
	}
}
