package domain

import (
	"context"
)

type IMigrate interface {
	Migrate(ctx context.Context, resourceName string) (bool, error)
}

type IExtractorUsecase interface {
	Extract(ctx context.Context, resourceName string) (*Collection, error)
}

type IImporterUsecase interface {
	Import(ctx context.Context, collection *Collection) (bool, error)
}

type IMigrationStateRepository interface {
	GetPendingItems() map[string]string
	UpdateItemToMigrated(itemID string)
	AddItem(item *Music, address string)
	Read() (bool, error)
	Save() error
}
