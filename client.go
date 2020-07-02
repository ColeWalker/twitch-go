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
	conn      *websocket.Conn
	Token     string
	LastPing  time.Time
	LastPong  time.Time
	Topics    []string
	writeLock *sync.Mutex
	Wheel     *Wheel
}

func newClient(auth string) *PubSubClient {
	return &PubSubClient{
		Token:     auth,
		Topics:    []string{},
		LastPing:  time.Now(),
		LastPong:  time.Now(),
		writeLock: &sync.Mutex{}}
}

//PingLoop constantly pings twitch
func (client *PubSubClient) PingLoop() {
	for {
		if time.Now().Unix()%60 == 0 {
			client.ping()
		}

		// if time.Now().Unix()%10 == 0 {

		if client.LastPong.Sub(client.LastPing).Seconds() > (20 * time.Second).Seconds() {
			client.reconnect()
			fmt.Println(client.LastPong.Sub(client.LastPing).Seconds)
			fmt.Println(20 * time.Second)
			return
		}
		// }

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

	client.LastPing = time.Now()
}

func (client *PubSubClient) Connect() {
	var err error
	client.conn, _, err = websocket.DefaultDialer.Dial("wss://pubsub-edge.twitch.tv", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

}

func (client *PubSubClient) reconnect() {
	if client.conn != nil {
		err := client.conn.Close()
		if err != nil {
			fmt.Println("error closing old connection")
		}
	}

	conn, _, err := websocket.DefaultDialer.Dial("wss://pubsub-edge.twitch.tv", nil)
	client.conn = conn
	fmt.Println("reconnecting...")
	if err != nil {
		return
	}

	//reset the timestamps to current timestamp
	client.LastPing = time.Now()
	client.LastPong = time.Now()

	go client.ReadLoop()
	go client.PingLoop()
	client.ping()

	//sub to topics
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
			client.reconnect()
			return
		}

		log.Printf("recv: %+v", inc)

		if inc.Data.Topic == "channel-points-channel-v1.23573216" {
			wrapper := &RedemptionWrapper{}
			err := json.Unmarshal([]byte(inc.Data.Message), &wrapper)

			if err != nil {
				log.Fatal("error unmarshaling commerce msg")
				return
			}

			//TODO -> regular expression this or strings.prefix + strings.suffix ?
			lowerTitle := strings.ToLower(wrapper.Data.Redemption.Reward.Title)
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
		} else if inc.Type == "PONG" {
			client.LastPong = time.Now()
		}
	}
}
