package main

import (
	"fmt"
	"regexp"
)

//Wheel stores information about Wheel and Users
type Wheel struct {
	URL        string
	NumOptions int
	Users      map[string]WheelUser
}

//WheelUser stores wheel related user information
type WheelUser struct {
	Name           string
	LastAdd3       int64
	LastAdd2       int64
	LastAdd1       int64
	ResetTimeout   int64
	RequestedGames []string
}

//AddToWheel adds a game to the wheel object
func (wheel *Wheel) AddToWheel(game string) {
	newURL := wheel.URL
	illegalChars := regexp.MustCompile(`/( |\?|\&|\/|\\|\.|\'|\:|\;)/g`)

	parsedGame := illegalChars.Copy().ReplaceAllLiteralString(game, "")

	if wheel.NumOptions == 0 {
		newURL = newURL + "c1=" + parsedGame
		wheel.NumOptions++
	} else {
		newURL = newURL + fmt.Sprintf("c%d=%s", wheel.NumOptions+1, parsedGame)
		wheel.NumOptions++
	}
	wheel.URL = newURL
}
