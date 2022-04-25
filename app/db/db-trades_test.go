package db

import (
	"algo-ex-mgr/app/srv"
	"fmt"
	"os"
	"testing"
)

func TestReadOrderBookFromDb(t *testing.T) {

	srv.Init()
	mydir, _ := os.Getwd()
	srv.LoadEnvVariables(mydir+"/../../userSettings.env", false)
	DbInit()

	a, b := ReadOrderBookFromDb(1)

	fmt.Println(a)
	fmt.Println(b)
}
