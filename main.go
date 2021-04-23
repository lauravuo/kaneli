package main

import (
	"fmt"
	"os"
)

const (
	argIndexPlayListID = 1
	argIndexRadioType  = 2
	sleepSeconds       = 5
)

type Track struct {
	Artist string
	Title  string
}

func main() {
	if len(os.Args) < (argIndexRadioType + 1) {
		panic("playlist id and radio type (christmas/esc) missing")
	}

	playlistID := os.Args[argIndexPlayListID]
	radioType := os.Args[argIndexRadioType]
	sToken := fetchUserToken()

	if radioType == "christmas" {
		addChristmasRadioLists(playlistID, sToken)
	} else if radioType == "esc" {
		addLatestEscRadioSongs(playlistID, sToken)
	}
	fmt.Printf("https://open.spotify.com/playlist/%s\n", playlistID)
}
