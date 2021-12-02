package kite

import (
	"fmt"
	"log"
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
		err = godotenv.Write(env, "./app/config/ENV_accesstoken.env")
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

	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()

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
	req.SetTimeout(5)
	resp, err = req.Get(reqTokenUrl)
	if err != nil {
		println(err.Error())
		arr := strings.Split(err.Error(), `"`) // split on '&'
		requestToken = extractKeyValue(arr[1], "request_token")
		if requestToken == "" {
			requestToken = "ERR: Cannot fetch request token"
			return requestToken
		}
	} else {

		m, err := url.ParseQuery(resp.R.Request.URL.RawQuery)
		if (err != nil) || (resp.R.StatusCode != 200) {

			requestToken = "ERR: Cannot fetch request token"
			return requestToken
		}

		fmt.Println("parsed m:", m)
		requestToken = m["request_token"][0]
	}

	fmt.Println("extraced req token:", requestToken)

	return requestToken
}

func extractValue(body string, key string) string {
	keystr := "\"" + key + "\":[^,;\\]}]*"
	r, _ := regexp.Compile(keystr)
	match := r.FindString(body)
	keyValMatch := strings.Split(match, ":")
	return strings.ReplaceAll(keyValMatch[1], "\"", "")
}

// Find value based on key, split on '='
// Example string - https://pathtonowhere.com/?type=login&status=success&request_token=tTy0wqusPbDObGf2zz7J0Wx9J5OYkFlp&action=login":

func extractKeyValue(body string, key string) string {
	arr := strings.Split(body, `&`) // split on '&'

	for index, _ := range arr {
		if strings.Contains(arr[index], key) { // if key is found
			//fmt.Println(arr[index])
			arrVal := strings.Split(arr[index], `=`)
			//fmt.Println("Result 1: ", arrVal[1])
			return arrVal[1] // Extract value
		}
	}
	return ""
}
