package adapters

import (
	"context"
	"log"
	"spotify_migration/entities"
	"spotify_migration/entities/data"
	"spotify_migration/usecases"

	"google.golang.org/api/youtube/v3"
)

type youtubeAlbumSave struct {
	service *youtube.Service
}

func NewYoutubeAlbumSave(service *youtube.Service) usecases.ISaveAlbums {
	return &youtubeAlbumSave{
		service: service,
	}
}

func (s *youtubeAlbumSave) SaveAlbums(ctx context.Context, albums []*data.Album) error {
	albumIDs, err := s.searchAlbums(ctx, albums)
	if err != nil {
		return err
	}

	for range albumIDs {
		// TODO: Save playlist (??????????)
	}
	return nil
}

func (s *youtubeAlbumSave) searchAlbums(ctx context.Context, albums []*data.Album) (albumsID []string, err error) {
	for _, album := range albums {
		call := s.service.Search.List([]string{"id", "snippet"}).
			Q(album.Title + " " + entities.List(album.Artists) + " album oficial channel").MaxResults(1).
			Type("playlist").Context(ctx)

		response, err := call.Context(ctx).Do()
		if err != nil {
			return nil, err
		}

		if len(response.Items) == 0 {
			log.Printf("No YouTube results found for album: %s\n", album.Title)
			continue
		}

		if response.Items[0].Id.Kind == "youtube#playlist" {
			albumsID = append(albumsID, response.Items[0].Id.PlaylistId)
		}
	}
	return albumsID, nil
}
