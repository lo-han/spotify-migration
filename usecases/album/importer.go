package album

import (
	"context"
	"errors"
	"log"
	"spotify_migration/entities"
	"spotify_migration/entities/data"
	"spotify_migration/usecases"
)

func NewImporter(
	target usecases.ISaveAlbuns,
) entities.IImporterUsecase {
	return &playlistImporter{
		target: target,
	}
}

type playlistImporter struct {
	target usecases.ISaveAlbuns
}

func (s *playlistImporter) Import(ctx context.Context, collection any) (bool, error) {
	if collection == nil {
		return false, nil
	}

	albuns, ok := collection.([]*data.Album)
	if !ok {
		return false, errors.New("invalid album data")
	}

	if len(albuns) == 0 {
		log.Println("No items to import.")
		return false, nil
	}

	log.Println("Importing items...")

	err := s.target.SaveAlbuns(ctx, albuns)

	return err == nil, err
}
