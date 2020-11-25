package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	path "net/url"
	"os"
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
			ID  string `json:"id"`
			Uri string `json:"uri"`
		} `json:"items"`
	} `json:"tracks"`
}

type SpotifyRemoveItem struct {
	Uri string `json:"uri"`
}

type SpotifyPlaylistDelete struct {
	Tracks []*SpotifyRemoveItem `json:"tracks"`
}

type SpotifyPlaylistModify struct {
	Uris     []string `json:"uris"`
	Position int      `json:"position"`
}

const (
	argIndexPlayListID = 1
)

func main() {
	if len(os.Args) < (argIndexPlayListID + 1) {
		panic("please provide playlist id as command line argument")
	}

	playlistID := os.Args[argIndexPlayListID]
	sToken := fetchUserToken()

	// loop all interesting lists
	url := "https://jouluradio-wp.production.geniem.io/viimeksi-soitetut/"

	songIds := make([]string, 0)
	removeIds := make([]*SpotifyRemoveItem, 0)
	data, err := doGetRequest(url, "")
	authHeader = fmt.Sprintf("Bearer %s", sToken)
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
				if r, err := doGetRequest(fmt.Sprintf("https://api.spotify.com/v1/search?type=track&q=%s", q), authHeader); err == nil {
					resp := SpotifyResponse{}
					if jsonErr = json.Unmarshal(r, &resp); jsonErr == nil && len(resp.Tracks.Items) > 0 {
						songIds = append(songIds, resp.Tracks.Items[0].Uri)
						removeIds = append(removeIds, &SpotifyRemoveItem{Uri: resp.Tracks.Items[0].Uri})
					} else {
						fmt.Println(jsonErr)
					}
				}
			}

			fmt.Println(songIds)
			apiPath := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)
			if _, err = doJSONRequest(apiPath, http.MethodDelete, &SpotifyPlaylistDelete{Tracks: removeIds}, authHeader); err == nil {
				data, err := doJSONRequest(apiPath, http.MethodPost, &SpotifyPlaylistModify{Uris: songIds, Position: 0}, authHeader)
				fmt.Println(string(data))
				fmt.Println(err)
				//https://api.spotify.com/v1/playlists/{playlist_id}/tracks
			}
		}
	}
}
