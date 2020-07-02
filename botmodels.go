package main

import (
	"net"
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
