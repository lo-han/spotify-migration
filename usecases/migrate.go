package usecases

import "spotify_migration/ports"

type Migration struct {
	Extractor ports.IExtractor
	Importer  ports.IImporter
}

func NewMigration(extractor ports.IExtractor, importer ports.IImporter) *Migration {
	return &Migration{
		Extractor: extractor,
		Importer:  importer,
	}
}

func (m *Migration) Migrate(resourceName string) (bool, error) {
	resourceData, err := m.Extractor.Extract(resourceName)
	if err != nil {
		return false, err
	}

	return m.Importer.Import(resourceData)
}
