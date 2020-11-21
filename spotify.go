package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/pkg/browser"
)

type Token struct {
	AccessToken string `json:"access_token"`
}

var (
	clientID     = os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	authHeader   = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret)))
)

func fetchSpotifyClientToken() string {
	if clientID == "" && clientSecret == "" {
		panic(fmt.Errorf("Spotify client ID and secret missing!"))
	}

	data, err := doPostRequest("https://accounts.spotify.com/api/token", url.Values{"grant_type": {"client_credentials"}}, authHeader)

	if err == nil {
		fmt.Println(string(data))
		response := Token{}
		if err = json.Unmarshal(data, &response); err == nil {
			return response.AccessToken
		}
	}

	panic(fmt.Errorf("Unable to acquire Spotify token!"))
}

func fetchSpotifyUserToken() string {
	if clientID == "" && clientSecret == "" {
		panic(fmt.Errorf("spotify client ID and secret missing"))
	}

	code := ""
	state := strconv.FormatInt(int64(rand.Int()), 10)
	scope := "playlist-modify-public&playlist-modify-private"
	const (
		redirectURL     = "http://localhost:4321"
		spotifyLoginURL = "https://accounts.spotify.com/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=%s&state=%s"
	)
	path := fmt.Sprintf(spotifyLoginURL, clientID, redirectURL, scope, state)

	if err := browser.OpenURL(path); err != nil {
		panic(fmt.Errorf("failed to open browser for authentication %s", err.Error()))
	}
	server := &http.Server{Addr: ":4321"}
	messages := make(chan bool)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Spotify callback", r.Method, r.URL.Query())
		if s, ok := r.URL.Query()["state"]; ok && s[0] == state {
			if codes, ok := r.URL.Query()["code"]; ok {
				code = codes[0]
				messages <- true
			}
		}
		http.Redirect(w, r, "https://www.spotify.com/", http.StatusSeeOther)
	})

	go func() {
		okToClose := <-messages
		if okToClose {
			if err := server.Shutdown(context.Background()); err != nil {
				log.Println("Failed to shutdown server", err)
			}
		}
	}()
	log.Println(server.ListenAndServe())

	params := url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("code", code)
	params.Add("redirect_uri", redirectURL)
	data, err := doPostRequest(
		"https://accounts.spotify.com/api/token",
		params,
		authHeader,
	)
	if err == nil {
		response := Token{}
		if err = json.Unmarshal(data, &response); err == nil {
			return response.AccessToken
		}
	}
	panic(fmt.Errorf("unable to acquire Spotify user token"))
}
