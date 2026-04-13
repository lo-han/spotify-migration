package album

import (
	"context"
	"log"
	"spotify_migration/entities"
	domain "spotify_migration/entities"
	"spotify_migration/entities/data"
	"spotify_migration/usecases"
)

func NewExtractor(origin usecases.IGetAlbums) domain.IExtractorUsecase {
	return &albumExtractor{
		origin: origin,
	}
}

type albumExtractor struct {
	origin usecases.IGetAlbums
}

func (s *albumExtractor) Extract(ctx context.Context, resourceName string) (*data.Collection, error) {
	albums, err := s.origin.GetAlbums(ctx)
	if err != nil {
		return nil, err
	}

	if albums != nil {
		log.Println("Found albums:")

		for album := albums.Current(); album != nil; album = album.Next() {
			log.Printf("- %s by %s\n", album.Name, entities.ReadStr(album.Artist))
		}
	}

	return albums, nil
}
