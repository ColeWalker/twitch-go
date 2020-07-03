package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

//PubSubClient is a twitch pubsub bot with all necessary functions
type PubSubClient struct {
	//Websocket connection
	conn *websocket.Conn
	//oauth token from twitch
	Token string
	//Channel id to listen for events
	ChannelID int
	//used to tell when we need to issue pings or reconnect to twitch server
	LastPing time.Time
	LastPong time.Time
	//Topics to listen to, currently unused
	Topics []string
	//Gorilla doesn't support concurrent writes
	writeLock *sync.Mutex
	//Wheel object to be modified
	Wheel *Wheel
}

//constructor
func newClient(auth string) *PubSubClient {
	return &PubSubClient{
		Token:     auth,
		Topics:    []string{},
		LastPing:  time.Now(),
		LastPong:  time.Now(),
		writeLock: &sync.Mutex{}}
}

//PingLoop constantly pings twitch, issues reconnect if last pong was too long ago
func (client *PubSubClient) PingLoop() {
	for {
		if time.Now().Unix()%60 == 0 {
			client.ping()
		}

		if client.LastPong.Sub(client.LastPing).Seconds() > (20 * time.Second).Seconds() {
			client.Connect()
			return
		}

		time.Sleep(time.Second * 1)
	}
}

//ping twitch server
func (client *PubSubClient) ping() {
	client.writeLock.Lock()
	defer client.writeLock.Unlock()

	err := client.conn.WriteJSON(&OutgoingMessage{
		Type: "PING",
	})

	if err != nil {
		fmt.Println("Error sending ping")
	}

	client.LastPing = time.Now()
}

//Connect to twitch server
func (client *PubSubClient) Connect() {
	//if client connection exists, close it
	if client.conn != nil {
		err := client.conn.Close()
		if err != nil {
			fmt.Println("error closing old connection")
		}
	}

	conn, _, err := websocket.DefaultDialer.Dial("wss://pubsub-edge.twitch.tv", nil)
	client.conn = conn
	fmt.Println("Connecting to PubSub server")
	if err != nil {
		return
	}

	//reset the timestamps to current timestamp
	client.LastPing = time.Now()
	client.LastPong = time.Now()

	go client.ReadLoop()
	go client.PingLoop()
	client.ping()

	//sub to topic
	//TODO -> set topic in slice, and use enumerated types to represent this junk, auto append channel id
	data := &RequestData{Topics: []string{"channel-points-channel-v1.23573216"}, Auth: client.Token}
	requestMessage := &Request{Type: "LISTEN", Data: *data}

	client.writeLock.Lock()
	defer client.writeLock.Unlock()
	err = client.conn.WriteJSON(requestMessage)

	if err != nil {
		log.Fatal("issue sending JSON")
	}

}

//ReadLoop handles websocket reads
func (client *PubSubClient) ReadLoop() {
	for {
		inc := &IncomingMessage{}
		err := client.conn.ReadJSON(inc)

		if err != nil {
			log.Println("closed read loop due to error")
			go client.Connect()
			return
		}

		if inc.Type == "PONG" {
			client.LastPong = time.Now()
		} else if inc.Data.Topic == "channel-points-channel-v1.23573216" {
			//enum would be nice here
			wrapper := &RedemptionWrapper{}
			err := json.Unmarshal([]byte(inc.Data.Message), &wrapper)

			if err != nil {
				log.Fatal("error unmarshaling commerce msg")
				return
			}

			lowerTitle := strings.ToLower(wrapper.Data.Redemption.Reward.Title)

			//add x to wheel
			addRewardRegex := regexp.MustCompile(`add (\d*[1-9]+) to wheel`)

			if addRewardRegex.MatchString(lowerTitle) {
				numRegex := regexp.MustCompile(`\d*[1-9]+`)

				raw := numRegex.FindString(lowerTitle)
				num, err := strconv.Atoi(raw)
				if err != nil {
					fmt.Println("Cannot convert " + raw + "to number")
				}

				game := wrapper.Data.Redemption.Input
				client.Wheel.AddMultiple(game, num)
			}
		}
	}
}
