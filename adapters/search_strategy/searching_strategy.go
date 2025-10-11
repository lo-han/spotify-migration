package searching_strategy

import (
	"context"
	"errors"
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
	call := s.service.Search.List([]string{"id", "snippet"}).Q(s.buildSearchQuery(music)).MaxResults(1).
		Type("video").Context(ctx)

	response, err := call.Context(ctx).Do()
	if err != nil {
		return "", err
	}

	if len(response.Items) == 0 {
		return "", errors.New("item " + music.Title + " not found")
	}

	itemID = response.Items[0].Id.VideoId
	return itemID, nil
}

func (s *StandardSearchStrategy) buildSearchQuery(music *domain.Music) string {
	return music.Title + " " + music.Artist + " " + music.Album + " audio"
}
