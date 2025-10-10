package importing_strategy

import (
	"context"

	"google.golang.org/api/youtube/v3"
)

func NewYoutubeMemsetUpdater(service *youtube.Service) *YoutubeMemsetUpdater {
	return &YoutubeMemsetUpdater{service: service}
}

type YoutubeMemsetUpdater struct {
	service *youtube.Service
}

func (u *YoutubeMemsetUpdater) UpdateItems(ctx context.Context, collectionID string, itemIDs []string) error {
	err := u.deleteAll(ctx, collectionID)
	if err != nil {
		return err
	}

	err = u.insertAll(ctx, collectionID, itemIDs)
	if err != nil {
		return err
	}

	return nil
}

func (u *YoutubeMemsetUpdater) deleteAll(ctx context.Context, collectionID string) error {
	// Simulate deleting all items from the playlist
	if collectionID == "" {
		return nil
	}
	return nil
}

func (u *YoutubeMemsetUpdater) insertAll(ctx context.Context, collectionID string, itemIDs []string) error {
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

func (u *YoutubeMemsetUpdater) addItemToPlaylist(ctx context.Context, collectionID string, itemID string) error {
	// Simulate adding an item to the playlist
	if collectionID == "" || itemID == "" {
		return nil
	}
	return nil
}
