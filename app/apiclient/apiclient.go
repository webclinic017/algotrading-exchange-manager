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
		"pos_entr_price": fmt.Sprintf("%f", ob.ApiSignalEntr.Entry),
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

				if mode == "entr" {
					ob.Dir = apiSig[0].Dir
					copyApiSignalvalues(&apiSig[0], &ob.ApiSignalEntr)

				} else if mode == "exit" {
					ob.Exit_reason = apiSig[0].ExitReason
					copyApiSignalvalues(&apiSig[0], &ob.ApiSignalExit)
				}
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

func copyApiSignalvalues(rec *appdata.ApiSignal_S, d *appdata.ApiSignal_S) {

	d.Status = rec.Status
	d.Id = rec.Id
	d.Date = rec.Date
	d.Instr = rec.Instr
	d.Strategy = rec.Strategy
	d.Dir = rec.Dir
	d.Entry = rec.Entry
	d.Target = rec.Target
	d.Stoploss = rec.Stoploss
	d.DebugEntr = rec.DebugEntr
	d.EntryTime = rec.EntryTime
	d.TriggerValue = rec.TriggerValue
	d.Exit = rec.Exit
	d.ExitTime = rec.ExitTime
	d.ExitReason = rec.ExitReason
	d.Debug = rec.Debug
	d.Gain = rec.Gain
	d.TimeDiff = rec.TimeDiff
	// rec.
}

func Services(service string, date time.Time) bool {

	p := requests.Params{
		"sid":   service,
		"date":  date.Format("2006-01-02"),
		"table": "all",
	}
	req := requests.Requests()
	req.SetTimeout(120) // set timeout to 120 seconds, candle service requires time to respond
	_, err := req.Get(appdata.Env["ALGO_ANALYSIS_ADDRESS"]+"services/", p)

	return err == nil
}
