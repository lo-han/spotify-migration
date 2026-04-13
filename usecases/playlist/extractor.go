package playlist

import (
	"context"
	"log"
	domain "spotify_migration/entities"
	"spotify_migration/entities/data"
	"spotify_migration/usecases"
)

func NewExtractor(origin usecases.ISourceGetter) domain.IExtractorUsecase {
	return &playlistExtractor{
		origin: origin,
	}
}

type playlistExtractor struct {
	origin usecases.ISourceGetter
}

func (s *playlistExtractor) Extract(ctx context.Context, resourceName string) (*data.Collection, error) {
	playlistID, err := s.origin.GetPlaylistID(ctx, resourceName)
	if err != nil {
		return nil, err
	}
	log.Println("Found playlist", resourceName, "with ID:", playlistID)

	playlistItems, err := s.origin.GetPlaylistItems(ctx, resourceName, playlistID)
	if err != nil {
		return nil, err
	}
	log.Println("Found", len(playlistItems.Musics), "items in playlist", resourceName)

	return playlistItems, nil
}
