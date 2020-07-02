package main

import (
	"fmt"
	"strings"
)

//CommandInterpreter receives username and message and interprets commands
func (bot *Bot) CommandInterpreter(wheel *Wheel, username string, message string, moderator bool) {
	username = strings.ToLower(username)
	message = strings.ToLower(message)

	if strings.HasPrefix(message, "!add1 ") {
		game := strings.TrimPrefix(message, "!add1")

		success := wheel.ChatAdd(username, game, 1)
		if success {
			bot.Message("Add successful")
		} else {
			bot.Message("Add unsuccessful")
		}

		fmt.Printf("%+v\n", wheel)
	} else if strings.HasPrefix(message, "!add2 ") {
		game := strings.TrimPrefix(message, "!add2 ")

		success := wheel.ChatAdd(username, game, 2)
		if success {
			bot.Message("Add successful")
		} else {
			bot.Message("Add unsuccessful")
		}
		fmt.Printf("%+v\n", wheel)
	} else if strings.HasPrefix(message, "!add3 ") {
		game := strings.TrimPrefix(message, "!add3 ")

		success := wheel.ChatAdd(username, game, 3)

		if success {
			bot.Message("Add successful")
		} else {
			bot.Message("Add unsuccessful")
		}

		fmt.Printf("%+v\n", wheel)
	} else if strings.HasPrefix(message, "!geturl") && moderator {
		reply := wheel.GetURL()
		bot.Message(reply)
		fmt.Println("get url called")
	}
	//TODO -> !geturl command for mods only

}
