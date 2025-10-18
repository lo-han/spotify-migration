package domain

func ID(music *Music) string {
	return music.Title + "_" + music.Artist + "_" + music.Album
}
