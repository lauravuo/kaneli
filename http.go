package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var client = http.Client{
	Timeout: time.Second * 2, // Timeout after 2 seconds
}

func doGetRequest(url, authHeader string) (data []byte, err error) {
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, url, nil)
	if err != nil {
		return
	}

	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	return
}

func doPostRequest(url string, values url.Values, authHeader string) (data []byte, err error) {
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, url, strings.NewReader(values.Encode()))
	if err != nil {
		return
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return
	}

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	return

}
