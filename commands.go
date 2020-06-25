package main

import "strings"

//CommandInterpreter receives username and message and interprets commands
func (bot *Bot) CommandInterpreter(username string, message string) {
	username = strings.ToLower(username)
	message = strings.ToLower(message)

	if strings.Contains(message, "hey") {
		bot.Message("Hello " + username)
	} else if strings.Contains(message, "!peepo") {
		bot.Message("Poop")
	}

}
