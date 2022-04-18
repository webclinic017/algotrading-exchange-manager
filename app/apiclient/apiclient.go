package apiclient

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/srv"
	"encoding/json"
	"net/http"
	"time"

	"github.com/asmcos/requests"
)

func SignalAnalyzer(tr *appdata.TradeSignal, mode string) bool {

	p := requests.Params{
		"multisymbol": "false",
		"algo":        tr.Strategy + mode,
		"symbol":      tr.Instr,
		"date":        time.Now().Format("2006-01-02"),
	}
	resp, err := requests.Get(appdata.Env["ALGO_ANALYSIS_ADDRESS"]+"tradesignals/", p)
	// resp, err := requests.Get("http://localhost:5000/tradesignals/", p)

	println(resp, resp)

	if err != nil {
		srv.WarningLogger.Println(tr.Instr, "-", tr.Strategy, "]", err.Error())
		return false
	}

	if resp.R.StatusCode == http.StatusOK {

		apiSig := appdata.ApiSignal{}

		json.Unmarshal([]byte(resp.Text()), &apiSig)
		if apiSig.Status == "signal-processed" { // register only if processed correctly
			tr.Dir = apiSig.Dir
			tr.Entry = apiSig.Entry
			tr.Stoploss = apiSig.Stoploss
			return true
		} else {
			srv.WarningLogger.Println(tr.Instr, "-", tr.Strategy, "]", apiSig.Status)
			return false
		}
	} else {

		srv.WarningLogger.Println(tr.Instr, "-", tr.Strategy, "]", resp.R.StatusCode)
		return false
	}
}

func Services(service string, date time.Time) bool {

	p := requests.Params{
		"sid":  service,
		"date": date.Format("2006-01-02"),
	}
	req := requests.Requests()
	req.SetTimeout(120) // set timeout to 120 seconds, candle service requires time to respond
	_, err := req.Get(appdata.Env["ALGO_ANALYSIS_ADDRESS"]+"services/", p)

	return err == nil
}
