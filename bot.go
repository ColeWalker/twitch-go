package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/textproto"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

//Connect bot to IRC server
func (bot *Bot) Connect(token string) {
	var err error
	fmt.Println("Attempting to connect to Twitch IRC server!")
	bot.conn, err = net.Dial("tcp", bot.server+":"+bot.port)

	if err != nil {
		fmt.Printf("Unable to connect to Twitch IRC server! Reconnecting in 10 seconds...\n")
		time.Sleep(10 * time.Second)
		bot.Connect(token)
	}
	fmt.Printf("Connected to IRC server %s\n", bot.server)

	fmt.Fprintf(bot.conn, "USER %s 8 * :%s\r\n", bot.nick, bot.nick)
	fmt.Fprintf(bot.conn, "PASS oauth:%s\r\n", token)
	fmt.Fprintf(bot.conn, "NICK %s\r\n", bot.nick)
	fmt.Fprintf(bot.conn, "JOIN %s\r\n", bot.channel)
}

//Message sent to server
func (bot *Bot) Message(message string) {
	if message != "" {
		fmt.Fprintf(bot.conn, "PRIVMSG "+bot.channel+" :"+message+"\r\n")
	}
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
		log.Println("Error in Unmarshal")
		log.Fatalln(err)
	}

	if NewToken.Token == "" {
		log.Fatalln(string(body))
	}
	return NewToken.Token
}

//Connect pubsubclient to websocket
func (client *PubSubClient) Connect() {
	var err error
	client.conn, _, err = websocket.DefaultDialer.Dial("wss://pubsub-edge.twitch.tv", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

}
func main() {
	token := refreshAuth(os.Getenv("twitch-bot-refresh-token"), os.Getenv("twitch-bot-client-id"), os.Getenv("twitch-bot-client-secret"))
	fmt.Println(token)

	bot := newBot()
	bot.Connect(token)

	defer bot.conn.Close()

	reader := bufio.NewReader(bot.conn)
	tp := textproto.NewReader(reader)

	client := newClient(token)
	client.Connect()

	defer client.conn.Close()

	done := make(chan struct{})

	data := &RequestData{Topics: []string{"channel-points-channel-v1.23573216"}, Auth: token}
	requestMessage := &Request{Type: "LISTEN", Data: *data}

	err := client.conn.WriteJSON(requestMessage)

	if err != nil {
		log.Fatal("issue sending JSON")
	}

	wheel := newWheel()
	wheel.Add("hello world/\\\"'+&")
	fmt.Printf("%+v\n", wheel)
	go func() {
		defer close(done)
		for {
			inc := &IncomingMessage{}
			err := client.conn.ReadJSON(inc)
			if err != nil {
				log.Println("read:", err)
				return
			}

			log.Printf("recv: %+v\n", inc.Data.Message)

			if inc.Data.Topic == "channel-points-channel-v1.23573216" {
				wrapper := &RedemptionWrapper{}
				err := json.Unmarshal([]byte(inc.Data.Message), &wrapper)

				if err != nil {
					log.Fatal("error unmarshaling commerce msg")
					return
				}

				if wrapper.Data.Redemption.Reward.Title == "Add 1 to wheel" {
					// wheel.Add(wrapper.Data.Redemption.Input)
					// fmt.Printf("%+v\n", wheel)
				}
			}
		}
	}()

	for {
		line, err := tp.ReadLine()
		if err != nil {
			break
		} else if strings.Contains(line, "PING") {
			fmt.Fprintf(bot.conn, "PONG :tmi.twitch.tv")
		} else if strings.Contains(line, ".tmi.twitch.tv PRIVMSG "+bot.channel) {
			messageContents := strings.Split(line, ".tmi.twitch.tv PRIVMSG "+bot.channel)
			username := strings.Split(messageContents[0], "@")[1]
			message := messageContents[1][2:len(messageContents[1])]

			go bot.CommandInterpreter(wheel, username, message)
		}
	}
}
