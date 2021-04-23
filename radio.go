package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type SpotifyResponse struct {
	Tracks struct {
		Items []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			ID  string `json:"id"`
			URI string `json:"uri"`
		} `json:"items"`
	} `json:"tracks"`
}

type SpotifyRemoveItem struct {
	URI string `json:"uri"`
}

type SpotifyPlaylistDelete struct {
	Tracks []*SpotifyRemoveItem `json:"tracks"`
}

type SpotifyPlaylistModify struct {
	URIs     []string `json:"uris"`
	Position int      `json:"position"`
}

func addSongsFromRadioToPlaylist(tracks []*Track, playlistID, spotifyToken string) (err error) {
	songIds := make([]string, 0)
	removeIds := make([]*SpotifyRemoveItem, 0)
	bearerHeader := fmt.Sprintf("Bearer %s", spotifyToken)

	for _, track := range tracks {
		query := url.QueryEscape(fmt.Sprintf("artist:%s track:%s", track.Artist, track.Title))
		// search for song by artist and title
		trackResponse, trackErr := doGetRequest(fmt.Sprintf("https://api.spotify.com/v1/search?type=track&q=%s", query), bearerHeader)
		if trackErr != nil {
			fmt.Printf("Unable to fetch track data %s\n", err.Error())
			continue
		}

		trackData := SpotifyResponse{}
		err = json.Unmarshal(trackResponse, &trackData)
		if err != nil {
			fmt.Printf("Unable to parse track data %s\n", err.Error())
			continue
		}

		// just pick the first found track
		if len(trackData.Tracks.Items) > 0 {
			item := trackData.Tracks.Items[0]
			fmt.Printf("Add track %s: %s\n", track.Artist, track.Title)
			songIds = append(songIds, item.URI)
			removeIds = append(removeIds, &SpotifyRemoveItem{URI: item.URI})
		}
	}

	apiPath := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)

	// first delete all tracks with similar id to avoid duplicates
	_, err = doJSONRequest(apiPath, http.MethodDelete, &SpotifyPlaylistDelete{Tracks: removeIds}, bearerHeader)
	if err != nil {
		return err
	}

	// then add all tracks to the start of the list
	_, err = doJSONRequest(apiPath, http.MethodPost, &SpotifyPlaylistModify{URIs: songIds, Position: 0}, bearerHeader)
	return err
}
