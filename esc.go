package main

import (
	"fmt"
	"regexp"
	"strings"
)

func trimEscValue(value string) string {
	res := strings.ReplaceAll(value, "CDATA[ ", "")
	return strings.ReplaceAll(res, " ]", "")
}

func addLatestEscRadioSongs(playlistID, sToken string) {
	const url = "https://www.escradio.com/_playlist/playlist.xml"
	data, err := doGetRequest(url, "")
	if err != nil {
		panic(err)
	}

	tracks := make([]*Track, 0)

	regex := regexp.MustCompile(`CDATA\[(.*?)\]`)
	matches := regex.FindAllString(string(data), -1)
	for i := 0; i < len(matches); i += 2 {
		artist := trimEscValue(matches[i])
		title := trimEscValue(matches[i+1])
		cutIndex := strings.Index(title, "(")
		if cutIndex > 0 {
			title = title[:cutIndex]
		}
		tracks = append(tracks, &Track{Artist: artist, Title: title})
		fmt.Println(artist, title)
	}

	if err = addSongsFromRadioToPlaylist(tracks, playlistID, sToken); err != nil {
		panic(err)
	}

	fmt.Println("All done! Happy Eurovision!")
}
