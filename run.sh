#!/bin/bash

cd $(dirname "${BASH_SOURCE[0]}")
go install
source .envrc
~/go/bin/kaneli $SPOTIFY_PLAYLIST_ID $1