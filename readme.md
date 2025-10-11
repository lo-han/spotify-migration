# Spotify Migration

Automation to migrate playlist or albums from Spotify to Youtube

_**Note**: Currently album migration is not supported_

## Environment setting

### Spotify

After register in https://developer.spotify.com/dashboard, you should add your client id and secret in a new `.env` file.

Example:

```
SPOTIFY_ID=<id>
SPOTIFY_SECRET=<secret>
```

### Youtube

After creating a service account key in https://cloud.google.com/iam/docs/keys-create-delete, you should add it to a new `keyfile.json`, in spotify-migration root directory.


## Run the app

Migrate playlist

``` bash
go run . playlist <your_playlist_name>
```

Migrate album
``` bash
go run . album <album_name>
```

## Contributing

Please don't commit to *master* branch. Thank you.
