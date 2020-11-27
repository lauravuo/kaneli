# kaneli

Kaneli (cinnamon in Finnish) is a command line program that adds songs from [Jouluradio](https://www.jouluradio.fi/)'s playlist to a user chosen Spotify playlist.

## Setup

1. Install [golang](https://golang.org/)
1. Create new application from [Spotify developer dashboard](https://developer.spotify.com/dashboard/)
1. Copy the app client id and secret from the Spotify developer UI. Define following env variables:

```bash
export SPOTIFY_CLIENT_ID=xxx
export SPOTIFY_CLIENT_SECRET=xxx
```

## Usage

1. In the Spotify desktop UI: create a new playlist (or use an existing playlist.) Copy the Spotify URI for the playlist (e.g. `spotify:playlist:5x5mdsVit4ngNyvglqkO8f`).
1. Run the app by giving the playlist id (without the `spotify:playlist:` part) as command line argument:

    ```bash
    go run . 5x5mdsVit4ngNyvglqkO8f
    ```
    
1. The app will open browser for Spotify authentication. Log in and give app permission to modify your playlists.
1. Bunch of random Christmas songs are added to your chosen playlist:

    ```bash
    go run . 5x5mdsVit4ngNyvglqkO8f
    2020/11/27 21:55:53 http: Server closed
    ...
    Add track Oscar Peterson: God Rest Ye Merry, Gentlemen
    Add track Sister Rosetta Tharpe: O Little Town Of Bethlehem
    Add track Dianne Reeves: Jingle Bells
    Sleeping a while...
    ...
    All done! Merry Christmas! Hyvää joulua!
    https://open.spotify.com/playlist/5x5mdsVit4ngNyvglqkO8f
    ```
1. Start listening and find new favourite Christmas tunes 🎄!