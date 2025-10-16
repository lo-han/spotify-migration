package domain

import (
	"context"
)

type IMigrate interface {
	Migrate(ctx context.Context, resourceName string, migrationState IMigrationStateRepository) (bool, error)
}

type IExtractorUsecase interface {
	Extract(ctx context.Context, resourceName string) (*Collection, error)
}

type IImporterUsecase interface {
	Import(ctx context.Context, collection *Collection, migrationState IMigrationStateRepository) (bool, error)
}

type IMigrationStateRepository interface {
	GetPendingItems() []string
	UpdateItemToMigrated(itemID string)
	Read() (bool, error)
	Save() error
}
