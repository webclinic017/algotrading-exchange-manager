package db

import (
	"algo-ex-mgr/app/srv"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestReadUserStrategiesFromDb(t *testing.T) {
	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	DbInit()

	a := ReadUserStrategiesFromDb()

	j, _ := json.Marshal(a)
	fmt.Println(string(j))
}
