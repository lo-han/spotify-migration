# Spotify Migration

Automation to migrate playlist or albums from Spotify to Youtube

_**Note**: Currently album migration is not supported_

## Setting

### Server

Create the necessary certificates to run local server

```
openssl genrsa -out server.key 2048
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```

### Spotify

Register at https://developer.spotify.com/dashboard.

### Youtube

1. Register a new resource at https://console.cloud.google.com/cloud-resource-manager
2. Register a new client at https://console.cloud.google.com/auth/clients
3. Add Youtube API V3 scope to the new client at https://console.cloud.google.com/auth/scopes

### Environment

After Youtube and Spotify setup, add your credentials to a new `.env` file, like the following:

```
SPOTIFY_ID=
SPOTIFY_SECRET=

YOUTUBE_ID=
YOUTUBE_SECRET=
```

Then export your environment

``` bash
export $(cat .env | xargs)
```

## Run the app

Migrate playlist

``` bash
go run . playlist "<your_playlist_name>"
```

Migrate album
``` bash
go run . album "<album_name>"
```

_**Note**:You must authorize both Spotify and Youtube clicking the printed CLI links_

## Contributing

Please don't commit to *master* branch. Thank you.
