package apiclient

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/srv"
	"encoding/json"
	"net/http"
	"time"

	"github.com/asmcos/requests"
)

// todo: add multi symbol support in response? current hardcoded to first signal
func SignalAnalyzer(tr *appdata.OrderBook_S, mode string) bool {

	p := requests.Params{
		"multisymbol": "false",
		"algo":        tr.Strategy + mode,
		"symbol":      tr.Instr,
		"date":        time.Now().Format("2006-01-02"),
	}
	resp, err := requests.Get(appdata.Env["ALGO_ANALYSIS_ADDRESS"]+"tradesignals/", p)

	if err != nil {
		// srv.WarningLogger.Println(tr.Instr, "-", tr.Strategy, "]", err.Error())
		return false
	}

	if resp.R.StatusCode == http.StatusOK {

		var apiSig []appdata.ApiSignal

		// fmt.Println(resp.Text())
		err := json.Unmarshal([]byte(resp.Text()), &apiSig)
		if err != nil {
			srv.WarningLogger.Println(tr.Instr, "-", tr.Strategy, "]", err.Error())
			return false
		}
		// check if signal processed for the same as requested
		if apiSig[0].Status == "signal-processed" &&
			apiSig[0].Instr == tr.Instr &&
			apiSig[0].Strategy == tr.Strategy { // register only if processed correctly

			tr.Dir = apiSig[0].Dir
			/*tr.Entry = apiSig[0].Entry
			tr.Stoploss = apiSig[0].Stoploss
			tr.Target = apiSig[0].Target*/
			return true
		} else {
			srv.WarningLogger.Println(tr.Instr, "-", tr.Strategy, "]", apiSig[0].Status)
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
