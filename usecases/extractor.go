package usecases

import (
	"context"
	"log"
	domain "spotify_migration/entities"
	"spotify_migration/entities/data"
)

func NewPlaylistExtractor(origin ISourceGetter) domain.IExtractorUsecase {
	return &extractor{
		origin: origin,
	}
}

type extractor struct {
	origin ISourceGetter
}

func (s *extractor) Extract(ctx context.Context, resourceName string) (*data.Collection, error) {
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
