package importing_strategy

import (
	"context"
	"spotify_migration/ports"

	"google.golang.org/api/youtube/v3"
)

func NewYoutubeJustInsert(service *youtube.Service) ports.IImportingStrategy {
	return &youtubeJustInsert{service: service}
}

type youtubeJustInsert struct {
	service *youtube.Service
}

func (u *youtubeJustInsert) UpdateItems(ctx context.Context, resourceName string, collectionID string, itemIDs []string) error {
	err := u.insertAll(ctx, collectionID, itemIDs)
	if err != nil {
		return err
	}

	return nil
}

func (u *youtubeJustInsert) insertAll(ctx context.Context, collectionID string, itemIDs []string) error {
	if collectionID == "" || len(itemIDs) == 0 {
		return nil
	}
	for _, itemID := range itemIDs {
		err := u.addItemToPlaylist(ctx, collectionID, itemID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *youtubeJustInsert) addItemToPlaylist(ctx context.Context, collectionID string, itemID string) error {
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
