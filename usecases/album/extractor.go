package album

import (
	"context"
	"log"
	"spotify_migration/entities"
	domain "spotify_migration/entities"
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

func (s *albumExtractor) Extract(ctx context.Context, resourceName string) (any, error) {
	albums, err := s.origin.GetAlbums(ctx)
	if err != nil {
		return nil, err
	}

	if len(albums) != 0 {
		log.Println("Found albums:")
		for _, album := range albums {
			log.Printf("- %s by %s\n", album.Title, entities.List(album.Artists))
		}
	}
	return albums, nil
}
