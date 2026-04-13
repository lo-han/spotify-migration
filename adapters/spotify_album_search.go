package adapters

import (
	"context"
	"spotify_migration/entities"
	"spotify_migration/entities/data"
	"spotify_migration/usecases"
	"strings"

	"github.com/lo-han/spotify/v2"
)

type spotifyAlbumSearch struct {
	client *spotify.Client
}

func NewSpotifyAlbumSearch(client *spotify.Client) usecases.IGetAlbums {
	return &spotifyAlbumSearch{client: client}
}

func (s *spotifyAlbumSearch) GetAlbums(ctx context.Context) (albums *data.Collection, err error) {
	albumPage, err := s.client.CurrentUsersAlbums(ctx)
	if err != nil {
		return nil, err
	}

	if len(albumPage.Albums) == 0 {
		return nil, nil
	}

	currentAlbum := &data.Collection{
		Name: albumPage.Albums[0].Name,
	}
	currentAlbum.Artist = s.getArtists(albumPage.Albums[0])
	albums = currentAlbum

	// O(n*m) where n is the number of albums and m is the number of artists per album
	for i := 1; i < len(albumPage.Albums); i++ {
		currentPage := albumPage.Albums[i]

		album := &data.Collection{
			Name: currentPage.Name,
		}
		album.Artist = s.getArtists(currentPage)

		currentAlbum.Append(album)
		currentAlbum = currentAlbum.Next()
	}
	return albums, nil
}

func (s *spotifyAlbumSearch) getArtists(album spotify.SavedAlbum) *string {
	artists := make([]string, len(album.Artists))

	if len(album.Artists) == 0 {
		return entities.PtrStr("")
	}

	for i, artist := range album.Artists {
		artists[i] = artist.Name
	}
	return entities.PtrStr(strings.Join(artists, ", "))
}
