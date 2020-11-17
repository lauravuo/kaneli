package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

type Token struct {
	AccessToken string `json:"access_token"`
}

func fetchSpotifyToken() string {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if clientID == "" && clientSecret == "" {
		panic(fmt.Errorf("Spotify client ID and secret missing!"))
	}

	auth := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret)))

	data, err := doPostRequest("https://accounts.spotify.com/api/token", url.Values{"grant_type": {"client_credentials"}}, auth)

	if err == nil {
		fmt.Println(string(data))
		response := Token{}
		if err = json.Unmarshal(data, &response); err == nil {
			return response.AccessToken
		}
	}

	panic(fmt.Errorf("Unable to acquire Spotify token!"))
}
