package main

import "testing"

func TestUserAuth(t *testing.T) {
	token := fetchUserToken()
	if token == "" {
		t.Error("error fetching token")
	}
	t.Log(token)
}
