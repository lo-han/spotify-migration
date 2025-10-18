package usecases

import (
	"context"
	"spotify_migration/domain"
)

type ISourceGetter interface {
	GetPlaylistID(ctx context.Context, resourceName string) (string, error)
	GetPlaylistItems(ctx context.Context, resourceName, id string) (collection *domain.Collection, err error)
}

type ITargetSearch interface {
	SearchItem(ctx context.Context, music *domain.Music) (itemID string, err error)
}

type ITargetCollection interface {
	CheckIfCollectionExists(ctx context.Context, playlistName string) (collectionID string, err error)
	CreateCollection(ctx context.Context, name string) (collectionID string, err error)
}

type ITargetWriter interface {
	AddItemToPlaylist(ctx context.Context, collectionID string, itemID string) error
}
