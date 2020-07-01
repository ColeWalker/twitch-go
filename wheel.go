package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

//Wheel stores information about Wheel and Users
type Wheel struct {
	URL        string
	NumOptions int
	Users      map[string]WheelUser
}

func newWheel() *Wheel {
	return &Wheel{
		URL:        "https://wheeldecide.com/index.php?",
		NumOptions: 0,
		Users:      make(map[string]WheelUser),
	}
}

//WheelUser stores wheel related user information
type WheelUser struct {
	Name           string
	HasUsedAdd1    bool
	HasUsedAdd2    bool
	HasUsedAdd3    bool
	ResetTimeout   time.Time
	RequestedGames []string
}

func newWheelUser(username string) *WheelUser {
	return &WheelUser{
		Name:           username,
		HasUsedAdd1:    false,
		HasUsedAdd2:    false,
		HasUsedAdd3:    false,
		ResetTimeout:   time.Now().Add(-5 * time.Minute),
		RequestedGames: []string{}}
}

//AddMultiple times to wheel
func (wheel *Wheel) AddMultiple(game string, times int) {
	for i := 0; i < times; i++ {
		wheel.Add(game)
	}
}

//Add adds a game to the wheel object
func (wheel *Wheel) Add(game string) {
	newURL := wheel.URL
	spaces := regexp.MustCompile(`\s`)

	illegalChars := regexp.MustCompile(`(\?|\&|\/|\\|\.|\'|\:|\;)`)
	parsedGame := spaces.ReplaceAllLiteralString(game, "+")
	parsedGame = illegalChars.ReplaceAllLiteralString(parsedGame, "")

	if wheel.NumOptions == 0 {
		newURL = newURL + "c1=" + parsedGame
		wheel.NumOptions++
	} else {
		newURL = newURL + fmt.Sprintf("c%d=%s", wheel.NumOptions+1, parsedGame)
		wheel.NumOptions++
	}
	wheel.URL = newURL
}

//Reset wheel
func (wheel *Wheel) Reset() {
	wheel.URL = "https://wheeldecide.com/index.php?"
	wheel.NumOptions = 0
	wheel.Users = make(map[string]WheelUser)
}

//Reset User
func (user *WheelUser) Reset() {
	user.HasUsedAdd1 = false
	user.HasUsedAdd2 = false
	user.HasUsedAdd3 = false
	user.ResetTimeout = time.Now().Add(-5 * time.Minute)
	user.RequestedGames = []string{}
}

//HasRequested returns whether user has requested a game
func (user *WheelUser) HasRequested(game string) bool {
	for _, v := range user.RequestedGames {
		isGame := false
		spaces, _ := regexp.Compile(`\s`)
		request := spaces.ReplaceAllLiteralString(v, "")
		request = strings.ToLower(request)
		parsedGame := spaces.ReplaceAllLiteralString(game, "")
		parsedGame = strings.ToLower(parsedGame)

		isGame = strings.Contains(request, parsedGame) || strings.Contains(parsedGame, request)

		if isGame {
			return isGame
		}
	}
	return false
}

//Add1 to user
func (user *WheelUser) Add1(game string) {
	user.HasUsedAdd1 = true
	user.RequestedGames = append(user.RequestedGames, game)
}

//Add2 to user
func (user *WheelUser) Add2(game string) {
	user.HasUsedAdd2 = true
	user.RequestedGames = append(user.RequestedGames, game)
}

//Add3 to user
func (user *WheelUser) Add3(game string) {
	user.HasUsedAdd3 = true
	user.RequestedGames = append(user.RequestedGames, game)
}
