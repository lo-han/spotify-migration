package usecases

import (
	"context"
	"spotify_migration/entities/data"
)

type IGetAlbums interface {
	GetAlbums(ctx context.Context) (albums *data.Collection, err error)
}

type ISaveAlbums interface {
	SaveAlbums(ctx context.Context, albums *data.Collection) error
}

type ISourceGetter interface {
	GetPlaylistID(ctx context.Context, resourceName string) (string, error)
	GetPlaylistItems(ctx context.Context, resourceName, id string) (collection *data.Collection, err error)
}

type ITargetSearch interface {
	SearchItem(ctx context.Context, music *data.Music) (itemID string, err error)
}

type ITargetCollection interface {
	CheckIfCollectionExists(ctx context.Context, playlistName string) (collectionID string, err error)
	CreateCollection(ctx context.Context, name string) (collectionID string, err error)
}

type ITargetWriter interface {
	AddItemToPlaylist(ctx context.Context, collectionID string, itemID string) error
}
