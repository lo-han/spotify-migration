package entities

import "spotify_migration/entities/data"

func ID(music *data.Music) string {
	return music.Title + "_" + music.Artist + "_" + music.Album
}
