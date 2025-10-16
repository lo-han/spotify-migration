package domain

const (
	PlaylistKind = "playlist"
	AlbumKind    = "album"
)

type Collection struct {
	Name   string
	Musics []*Music
}

type Music struct {
	Title  string
	Artist string
	Album  string
}
