package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//CommandInterpreter receives username and message and interprets commands
func (bot *Bot) CommandInterpreter(username string, message string, moderator bool) {
	lowerUsername := strings.ToLower(username)
	lowerMessage := strings.ToLower(message)

	addCommand := regexp.MustCompile(`\!add[1-3] .*`)

	if addCommand.MatchString(message) {
		numRaw := strings.TrimPrefix(message, "!add")[0]
		num, _ := strconv.Atoi(string(numRaw))

		game := strings.SplitN(message, " ", 2)[1]

		err := bot.Wheel.ChatAdd(lowerUsername, game, num)
		if err != nil {
			bot.Message("Add unsuccessful: " + err.Error())
		} else {
			bot.Message("Add successful")
		}
	} else if strings.HasPrefix(lowerMessage, "!geturl") && moderator {
		reply := bot.Wheel.GetURL()
		fmt.Println(reply)
		go bot.Message(reply)
		fmt.Println("get url called")
	}

}
