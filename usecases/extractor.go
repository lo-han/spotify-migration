package usecases

import (
	"context"
	"log"
	"spotify_migration/domain"
)

func NewPlaylistExtractor(ctx context.Context, origin ISourceGetter) domain.IExtractorUsecase {
	return &spotifyPlaylistExtractor{
		origin: origin,
	}
}

type spotifyPlaylistExtractor struct {
	origin ISourceGetter
}

func (s *spotifyPlaylistExtractor) Extract(ctx context.Context, resourceName string) (*domain.Collection, error) {
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
