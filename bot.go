package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"strings"
	"time"
)

//The Bot object which stores irc info
type Bot struct {
	//server to connect to
	server string
	//port to connect to
	port string
	//nickname of bot
	nick string
	//channel to join - #username format
	channel string
	//mods of channel -> currently unused
	mods map[string]bool
	//connection to irc
	conn net.Conn
	//Wheel to modify with commands
	Wheel *Wheel
	//oauth token containing chat:read and chat:write scopes
	AuthToken string
}

//constructor
func newBot(token string) *Bot {
	return &Bot{
		server:    "irc.twitch.tv",
		port:      "6667",
		nick:      "supcole",
		channel:   "#supcole",
		mods:      make(map[string]bool),
		conn:      nil,
		AuthToken: token}
}

//Connect and connect bot to IRC server
func (bot *Bot) Connect() {
	//close connection if it exists
	if bot.conn != nil {
		err := bot.conn.Close()
		if err != nil {
			fmt.Println("error closing old connection")
		}
	}

	var err error
	fmt.Println("Attempting to connect to Twitch IRC server!")
	bot.conn, err = net.Dial("tcp", bot.server+":"+bot.port)

	if err != nil {
		fmt.Printf("Unable to connect to Twitch IRC server! Reconnecting in 10 seconds...\n")
		time.Sleep(10 * time.Second)
		bot.Connect()
	}

	fmt.Printf("Connected to IRC server %s\n", bot.server)

	fmt.Fprintf(bot.conn, "USER %s 8 * :%s\r\n", bot.nick, bot.nick)
	fmt.Fprintf(bot.conn, "PASS oauth:%s\r\n", bot.AuthToken)
	fmt.Fprintf(bot.conn, "NICK %s\r\n", bot.nick)
	fmt.Fprintf(bot.conn, "JOIN %s\r\n", bot.channel)
	fmt.Fprintf(bot.conn, "CAP REQ :twitch.tv/membership\r\n")
	fmt.Fprintf(bot.conn, "CAP REQ :twitch.tv/tags\r\n")

	go bot.ReadLoop()
}

//Message to IRC channel
func (bot *Bot) Message(message string) {
	if message != "" {
		fmt.Fprintf(bot.conn, "PRIVMSG "+bot.channel+" :"+message+"\r\n")
	}
	fmt.Println("PRIVMSG " + bot.channel + " :" + message + "\r\n")
}

//ReadLoop for Bot
func (bot *Bot) ReadLoop() {
	reader := bufio.NewReader(bot.conn)
	tp := textproto.NewReader(reader)

	for {
		line, err := tp.ReadLine()

		if err != nil {
			fmt.Println("Bot read loop exited due to error")
			bot.Connect()
			break
		} else if strings.Contains(line, "PING") {
			fmt.Fprintf(bot.conn, "PONG :tmi.twitch.tv")
		} else if strings.Contains(line, ".tmi.twitch.tv PRIVMSG "+bot.channel) {
			//TODO -> clean this up with regex
			messageContents := strings.Split(line, ".tmi.twitch.tv PRIVMSG "+bot.channel)

			//TODO -> parse this into struct with info

			//rawValues[0] = IRC3 tags
			//rawValues[1] = raw message
			rawValues := strings.SplitN(messageContents[0], ":", 2)
			isMod := strings.Contains(messageContents[0], "mod=1")

			username := strings.Split(rawValues[1], "@")[1]
			message := messageContents[1][2:len(messageContents[1])]
			moderator := isMod || (strings.ToLower(username) == strings.TrimPrefix(bot.channel, "#"))

			//print basic message to console
			log.Printf("mod:%t %s: %s\n", moderator, username, message)

			//send message to command interpreter
			go bot.CommandInterpreter(username, message, moderator)
		}
	}

}
