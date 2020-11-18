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
	"time"

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
		panic(fmt.Errorf("Spotify client ID and secret missing!"))
	}

	token := ""
	state := strconv.FormatInt(int64(rand.Int()), 10)
	scope := "playlist-modify-public&playlist-modify-private"
	const redirectURI = "http://localhost:4321"
	path := fmt.Sprintf("https://accounts.spotify.com/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=%s&state=%s", clientID, redirectURI, scope, state)

	browser.OpenURL(path)
	server := &http.Server{Addr: ":4321"}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Spotify callback")
		if s, ok := r.URL.Query()["state"]; ok && s[0] == state {
			if code, ok := r.URL.Query()["code"]; ok {
				params := url.Values{}
				params.Add("grant_type", "authorization_code")
				params.Add("code", code[0])
				params.Add("redirect_uri", redirectURI)
				data, err := doPostRequest(
					"https://accounts.spotify.com/api/token",
					params,
					authHeader,
				)
				if err == nil {
					response := Token{}
					if err = json.Unmarshal(data, &response); err == nil {
						token = response.AccessToken
					}
				} else {
					panic(fmt.Errorf("Unable to acquire Spotify user token!"))
				}
			}
		}
		w.Write([]byte("OK, you can close this window"))

		// TODO: Do only after process completed
		time.AfterFunc(5*time.Second, func() {
			ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(ctxShutDown); err != nil {
				panic(fmt.Errorf("Unable to acquire Spotify user token!"))
			}
		})
	})
	log.Print(server.ListenAndServe())

	if token == "" {
		panic(fmt.Errorf("Unable to acquire Spotify user token!"))
	}
	return token
}
