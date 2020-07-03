package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

//Wheel stores information about Wheel and Users
type Wheel struct {
	//Wheeldecide URL to be added to
	URL string
	//Number of options passed into URL
	NumOptions int
	//Map of users that have added to the wheel
	Users map[string]*WheelUser
}

//Wheel constructor
func newWheel() *Wheel {
	return &Wheel{
		URL:        "https://wheeldecide.com/index.php?",
		NumOptions: 0,
		Users:      make(map[string]*WheelUser),
	}
}

//WheelUser stores wheel related user information
type WheelUser struct {
	//Username
	Name string
	//Stores whether user has used specific commands this wheel iteration
	Add1Available bool
	Add2Available bool
	Add3Available bool
	//After a few minutes, reset the user, this stores the time that user should be reset
	ResetTimeout time.Time
	//List of games already requested by user to prevent duplicates
	RequestedGames []string
	// Prevents user related race conditions
	Lock *sync.Mutex
}

//constructor
func newWheelUser(username string) *WheelUser {
	return &WheelUser{
		Name:           username,
		Add1Available:  true,
		Add2Available:  true,
		Add3Available:  true,
		ResetTimeout:   time.Now().Add(5 * time.Minute),
		RequestedGames: []string{},
		Lock:           &sync.Mutex{}}
}

//AddMultiple versions of game to wheel
func (wheel *Wheel) AddMultiple(game string, times int) {
	for i := 0; i < times; i++ {
		wheel.Add(game)
	}
}

//Add adds a game to the wheel object
func (wheel *Wheel) Add(game string) {
	newURL := wheel.URL

	//remove chars which will break url
	spaces := regexp.MustCompile(`\s`)
	illegalChars := regexp.MustCompile(`(\?|\&|\/|\\|\.|\'|\"|\:|\;|[\x{1F600}-\x{1F6FF}|[\x{2600}-\x{26FF}])|\%`)
	parsedGame := spaces.ReplaceAllLiteralString(game, "+")
	parsedGame = illegalChars.ReplaceAllLiteralString(parsedGame, "")

	//handle url construction
	if wheel.NumOptions == 0 {
		newURL = newURL + "c1=" + parsedGame
		wheel.NumOptions++
	} else {
		newURL = newURL + fmt.Sprintf("&c%d=%s", wheel.NumOptions+1, parsedGame)
		wheel.NumOptions++
	}
	wheel.URL = newURL
}

//ChatAdd handles additions done via chat
func (wheel *Wheel) ChatAdd(username string, game string, times int) error {
	user, ok := wheel.Users[username]

	//user lookup in wheel obj
	if ok {
		//timeout handling
		if time.Now().After(user.ResetTimeout) {
			fmt.Println(user.ResetTimeout)
			user.Reset()
			user.ResetTimeout = time.Now().Add(5 * time.Minute)
		}

	} else {
		user = newWheelUser(username)
		wheel.Users[username] = user
	}

	user.Lock.Lock()
	defer user.Lock.Unlock()

	if user.HasRequested(game) {
		return errors.New("You've already requested this game")
	}

	//Command interpreter, add1 = 1 time add2= 2 times, etc.
	switch times {
	case 1:
		if user.Add1Available {
			user.Add1(game)
			wheel.Add(game)
		} else {
			return errors.New("You've already used !add1")
		}
	case 2:
		if user.Add2Available {
			user.Add2(game)
			wheel.AddMultiple(game, 2)

		} else {
			return errors.New("You've already used !add2")
		}
	case 3:
		if user.Add3Available {
			user.Add3(game)
			wheel.AddMultiple(game, 3)

		} else {
			return errors.New("You've already used !add3")
		}
	}

	return nil
}

//Reset wheel and all users within wheel
func (wheel *Wheel) Reset() {
	wheel.URL = "https://wheeldecide.com/index.php?"
	wheel.NumOptions = 0
	wheel.Users = make(map[string]*WheelUser)
}

//GetURL from wheel and reset wheel
func (wheel *Wheel) GetURL() string {
	//append &time=30 decides animation time
	returnString := wheel.URL + "&time=30"
	fmt.Println("Wheel URL: " + wheel.URL)

	//checks if too long for twitch chat, sends it to pastebin
	if len(returnString) > 400 {
		returnString = CreatePaste(os.Getenv("pastebin-secret"), returnString)
	}

	wheel.Reset()

	fmt.Println("Returned URL: " + returnString)

	return returnString
}

//Reset User, set new ResetTimeout
func (user *WheelUser) Reset() {
	user.Add1Available = true
	user.Add1Available = true
	user.Add1Available = true
	user.ResetTimeout = time.Now().Add(5 * time.Minute)
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
	user.Add1Available = false
	user.RequestedGames = append(user.RequestedGames, game)
}

//Add2 to user
func (user *WheelUser) Add2(game string) {
	user.Add2Available = false
	user.RequestedGames = append(user.RequestedGames, game)
}

//Add3 to user
func (user *WheelUser) Add3(game string) {
	user.Add3Available = false
	user.RequestedGames = append(user.RequestedGames, game)
}
