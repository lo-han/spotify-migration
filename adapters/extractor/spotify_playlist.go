package extractor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"spotify_migration/domain"
	"spotify_migration/ports"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

func NewSpotifyPlaylistExtractor(ctx context.Context, auth *spotifyauth.Authenticator, token *oauth2.Token) ports.IExtractor {
	return &spotifyPlaylistExtractor{
		client: spotify.New(auth.Client(ctx, token)),
	}
}

type spotifyPlaylistExtractor struct {
	client *spotify.Client
}

func (s *spotifyPlaylistExtractor) Extract(ctx context.Context, resourceName string) (*domain.Collection, error) {
	playlistID, err := s.getPlaylistID(ctx, resourceName)
	if err != nil {
		return nil, err
	}
	log.Println("Found playlist", resourceName, "with ID:", playlistID)

	playlistItems, err := s.getPlaylistItems(ctx, resourceName, playlistID)
	if err != nil {
		return nil, err
	}
	log.Println("Found", len(playlistItems.Musics), "items in playlist", resourceName)

	return playlistItems, nil
}

func (s *spotifyPlaylistExtractor) getPlaylistID(ctx context.Context, resourceName string) (string, error) {
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

func (s *spotifyPlaylistExtractor) getPlaylistItems(ctx context.Context, resourceName, id string) (collection *domain.Collection, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()

	collection = &domain.Collection{
		Name: resourceName,
	}

	itensPage, err := s.client.GetPlaylistItems(ctx, spotify.ID(id))
	if err != nil {
		return nil, err
	}

	for _, item := range itensPage.Items {
		track := item.Track.Track.SimpleTrack

		collection.Musics = append(collection.Musics, &domain.Music{
			Title:  track.Name,
			Artist: track.Artists[0].Name,
			Album:  track.Album.Name,
		})
	}

	return collection, nil
}
