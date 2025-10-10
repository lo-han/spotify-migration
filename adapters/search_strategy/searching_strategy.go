package searching_strategy

import (
	"context"
	"spotify_migration/domain"
)

type StandardSearchStrategy struct {
}

func NewStandardSearchStrategy() *StandardSearchStrategy {
	return &StandardSearchStrategy{}
}

func (s *StandardSearchStrategy) SearchItem(ctx context.Context, music *domain.Music) (itemID string, err error) {
	// Simulate searching for an item
	if music == nil || music.Title == "" {
		return "", nil
	}
	return "found_item_id", nil
}
