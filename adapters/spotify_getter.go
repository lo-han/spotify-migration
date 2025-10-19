package adapters

import (
	"context"
	"errors"
	"spotify_migration/entities/data"
	"spotify_migration/usecases"

	"github.com/lo-han/spotify/v2"
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

	collection = &data.Collection{
		Name: resourceName,
	}

	itemsPage, err := s.client.GetPlaylistItems(ctx, spotify.ID(id))
	if err != nil {
		return nil, err
	}

	for ; itemsPage != nil; itemsPage, err = itemsPage.NextItems(ctx, s.client) {

		if err != nil {
			return nil, err
		}

		for _, item := range itemsPage.Items {
			track := item.Track.Track.SimpleTrack

			music := &data.Music{
				Title: track.Name,
				Album: track.Album.Name,
			}
			if len(track.Artists) > 0 {
				music.Artist = track.Artists[0].Name
			}

			collection.Musics = append(collection.Musics, music)
		}
	}

	return collection, nil
}
