package main

import (
	"fmt"
	"os"
	"time"
)

const (
	argIndexPlayListID = 1
	sleepSeconds       = 5
)

func main() {
	if len(os.Args) < (argIndexPlayListID + 1) {
		panic("please provide playlist id as command line argument")
	}

	playlistID := os.Args[argIndexPlayListID]
	sToken := fetchUserToken()

	// loop all interesting lists
	radioUrls := []string{
		"https://jouluradio-wp.production.geniem.io/viimeksi-soitetut/",
		"https://jouluradio-wp.production.geniem.io/viimeksi-soitetut/?kanava=indiejoulu",
		"https://jouluradio-wp.production.geniem.io/viimeksi-soitetut/?kanava=jazzjoulu",
		"https://jouluradio-wp.production.geniem.io/viimeksi-soitetut/?kanava=klassinen-joulu",
	}
	for _, radioURL := range radioUrls {
		fmt.Printf("Adding songs from %s\n", radioURL)
		if err := addSongsFromRadioToPlaylist(radioURL, playlistID, sToken); err != nil {
			fmt.Printf("Error when adding songs %s\n", err.Error())
		}
		fmt.Printf("Sleeping a while...\n")
		time.Sleep(sleepSeconds * time.Second)
	}

	fmt.Println("All done! Merry Christmas! Hyvää joulua!")
	fmt.Printf("https://open.spotify.com/playlist/%s\n", playlistID)
}
