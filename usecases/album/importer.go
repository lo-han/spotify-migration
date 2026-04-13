package album

import (
	"context"
	"log"
	"spotify_migration/entities"
	"spotify_migration/entities/data"
	"spotify_migration/usecases"
)

func NewImporter(
	target usecases.ISaveAlbums,
) entities.IImporterUsecase {
	return &playlistImporter{
		target: target,
	}
}

type playlistImporter struct {
	target usecases.ISaveAlbums
}

func (s *playlistImporter) Import(ctx context.Context, collection *data.Collection) (bool, error) {
	if collection == nil {
		log.Println("No items to import.")
		return false, nil
	}

	log.Println("Importing items...")

	err := s.target.SaveAlbums(ctx, collection)

	return err == nil, err
}
