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
	target usecases.ISaveAlbums,
) entities.IImporterUsecase {
	return &playlistImporter{
		target: target,
	}
}

type playlistImporter struct {
	target usecases.ISaveAlbums
}

func (s *playlistImporter) Import(ctx context.Context, collection any) (bool, error) {
	if collection == nil {
		return false, nil
	}

	albums, ok := collection.([]*data.Album)
	if !ok {
		return false, errors.New("invalid album data")
	}

	if len(albums) == 0 {
		log.Println("No items to import.")
		return false, nil
	}

	log.Println("Importing items...")

	err := s.target.SaveAlbums(ctx, albums)

	return err == nil, err
}
