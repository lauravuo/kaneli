run:
	go run . $(SPOTIFY_PLAYLIST_ID) christmas

build:
	go build ./...

lint:
	golangci-lint run