package domain

import (
	"context"
)

type Migration struct {
	Extractor IExtractorUsecase
	Importer  IImporterUsecase
}

func NewMigration(extractor IExtractorUsecase, importer IImporterUsecase) IMigrate {
	return &Migration{
		Extractor: extractor,
		Importer:  importer,
	}
}

func (m *Migration) Migrate(ctx context.Context, resourceName string, migrationState IMigrationStateRepository) (bool, error) {
	resourceData, err := m.Extractor.Extract(ctx, resourceName)
	if err != nil {
		return false, err
	}

	return m.Importer.Import(ctx, resourceData, migrationState)
}
