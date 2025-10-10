package searching_strategy

import (
	"context"
	"spotify_migration/domain"

	"google.golang.org/api/youtube/v3"
)

type StandardSearchStrategy struct {
	service *youtube.Service
}

func NewStandardSearchStrategy(service *youtube.Service) *StandardSearchStrategy {
	return &StandardSearchStrategy{service: service}
}

func (s *StandardSearchStrategy) SearchItem(ctx context.Context, music *domain.Music) (itemID string, err error) {
	// Simulate searching for an item
	if music == nil || music.Title == "" {
		return "", nil
	}
	return "found_item_id", nil
}
