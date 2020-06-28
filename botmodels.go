package main

import (
	"net"
	"time"

	"github.com/gorilla/websocket"
)

//AccessToken stores a single access token
type AccessToken struct {
	Token string `json:"access_token"`
}

//The Bot object which stores websocket info
type Bot struct {
	server  string
	port    string
	nick    string
	channel string
	mods    map[string]bool
	conn    net.Conn
}

func newBot() *Bot {
	return &Bot{
		server:  "irc.twitch.tv",
		port:    "6667",
		nick:    "supcole",
		channel: "#supcole",
		mods:    make(map[string]bool),
		conn:    nil}
}

//RequestData for data field of Request
type RequestData struct {
	Topics []string `json:"topics"`
	Auth   string   `json:"auth_token"`
}

//Request for twitch websocket
type Request struct {
	Type string      `json:"type"`
	Data RequestData `json:"data"`
}

//PubSubClient is a twitch pubsub bot with all necessary functions
type PubSubClient struct {
	conn     *websocket.Conn
	Token    string
	LastPing int64
	LastPong int64
	Topics   []string
}

func newClient(auth string) *PubSubClient {
	return &PubSubClient{
		Token:    auth,
		Topics:   []string{},
		LastPing: time.Now().Unix(),
		LastPong: time.Now().Unix()}
}
