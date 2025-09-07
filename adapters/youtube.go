package adapters

import (
	"spotify_migration/domain"
	"spotify_migration/ports"
)

func NewYoutubeImporter() *YoutubeImporter {
	return &YoutubeImporter{
		updater:  NewYoutubeMemsetUpdater(),
		searcher: NewStandardSearchStrategy(),
	}
}

type YoutubeImporter struct {
	updater  ports.IImportingStrategy
	searcher ports.ISearchStrategy
}

func (s *YoutubeImporter) Import(collection *domain.Collection) (bool, error) {
	if collection == nil {
		return false, nil
	}

	collectionID, err := s.checkIfPCollectionExists(collection.Name)
	if err != nil {
		return false, err
	}

	if collectionID == "" {
		collectionID, err = s.createCollection(collection.Name)
		if err != nil {
			return false, err
		}
	}

	var itemIDs []string

	for _, music := range collection.Musics {
		itemID, err := s.searcher.SearchItem(music)
		if err != nil {
			return false, err
		}
		itemIDs = append(itemIDs, itemID)
	}

	err = s.updater.UpdateItems(collectionID, itemIDs)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *YoutubeImporter) checkIfPCollectionExists(playlistName string) (collectionID string, err error) {
	// Simulate a check for the playlist's existence
	if playlistName == "" {
		return "", nil
	}
	return "existing_playlist_id", nil
}

func (s *YoutubeImporter) createCollection(name string) (collectionID string, err error) {
	// Simulate playlist creation
	if name == "" {
		return "", nil
	}
	return "new_playlist_id", nil
}
