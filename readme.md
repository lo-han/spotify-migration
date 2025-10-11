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

Register in https://developer.spotify.com/dashboard.
Then add your client id and secret to a new `.env` file.

Example:

```
SPOTIFY_ID=<id>
SPOTIFY_SECRET=<secret>
```

![alt text](/assets/image.png)


### Youtube

_**Note**: Currently not working!!_

### Environment

Run
``` bash
export $(cat .env | xargs)
```

## Run the app

Migrate playlist

``` bash
go run . playlist <your_playlist_name>
```

Migrate album
``` bash
go run . album <album_name>
```

And click the link in terminal to authenticate Spotify

## Contributing

Please don't commit to *master* branch. Thank you.
