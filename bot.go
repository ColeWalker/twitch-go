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
	"net/url"
	"os"
	"strings"
	"time"
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
	fmt.Fprintf(bot.conn, "CAP REQ :twitch.tv/membership\r\n")
	fmt.Fprintf(bot.conn, "CAP REQ :twitch.tv/tags\r\n")
}

//Message sent to server
func (bot *Bot) Message(message string) {
	if message != "" {
		fmt.Fprintf(bot.conn, "PRIVMSG "+bot.channel+" :"+message+"\r\n")
	}
	fmt.Println("PRIVMSG " + bot.channel + " :" + message + "\r\n")
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

func CreatePaste(secret string, content string) string {
	hc := http.Client{}

	form := url.Values{}
	form.Add("api_dev_key", secret)
	form.Add("api_option", "paste")
	form.Add("api_paste_code", content)
	req, err := http.NewRequest("POST", "https://pastebin.com/api/api_post.php", strings.NewReader(form.Encode()))
	if err != nil {
		log.Panic("error making newrequest")
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded;charset=utf-8")

	resp, err := hc.Do(req)

	if err != nil {
		fmt.Println("error contacting pastebin service")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error in ReadAll")
		log.Fatalln(err)
	}

	return string(body)
}

//Connect pubsubclient to websocket

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

	go client.PingLoop()
	go client.ReadLoop()

	for {
		line, err := tp.ReadLine()

		if err != nil {
			break
		} else if strings.Contains(line, "PING") {
			fmt.Fprintf(bot.conn, "PONG :tmi.twitch.tv")
		} else if strings.Contains(line, ".tmi.twitch.tv PRIVMSG "+bot.channel) {
			messageContents := strings.Split(line, ".tmi.twitch.tv PRIVMSG "+bot.channel)
			rawValues := strings.SplitN(messageContents[0], ":", 2)
			isMod := strings.Contains(messageContents[0], "mod=1")

			username := strings.Split(rawValues[1], "@")[1]
			message := messageContents[1][2:len(messageContents[1])]
			moderator := isMod || (strings.ToLower(username) == strings.TrimPrefix(bot.channel, "#"))
			log.Printf("mod:%t %s: %s\n", moderator, username, message)
			go bot.CommandInterpreter(wheel, username, message, moderator)
		}
	}
}
