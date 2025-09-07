package ports

import "spotify_migration/domain"

type IExtractor interface {
	Extract(resourceName string) (*domain.Collection, error)
}

type IImporter interface {
	Import(*domain.Collection) (bool, error)
}

type IImportingStrategy interface {
	UpdateItems(collectionID string, itemIDs []string) error
}

type ISearchStrategy interface {
	SearchItem(music *domain.Music) (itemID string, err error)
}
