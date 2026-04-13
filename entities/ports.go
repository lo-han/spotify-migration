package entities

import (
	"context"
	"spotify_migration/entities/data"
)

type IExtractorUsecase interface {
	Extract(ctx context.Context, resourceName string) (*data.Collection, error)
}

type IImporterUsecase interface {
	Import(ctx context.Context, collection *data.Collection) (bool, error)
}
