run:
	go run ./... $(SPOTIFY_PLAYLIST_ID)

build:
	go build ./...

lint:
	golangci-lint run