package ports

import (
	"context"
	"spotify_migration/domain"
)

type IExtractor interface {
	Extract(ctx context.Context, resourceName string) (*domain.Collection, error)
}

type IImporter interface {
	Import(ctx context.Context, collection *domain.Collection) (bool, error)
}

type IImportingStrategy interface {
	UpdateItems(ctx context.Context, resourceName string, collectionID string, itemIDs []string) error
}

type ISearchStrategy interface {
	SearchItem(ctx context.Context, music *domain.Music) (itemID string, err error)
}
