package apiclient

import (
	"algo-ex-mgr/app/appdata"
	"encoding/json"
	"time"

	"github.com/asmcos/requests"
)

func SignalAnalyzer(multiSymbol string, algo string, symbol string, date string) (bool, string) {

	p := requests.Params{
		"multisymbol": multiSymbol,
		"algo":        algo,
		"symbol":      symbol,
		"date":        date,
	}
	resp, err := requests.Get(appdata.Env["ALGO_ANALYSIS_ADDRESS"]+"tradesignals/", p)
	// resp, err := requests.Get("http://localhost:5000/tradesignals/", p)

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

func Services(service string, date time.Time) bool {

	p := requests.Params{
		"sid":  service,
		"date": date.Format("2006-01-02"),
	}
	_, err := requests.Get(appdata.Env["ALGO_ANALYSIS_ADDRESS"]+"services/", p)

	if err != nil {
		return false
	}
	return true
}
