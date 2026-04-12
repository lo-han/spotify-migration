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
			log.Printf("- %s by %s\n", album.Title, s.listArtists(album.Artists))
		}
	}
	return albuns, nil
}

func (s *albumExtractor) listArtists(artists []string) string {
	if len(artists) == 0 {
		return ""
	}
	if len(artists) == 1 {
		return artists[0]
	}
	result := artists[0]
	for i := 1; i < len(artists); i++ {
		result += ", " + artists[i]
	}
	return result
}
