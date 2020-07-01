package main

import (
	"fmt"
	"strings"
	"time"
)

//CommandInterpreter receives username and message and interprets commands
func (bot *Bot) CommandInterpreter(wheel *Wheel, username string, message string) {
	username = strings.ToLower(username)
	message = strings.ToLower(message)

	if strings.Contains(message, "hey") {
		bot.Message("Hello " + username)
	} else if strings.Contains(message, "!peepo") {
		bot.Message("Poop")
	} else if strings.Contains(message, "!add1") {
		game := strings.ReplaceAll(message, "!add1", "")
		wheel.Add(game)
		user := newWheelUser(username)
		user.Add1(game)
		user.ResetTimeout = time.Now().Add(5 * time.Minute)
		wheel.Users[username] = *user

		fmt.Printf("%+v\n", wheel)
	}

}
