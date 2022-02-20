package apiclient

import (
	"fmt"

	"github.com/asmcos/requests"
)

func ExecuteApi() {

	p := requests.Params{
		"algo":   "S001-ORB-001",
		"symbol": "BANKNIFTY",
		"date":   "2022-02-09",
	}
	resp, err := requests.Get("https://algoanalysis.wyealth.com/tradesignals/", p)
	// resp, err := requests.Get("https://jsonplaceholder.typicode.com/todos/1")
	if err != nil {
		return
	}
	fmt.Println(resp.Text())

	var json map[string]interface{}
	resp.Json(&json)

	for k, v := range json {
		fmt.Println(k, v)
	}
}
