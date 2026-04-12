package data

const (
	PlaylistKind = "playlist"
	AlbumKind    = "albums"
)

type Album struct {
	Title  string
	Artist string
}

type Collection struct {
	Name   string
	Musics []*Music
}

type Music struct {
	Title  string
	Artist string
	Album  string
}
