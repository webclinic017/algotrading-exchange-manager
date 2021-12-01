package kite

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/asmcos/requests"
	"github.com/joho/godotenv"
	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

const (
	URL      = "https://kite.zerodha.com/api/login"
	twofaUrl = "https://kite.zerodha.com/api/twofa"
)

func LoginKite() (bool, string, string) {
	// := os.Getenv("TFA_AUTH")
	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")

	requestToken := KiteGetRequestToken()

	if strings.Contains(requestToken, "ERR:") {
		fmt.Println("Authentication", requestToken)
	} else {

		fmt.Println("Authentication Succesful:", requestToken)
		// Create a new Kite connect instance
		kc := kiteconnect.New(apiKey)

		// Get user details and access token
		data, err := kc.GenerateSession(requestToken, apiSecret)
		if err != nil {
			fmt.Printf("Session Err: %v", err)
			return false, "", ""
		}

		// Set access token
		kc.SetAccessToken(data.AccessToken)
		fmt.Println("AccessToken", data.AccessToken)

		// keypair := strings.Join("accessToken", data.AccessToken)
		env, _ := godotenv.Unmarshal("accessToken=" + data.AccessToken)
		err = godotenv.Write(env, "./config/ENV_accesstoken.env")
		if err != nil {
			fmt.Println("Cannot write to accesstoken.env", err)
		}

		// Get margins
		margins, err := kc.GetUserMargins()
		if err != nil {
			fmt.Printf("Error getting margins: %v", err)
			//return false, "", ""
		}
		fmt.Println("margins: ", margins)

		return true, apiKey, data.AccessToken

	}
	return false, "", ""
}

func KiteGetRequestToken() string {

	tfAuth := os.Getenv("TFA_AUTH")
	userId := os.Getenv("USER_ID")
	userPwd := os.Getenv("PASSWORD")
	reqTokenUrl := os.Getenv("REQUEST_TOKEN_URL")

	requestToken := ""

	// 1. Start login, get reqId
	data := requests.Datas{
		"user_id":  userId,
		"password": userPwd,
	}

	req := requests.Requests()
	resp, err := req.Post(URL, data)

	if (err != nil) || (resp.R.StatusCode != 200) {
		requestToken = "ERR: User ID / Password error"
		println(resp.Text())
		return requestToken
	}

	//{"status":"success","data":{"user_id":"ZY7293","request_id":"mJHDflgxtNutiffR3jAZcIPA5SH4eXaLNkrwSgheQfoyb4syN9vJ5OJyGnPRKCi7","twofa_type":"pin","twofa_status":"active"}}
	reqID := extractValue(resp.Text(), "request_id")

	// 2. Do Two factor auth
	data = requests.Datas{
		"user_id":     userId,
		"request_id":  reqID,
		"twofa_value": tfAuth,
	}
	resp, err = req.Post(twofaUrl, data)

	if (err != nil) || (resp.R.StatusCode != 200) {
		println(resp.Text())
		requestToken = "ERR: Two factor auth failed"
		return requestToken
	}

	// 3. Post login, access URL to get requestToken
	resp, err = req.Get(reqTokenUrl)
	//action=login&type=login&status=success&request_token=5cs4Q9q52133Y4iL33FPwkkv37ewpLuV

	m, _ := url.ParseQuery(resp.R.Request.URL.RawQuery)
	if (err != nil) || (resp.R.StatusCode != 200) {
		println(resp.R.Request.URL.RawQuery)
		requestToken = "ERR: Cannot fetch request token"
		return requestToken
	}

	// fmt.Println(m)
	requestToken = m["request_token"][0]
	// fmt.Println(requestToken)

	return requestToken
}

func extractValue(body string, key string) string {
	keystr := "\"" + key + "\":[^,;\\]}]*"
	r, _ := regexp.Compile(keystr)
	match := r.FindString(body)
	keyValMatch := strings.Split(match, ":")
	return strings.ReplaceAll(keyValMatch[1], "\"", "")
}
