package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//AccessToken stores a single access token from refreshAuth
type AccessToken struct {
	Token string `json:"access_token"`
}

//refresh authentication
func refreshAuth(refreshToken string, clientID string, secret string) string {
	url := fmt.Sprintf("https://id.twitch.tv/oauth2/token?grant_type=refresh_token&refresh_token=%s&client_id=%s&client_secret=%s", refreshToken, clientID, secret)
	resp, err := http.Post(url, "application/json", nil)

	if err != nil {
		fmt.Println("Error refreshing auth token")
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error in parsing refresh token response")
		log.Fatalln(err)
	}

	var NewToken AccessToken
	err = json.Unmarshal(body, &NewToken)

	if err != nil {
		log.Println("Error in unmarshalling refresh token response")
		log.Fatalln(err)
	}

	if NewToken.Token == "" {
		log.Println("Didn't receive refresh token from twitch")
		log.Fatalln(string(body))
	}
	return NewToken.Token
}

//CreatePaste for PasteBin
func CreatePaste(secret string, content string) string {
	hc := http.Client{}

	form := url.Values{}
	form.Add("api_dev_key", secret)
	form.Add("api_option", "paste")
	form.Add("api_paste_code", content)
	req, err := http.NewRequest("POST", "https://pastebin.com/api/api_post.php", strings.NewReader(form.Encode()))

	if err != nil {
		log.Panic("error constructing request to pastebin")
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded;charset=utf-8")

	resp, err := hc.Do(req)

	if err != nil {
		fmt.Println("error contacting pastebin service")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error in reading pastebin response")
		log.Fatalln(err)
	}

	return string(body)
}
