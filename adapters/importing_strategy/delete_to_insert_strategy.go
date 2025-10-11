package importing_strategy

import (
	"context"
	"spotify_migration/ports"

	"google.golang.org/api/youtube/v3"
)

func NewYoutubeDeleteToInsert(service *youtube.Service) ports.IImportingStrategy {
	return &youtubeDeleteToInsert{service: service}
}

type youtubeDeleteToInsert struct {
	service *youtube.Service
	*youtubeJustInsert
}

func (u *youtubeDeleteToInsert) UpdateItems(ctx context.Context, resourceName string, collectionID string, itemIDs []string) error {
	err := u.deletePlaylist(ctx, collectionID)
	if err != nil {
		return err
	}

	collectionID, err = u.recreatePlaylist(ctx, resourceName)
	if err != nil {
		return err
	}

	err = u.insertAll(ctx, collectionID, itemIDs)
	if err != nil {
		return err
	}

	return nil
}

func (u *youtubeDeleteToInsert) deletePlaylist(ctx context.Context, collectionID string) error {
	return u.service.Playlists.Delete(collectionID).Context(ctx).Do()
}

func (u *youtubeDeleteToInsert) recreatePlaylist(ctx context.Context, collectionName string) (newID string, err error) {
	playlist, err := u.service.Playlists.Insert([]string{"snippet,status"}, &youtube.Playlist{
		Snippet: &youtube.PlaylistSnippet{
			Title: collectionName,
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
