package data

const (
	PlaylistKind = "playlist"
	AlbumKind    = "album"
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
