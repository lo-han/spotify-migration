package usecases

import (
	"context"
	"log"
	"spotify_migration/entities"
	"spotify_migration/entities/data"
)

const (
	API_LIMIT = 65
)

func NewImporter(
	searcher ITargetSearch, collection ITargetCollection, targetWriter ITargetWriter, migrationState entities.IMigrationStateRepository,
) entities.IImporterUsecase {

	return &youtubeImporter{
		searcher:       searcher,
		collection:     collection,
		apiLimit:       API_LIMIT,
		targetWriter:   targetWriter,
		migrationState: migrationState,
	}
}

type youtubeImporter struct {
	searcher       ITargetSearch
	collection     ITargetCollection
	targetWriter   ITargetWriter
	apiLimit       int
	searchedItems  int
	migrationState entities.IMigrationStateRepository
}

func (s *youtubeImporter) Import(ctx context.Context, collection *data.Collection) (bool, error) {
	if collection == nil {
		return false, nil
	}

	collectionID, err := s.getCollectionID(ctx, collection)
	if err != nil {
		return false, err
	}

	itemIDs, err := s.retrievePendingItems()
	if err != nil {
		return false, err
	}

	err = s.getNewItems(ctx, collection, itemIDs)
	if err != nil {
		return false, err
	}

	log.Println("Found", len(itemIDs), "items to import in collection", collection.Name)
	log.Println("Importing items...")

	err = s.insertAll(ctx, collectionID, itemIDs)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *youtubeImporter) getCollectionID(ctx context.Context, collection *data.Collection) (string, error) {
	collectionID, err := s.collection.CheckIfCollectionExists(ctx, collection.Name)
	if err != nil {
		return "", err
	}

	if collectionID == "" {
		log.Println("Collection", collection.Name, "does not exist. Creating it...")

		collectionID, err = s.collection.CreateCollection(ctx, collection.Name)
		if err != nil {
			return "", err
		}
	}
	return collectionID, nil
}

func (s *youtubeImporter) retrievePendingItems() (map[string]string, error) {
	itemIDs := make(map[string]string)

	exists, err := s.migrationState.Read()
	if err != nil {
		return nil, err
	}
	if exists {
		itemIDs = s.migrationState.GetPendingItems()
	}

	return itemIDs, nil
}

func (s *youtubeImporter) getNewItems(ctx context.Context, collection *data.Collection, itemIDs map[string]string) error {
	defer s.migrationState.Save()

	for _, music := range collection.Musics {
		if s.searchedItems >= s.apiLimit {
			break
		}

		id := entities.ID(music)

		if _, exists := itemIDs[id]; exists {
			continue
		}

		itemAddress, err := s.searcher.SearchItem(ctx, music)
		if err != nil {
			return err
		}

		itemIDs[id] = itemAddress
		s.migrationState.AddItem(music, itemAddress)

		s.searchedItems++
	}
	return nil
}

func (s *youtubeImporter) insertAll(ctx context.Context, collectionID string, itemIDs map[string]string) error {
	if collectionID == "" || len(itemIDs) == 0 {
		return nil
	}
	defer s.migrationState.Save()

	var insertedItems = 0

	for itemID, itemAddress := range itemIDs {
		err := s.targetWriter.AddItemToPlaylist(ctx, collectionID, itemAddress)
		if err != nil {
			return err
		}
		s.migrationState.UpdateItemToMigrated(itemID)
		insertedItems++

		if insertedItems >= s.apiLimit {
			break
		}
	}

	return nil
}
