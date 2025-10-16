package adapters

import (
	"context"

	"google.golang.org/api/youtube/v3"
)

type youtubeCollectionWriter struct {
	service *youtube.Service
}

func NewYoutubeCollectionWriter(service *youtube.Service) *youtubeCollectionWriter {
	return &youtubeCollectionWriter{service: service}
}

func (u *youtubeCollectionWriter) AddItemToPlaylist(ctx context.Context, collectionID string, itemID string) error {
	_, err := u.service.PlaylistItems.Insert([]string{"snippet"}, &youtube.PlaylistItem{
		Snippet: &youtube.PlaylistItemSnippet{
			PlaylistId: collectionID,
			ResourceId: &youtube.ResourceId{
				Kind:    "youtube#video",
				VideoId: itemID,
			},
		},
	}).Context(ctx).Do()

	if err != nil {
		return err
	}
	return nil
}
