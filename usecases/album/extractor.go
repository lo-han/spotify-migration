package album

import (
	"context"
	"log"
	domain "spotify_migration/entities"
	"spotify_migration/usecases"
)

func NewExtractor(origin usecases.IGetAlbuns) domain.IExtractorUsecase {
	return &albumExtractor{
		origin: origin,
	}
}

type albumExtractor struct {
	origin usecases.IGetAlbuns
}

func (s *albumExtractor) Extract(ctx context.Context, resourceName string) (any, error) {
	albuns, err := s.origin.GetAlbuns(ctx)
	if err != nil {
		return nil, err
	}

	if len(albuns) != 0 {
		log.Println("Found albuns:")
		for _, album := range albuns {
			log.Printf("- %s by %s\n", album.Title, album.Artist)
		}
	}
	return albuns, nil
}
