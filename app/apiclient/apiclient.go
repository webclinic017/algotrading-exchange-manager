package apiclient

import (
	"encoding/json"

	"github.com/asmcos/requests"
)

func ExecuteSingleSymbolApi(algo string, symbol string, date string) (bool, string) {

	p := requests.Params{
		"multisymbol": "true",
		"algo":        algo,
		"symbol":      symbol,
		"date":        date,
	}
	// resp, err := requests.Get("https://algoanalysis.wyealth.com/tradesignals/", p)
	resp, err := requests.Get("http://localhost:5000/tradesignals/", p)
	// resp, err := requests.Get("https://jsonplaceholder.typicode.com/todos/1")
	if err != nil {
		return false, "nil"
	}

	var js interface{}
	json.Unmarshal([]byte(resp.Text()), &js)

	if len(js.([]interface{})) > 0 {
		return true, resp.Text()
	} else {
		return false, "nil"
	}
}
