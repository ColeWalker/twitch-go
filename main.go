package main

import (
	"log"
	"os"
)

func main() {
	token := refreshAuth(os.Getenv("twitch-bot-refresh-token"), os.Getenv("twitch-bot-client-id"), os.Getenv("twitch-bot-client-secret"))

	bot := newBot(token)
	bot.Connect()

	client := newClient(token)
	client.Connect()

	defer client.conn.Close()

	data := &RequestData{Topics: []string{"channel-points-channel-v1.23573216"}, Auth: token}
	requestMessage := &Request{Type: "LISTEN", Data: *data}

	client.writeLock.Lock()
	err := client.conn.WriteJSON(requestMessage)
	client.writeLock.Unlock()

	if err != nil {
		log.Fatal("issue sending JSON")
	}
	wheel := newWheel()

	client.Wheel = wheel
	bot.Wheel = wheel
	go client.PingLoop()
	go client.ReadLoop()

	for {
	}
}
