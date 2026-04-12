package adapters

import (
	"context"
	"spotify_migration/entities/data"
	"spotify_migration/usecases"

	"github.com/lo-han/spotify/v2"
)

type spotifyAlbumSearch struct {
	client *spotify.Client
}

func NewSpotifyAlbumSearch(client *spotify.Client) usecases.IGetAlbuns {
	return &spotifyAlbumSearch{client: client}
}

func (s *spotifyAlbumSearch) GetAlbuns(ctx context.Context) (albums []*data.Album, err error) {
	albumPage, err := s.client.CurrentUsersAlbums(ctx)
	if err != nil {
		return nil, err
	}

	albums = make([]*data.Album, len(albumPage.Albums))

	// O(n*m) where n is the number of albums and m is the number of artists per album
	for index, album := range albumPage.Albums {
		albums[index] = &data.Album{
			Title: album.Name,
		}
		albums[index].Artists = make([]string, len(album.Artists))

		if len(album.Artists) > 0 {
			for i, artist := range album.Artists {
				albums[index].Artists[i] = artist.Name
			}
		}
	}
	return albums, nil
}
