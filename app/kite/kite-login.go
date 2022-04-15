package kite

import (
	"algo-ex-mgr/app/appdata"
	"algo-ex-mgr/app/srv"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/asmcos/requests"
	"github.com/joho/godotenv"
	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

var kc *kiteconnect.Client

const (
	URL      = "https://kite.zerodha.com/api/login"
	twofaUrl = "https://kite.zerodha.com/api/twofa"
)

func Init() bool {

	err := godotenv.Load("./zerodha-access-token.env")

	// ------------------------------------------------------------ Saved access token?, on failure login again
	if err != nil {
		srv.ErrorLogger.Println(err.Error())
		return loginKite()
	}

	// ------------------------------------------------------------ Set saved access token?, on failure login again
	if !setAccessToken(os.Getenv("kiteaccessToken")) {
		return loginKite()
	}
	return true
}

func loginKite() bool {

	srv.InfoLogger.Print(
		"\n\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~",
		"Zerodha Login",
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n")
	apiKey := appdata.Env["ZERODHA_API_KEY"]
	apiSecret := appdata.Env["ZERODHA_API_SECRET"]

	requestToken := KiteGetRequestToken()

	if strings.Contains(requestToken, "ERR:") {
		srv.ErrorLogger.Println("Authentication", requestToken)
	} else {

		srv.InfoLogger.Println("Authentication Succesful:", requestToken)
		// Create a new Kite connect instance
		kc = kiteconnect.New(apiKey)

		// Get user details and access token
		data, err := kc.GenerateSession(requestToken, apiSecret)
		if err != nil {
			srv.ErrorLogger.Printf("Session Err: %v", err)
			return false
		}

		// Set access token
		kc.SetAccessToken(data.AccessToken)
		srv.InfoLogger.Println("AccessToken", data.AccessToken)

		// keypair := strings.Join("accessToken", data.AccessToken)
		env, _ := godotenv.Unmarshal("kiteaccessToken=" + data.AccessToken)
		godotenv.Write(env, "./zerodha-access-token.env")
		os.Setenv("kiteaccessToken", data.AccessToken)

		// Get margins
		margins, err := kc.GetUserMargins()
		if err != nil {
			srv.ErrorLogger.Printf("Error getting margins: %v", err)
			//return false, "", ""
		}
		srv.InfoLogger.Println("Cash Balance (Net): ", margins.Equity.Net)
		return true
	}
	return false
}

func GetUserMargin() float64 {
	margins, err := kc.GetUserMargins()
	if err != nil {
		return 0
	}
	return margins.Equity.Net

}

func setAccessToken(accessToken string) bool {
	kc = kiteconnect.New(appdata.Env["ZERODHA_API_KEY"])
	kc.SetAccessToken(accessToken)
	margins, err := kc.GetUserMargins()
	if err != nil {
		srv.ErrorLogger.Printf("Error getting margins: %v", err)
		return false
	}
	srv.InfoLogger.Println("Cash Balance (Net): ", margins.Equity.Net)
	return true
}

func KiteGetRequestToken() string {

	defer func() {
		if err := recover(); err != nil {
			srv.ErrorLogger.Println("panic occurred:", err)
		}
	}()

	// ------------------------------------------------------------ 1. Start login, get reqId
	data := requests.Datas{
		"user_id":  appdata.Env["ZERODHA_USER_ID"],
		"password": appdata.Env["ZERODHA_PASSWORD"],
	}

	req := requests.Requests()
	resp, err := req.Post(URL, data)

	if (err != nil) || (resp.R.StatusCode != 200) {
		return "ERR: User ID / Password error"
	}

	// {"status":"success","data":{"user_id":"ZY7293",
	// "request_id":"mJHDflgxtNutiffR3jAZcIPA5SH4eXaLNkrwSgheQfoyb4syN9vJ5OJyGnPRKCi7",
	// "twofa_type":"pin","twofa_status":"active"}}
	reqID := extractValue(resp.Text(), "request_id")

	// ------------------------------------------------------------ 2. Do Two factor auth
	// If TOTP not enabled, use PIN for login
	twoAuth := appdata.Env["ZERODHA_TOTP_SECRET_KEY"]
	if twoAuth != "NOT-USED" {
		twoAuth = srv.GetTOTPToken(appdata.Env["ZERODHA_TOTP_SECRET_KEY"])
	} else {
		twoAuth = appdata.Env["ZERODHA_PIN"]
	}

	data = requests.Datas{
		"user_id":     appdata.Env["ZERODHA_USER_ID"],
		"request_id":  reqID,
		"twofa_value": twoAuth,
		// RULE - TWO FACTOR AUTH - MUST BE ENABLED
	}
	resp, err = req.Post(twofaUrl, data)

	if (err != nil) || (resp.R.StatusCode != 200) {
		return "ERR: Two factor auth failed"
	}

	// ------------------------------------------------------------ 3. Post login, access URL to get requestToken
	req.SetTimeout(5)
	resp, err = req.Get(appdata.Env["ZERODHA_REQ_TOKEN_URL"] + appdata.Env["ZERODHA_API_KEY"])
	if err != nil {
		srv.WarningLogger.Println(err.Error())
		arr := strings.Split(err.Error(), `"`)          // split on '&'
		return extractKeyValue(arr[1], "request_token") //
	} else {
		m, err := url.ParseQuery(resp.R.Request.URL.RawQuery)
		if (err != nil) || (resp.R.StatusCode != 200) {
			return "ERR: Cannot fetch request token"
		}
		return m["request_token"][0]
	}
}

func extractValue(body string, key string) string {
	keystr := "\"" + key + "\":[^,;\\]}]*"
	r, _ := regexp.Compile(keystr)
	match := r.FindString(body)
	keyValMatch := strings.Split(match, ":")
	return strings.ReplaceAll(keyValMatch[1], "\"", "")
}

// Find value based on key, split on '='
// Example string - https://pathtonowhere.com/?type=login&status=success&
//					request_token=tTy0wqusPbDObGf2zz7J0Wx9J5OYkFlp&action=login":

func extractKeyValue(body string, key string) string {
	arr := strings.Split(body, `&`) // split on '&'

	for index := range arr {
		if strings.Contains(arr[index], key) { // if key is found
			//fmt.Println(arr[index])
			arrVal := strings.Split(arr[index], `=`)
			//fmt.Println("Result 1: ", arrVal[1])
			return arrVal[1] // Extract value
		}
	}
	return ""
}
