package entities

import (
	"context"
)

type Migration struct {
	extractor        IExtractorUsecase
	playlistImporter IImporterUsecase
}

func NewMigration(extractor IExtractorUsecase, playlistImporter IImporterUsecase) *Migration {
	return &Migration{
		extractor:        extractor,
		playlistImporter: playlistImporter,
	}
}

func (m *Migration) Migrate(ctx context.Context, resourceName string) (bool, error) {
	resourceData, err := m.extractor.Extract(ctx, resourceName)
	if err != nil {
		return false, err
	}

	succeeded, err := m.playlistImporter.Import(ctx, resourceData)
	if err != nil {
		return false, err
	}

	return succeeded, nil
}
