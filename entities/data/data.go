package data

const (
	PlaylistKind = "playlist"
	AlbumKind    = "albums"
)

type Collection struct {
	Name   string
	Artist *string
	Musics []*Music
	next   *Collection
}

type Music struct {
	Name   string
	Artist string
	Album  *string
}

func (c *Collection) Append(collection *Collection) *Collection {
	if c == nil {
		return nil
	}
	c.next = collection
	return c.next
}

func (c *Collection) Next() *Collection {
	return c.next
}

func (c *Collection) Current() *Collection {
	return c
}
