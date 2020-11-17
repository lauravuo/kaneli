package main

import (
	"encoding/json"
	"fmt"
	path "net/url"
)

type Song struct {
	Artist  string `json:"artist"`
	Channel int32  `json:"channel"`
	Song    string `json:"song"`
	Time    string `json:"time"`
}

type RadioResponse struct {
	PageLastPlayed struct {
		Content struct {
			Hero struct {
				Text string `json:"text"`
			} `json:"Hero"`
			RecentlyPlayed struct {
				Songs []Song `json:"songs"`
			} `json:"recently_played"`
		} `json:"Content"`
	} `json:"PageLastPlayed"`
}

type SpotifyResponse struct {
	Tracks struct {
		Items []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			ID string `json:"id"`
		} `json:"items"`
	} `json:"tracks"`
}

func main() {

	sToken := fetchSpotifyToken()

	url := "https://jouluradio-wp.production.geniem.io/viimeksi-soitetut/"

	data, err := doGetRequest(url, "")
	if err != nil {
		fmt.Println(err)
	} else {
		response := RadioResponse{}
		jsonErr := json.Unmarshal(data, &response)
		if jsonErr != nil {
			fmt.Println(jsonErr)
		} else {
			fmt.Println(response)
			for _, track := range response.PageLastPlayed.Content.RecentlyPlayed.Songs {
				q := path.QueryEscape(fmt.Sprintf("artist:%s track:%s", track.Artist, track.Song))
				r, err := doGetRequest(fmt.Sprintf("https://api.spotify.com/v1/search?type=track&q=%s", q), fmt.Sprintf("Bearer %s", sToken))
				resp := SpotifyResponse{}
				jsonErr = json.Unmarshal(r, &resp)
				fmt.Println(resp)
				fmt.Println(err)
			}
		}
	}
}
