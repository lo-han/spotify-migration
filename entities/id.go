package entities

import "spotify_migration/entities/data"

func ID(music *data.Music) string {
	if music == nil {
		return ""
	}
	return music.Title + "_" + music.Artist + "_" + music.Album
}
