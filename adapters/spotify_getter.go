package adapters

import (
	"context"
	"errors"
	"fmt"
	"spotify_migration/entities/data"
	"spotify_migration/usecases"

	"github.com/zmb3/spotify/v2"
)

type spotifyGetter struct {
	client *spotify.Client
}

func NewSpotifyGetter(client *spotify.Client) usecases.ISourceGetter {
	return &spotifyGetter{client: client}
}

func (s *spotifyGetter) GetPlaylistID(ctx context.Context, resourceName string) (string, error) {
	playlistPage, err := s.client.CurrentUsersPlaylists(ctx)
	if err != nil {
		return "", err
	}

	for _, playlist := range playlistPage.Playlists {
		if playlist.Name == resourceName {
			return playlist.ID.String(), nil
		}
	}
	return "", errors.New("playlist not found")
}

func (s *spotifyGetter) GetPlaylistItems(ctx context.Context, resourceName, id string) (collection *data.Collection, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()

	collection = &data.Collection{
		Name: resourceName,
	}

	itensPage, err := s.client.GetPlaylistItems(ctx, spotify.ID(id))
	if err != nil {
		return nil, err
	}

	for _, item := range itensPage.Items {
		track := item.Track.Track.SimpleTrack

		collection.Musics = append(collection.Musics, &data.Music{
			Title:  track.Name,
			Artist: track.Artists[0].Name,
			Album:  track.Album.Name,
		})
	}

	return collection, nil
}
