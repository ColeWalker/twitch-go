package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

//AccessToken stores a single access token
type AccessToken struct {
	Token string `json:"access_token"`
}

type Bot struct {
	server  string
	port    string
	nick    string
	channel string
	mods    map[string]bool
	conn    net.Conn
}

func main() {
	token := refreshAuth(os.Getenv("twitch-bot-refresh-token"), os.Getenv("twitch-bot-client-id"), os.Getenv("twitch-bot-client-secret"))
	fmt.Println(token)
}

func refreshAuth(refreshToken string, clientID string, secret string) string {

	url := fmt.Sprintf("https://id.twitch.tv/oauth2/token?grant_type=refresh_token&refresh_token=%s&client_id=%s&client_secret=%s", refreshToken, clientID, secret)
	resp, err := http.Post(url, "application/json", nil)

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error in ReadAll")
		log.Fatalln(err)
	}

	var NewToken AccessToken
	err = json.Unmarshal(body, &NewToken)

	if err != nil {
		log.Println("Error in Unmarshall")
		log.Fatalln(err)
	}

	if NewToken.Token == "" {
		log.Fatalln(string(body))
	}
	return NewToken.Token
}
