package adapters

import (
	"context"
	"errors"
	"spotify_migration/entities/data"

	"google.golang.org/api/youtube/v3"
)

type youtubeSearch struct {
	service *youtube.Service
}

func NewYoutubeSearch(service *youtube.Service) *youtubeSearch {
	return &youtubeSearch{service: service}
}

func (s *youtubeSearch) SearchItem(ctx context.Context, music *data.Music) (itemID string, err error) {
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

func (s *youtubeSearch) buildSearchQuery(music *data.Music) string {
	return music.Title + " " + music.Artist + " " + music.Album + " audio"
}
