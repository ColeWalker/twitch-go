package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

//PubSubClient is a twitch pubsub bot with all necessary functions
type PubSubClient struct {
	conn      *websocket.Conn
	Token     string
	LastPing  int64
	LastPong  int64
	Topics    []string
	writeLock *sync.Mutex
}

func newClient(auth string) *PubSubClient {
	return &PubSubClient{
		Token:     auth,
		Topics:    []string{},
		LastPing:  time.Now().Unix(),
		LastPong:  time.Now().Unix(),
		writeLock: &sync.Mutex{}}
}

//PingLoop constantly pings twitch
func (client *PubSubClient) PingLoop() {
	for {
		if time.Now().Unix()%60 == 0 {
			client.ping()
		}
		time.Sleep(time.Second * 1)
	}
}

func (client *PubSubClient) ping() {
	client.writeLock.Lock()
	defer client.writeLock.Unlock()
	err := client.conn.WriteJSON(&OutgoingMessage{
		Type: "PING",
	})

	if err != nil {
		fmt.Println("Error sending ping")
	}
}
