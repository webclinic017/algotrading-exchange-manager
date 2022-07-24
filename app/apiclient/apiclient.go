package apiclient

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/srv"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/asmcos/requests"
)

// todo: add multi symbol support in response? current hardcoded to first signal
func SignalAnalyzer(ob *appdata.OrderBook_S, mode string) bool {

	p := requests.Params{
		"mode":           mode,
		"algo":           ob.Strategy,
		"symbol":         ob.Instr,
		"date":           time.Now().Format("2006-01-02"),
		"pos_dir":        ob.Dir,
		"pos_entr_price": fmt.Sprintf("%f", ob.Targets.EntrPrice),
		"pos_entr_time":  ob.Date.Format(time.RFC3339),
	}
	resp, err := requests.Get(appdata.Env["ALGO_ANALYSIS_ADDRESS"]+"tradesignals/", p)

	if err != nil {
		// srv.WarningLogger.Println(ob.Instr, "-", ob.Strategy, "]", err.Error())
		return false
	}

	if resp.R.StatusCode == http.StatusOK {

		var apiSig []appdata.ApiSignal_S

		// fmt.Println(resp.Text())
		err := json.Unmarshal([]byte(resp.Text()), &apiSig)
		if err != nil {
			srv.WarningLogger.Println(ob.Instr, "-", ob.Strategy, "]", err.Error())
			return false
		}
		// check if signal processed for the same as requested
		if len(apiSig) > 0 {
			if apiSig[0].Status == "signal-processed" &&
				apiSig[0].Instr == ob.Instr &&
				apiSig[0].Strategy == ob.Strategy { // register only if processed correctly

				ob.Dir = apiSig[0].Dir
				if mode == "-entr" {
					ob.Targets.EntrPrice = apiSig[0].TriggerValue
					ob.Targets.EntrTime = apiSig[0]

				} else if mode == "-exit" {
					ob.Exit_reason = apiSig[0].ExitReason
					ob.Targets.ExitPrice = apiSig[0].TriggerValue
				}

				/*ob.Entry = apiSig[0].Entry
				ob.Stoploss = apiSig[0].Stoploss
				ob.Target = apiSig[0].Target*/
				return true
			} else {
				srv.WarningLogger.Println(ob.Instr, "-", ob.Strategy, "]", apiSig[0].Status)
				return false
			}
		}
	}
	srv.WarningLogger.Println(ob.Instr, "-", ob.Strategy, "]", resp.R.StatusCode)
	return false
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
