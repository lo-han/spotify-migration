package importer

import (
	"context"
	"spotify_migration/adapters/importing_strategy"
	searching_strategy "spotify_migration/adapters/search_strategy"
	"spotify_migration/domain"
	"spotify_migration/ports"

	"google.golang.org/api/youtube/v3"
)

func NewYoutubeImporter(service *youtube.Service) ports.IImporter {
	return &youtubeImporter{
		searcher: searching_strategy.NewStandardSearchStrategy(service),
		service:  service,
	}
}

type youtubeImporter struct {
	service  *youtube.Service
	updater  ports.IImportingStrategy
	searcher ports.ISearchStrategy
}

func (s *youtubeImporter) Import(ctx context.Context, collection *domain.Collection) (bool, error) {
	if collection == nil {
		return false, nil
	}

	collectionID, err := s.checkIfCollectionExists(ctx, collection.Name)
	if err != nil {
		return false, err
	}

	if collectionID == "" {
		collectionID, err = s.createCollection(ctx, collection.Name)
		if err != nil {
			return false, err
		}
		s.updater = importing_strategy.NewYoutubeJustInsert(s.service)
	} else {
		s.updater = importing_strategy.NewYoutubeDeleteToInsert(s.service)
	}

	var itemIDs []string

	for _, music := range collection.Musics {
		itemID, err := s.searcher.SearchItem(ctx, music)
		if err != nil {
			return false, err
		}
		itemIDs = append(itemIDs, itemID)
	}

	err = s.updater.UpdateItems(ctx, collection.Name, collectionID, itemIDs)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *youtubeImporter) checkIfCollectionExists(ctx context.Context, playlistName string) (collectionID string, err error) {
	response, err := s.service.Playlists.List([]string{"id", "snippet"}).Context(ctx).Do()
	if err != nil {
		return "", err
	}

	for _, playlist := range response.Items {
		if playlist.Snippet.Title == playlistName {
			return playlist.Id, nil
		}
	}
	return "", nil
}

func (s *youtubeImporter) createCollection(ctx context.Context, name string) (collectionID string, err error) {
	playlist, err := s.service.Playlists.Insert([]string{"snippet,status"}, &youtube.Playlist{
		Snippet: &youtube.PlaylistSnippet{
			Title: name,
		},
		Status: &youtube.PlaylistStatus{
			PrivacyStatus: "private",
		},
	}).Context(ctx).Do()

	if err != nil {
		return "", err
	}

	return playlist.Id, nil
}
