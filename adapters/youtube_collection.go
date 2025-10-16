package adapters

import (
	"context"
	"spotify_migration/ports"

	"google.golang.org/api/youtube/v3"
)

type youtubeCollection struct {
	service *youtube.Service
}

func NewYoutubeCollection(service *youtube.Service) ports.ITargetCollection {
	return &youtubeCollection{service: service}
}

func (s *youtubeCollection) CheckIfCollectionExists(ctx context.Context, playlistName string) (collectionID string, err error) {
	response, err := s.service.Playlists.List([]string{"id", "snippet"}).Mine(true).Context(ctx).Do()
	if err != nil {
		return "", err
	}

	for _, playlist := range response.Items {
		if playlist.Snippet.Title == playlistName {
			return playlist.Id, nil
		}
	}
	return "", nil
}

func (s *youtubeCollection) CreateCollection(ctx context.Context, name string) (collectionID string, err error) {
	playlist, err := s.service.Playlists.Insert([]string{"snippet,status"}, &youtube.Playlist{
		Snippet: &youtube.PlaylistSnippet{
			Title: name,
		},
		Status: &youtube.PlaylistStatus{
			PrivacyStatus: "private",
		},
	}).Context(ctx).Do()

	if err != nil {
		return "", err
	}

	return playlist.Id, nil
}
