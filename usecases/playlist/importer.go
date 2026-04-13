package playlist

import (
	"context"
	"log"
	"spotify_migration/entities"
	"spotify_migration/entities/data"
	"spotify_migration/usecases"
)

const (
	API_LIMIT = 65
)

func NewImporter(
	searcher usecases.ITargetSearch,
	collection usecases.ITargetCollection,
	targetWriter usecases.ITargetWriter,
	migrationState data.IMigrationStateRepository,
) entities.IImporterUsecase {

	return &playlistImporter{
		searcher:       searcher,
		collection:     collection,
		apiLimit:       API_LIMIT,
		targetWriter:   targetWriter,
		migrationState: migrationState,
	}
}

type playlistImporter struct {
	searcher       usecases.ITargetSearch
	collection     usecases.ITargetCollection
	targetWriter   usecases.ITargetWriter
	apiLimit       int
	searchedItems  int
	migrationState data.IMigrationStateRepository
}

func (s *playlistImporter) Import(ctx context.Context, collection *data.Collection) (bool, error) {
	if collection == nil {
		return false, nil
	}

	collectionID, err := s.getCollectionID(ctx, collection)
	if err != nil {
		return false, err
	}

	pendingItemIDs, migratedItemIDs, err := s.retrieveItems()
	if err != nil {
		return false, err
	}

	err = s.getNewItems(ctx, collection, pendingItemIDs, migratedItemIDs)
	if err != nil {
		return false, err
	}

	log.Println("Found", len(pendingItemIDs), "items to import in collection", collection.Name)
	log.Println("Importing items...")

	err = s.insertAll(ctx, collectionID, pendingItemIDs)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *playlistImporter) getCollectionID(ctx context.Context, collection *data.Collection) (string, error) {
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

func (s *playlistImporter) retrieveItems() (pendingItems map[string]string, migratedItems map[string]string, err error) {
	pendingItemIDs := make(map[string]string)
	migratedItemIDs := make(map[string]string)

	exists, err := s.migrationState.Read()
	if err != nil {
		return nil, nil, err
	}
	if exists {
		pendingItemIDs = s.migrationState.GetPendingItems()
		migratedItemIDs = s.migrationState.GetMigratedItems()
	}

	return pendingItemIDs, migratedItemIDs, nil
}

func (s *playlistImporter) getNewItems(ctx context.Context, collection *data.Collection, pendingItems, migratedItems map[string]string) error {
	defer s.migrationState.Save()

	for _, music := range collection.Musics {
		if s.searchedItems >= s.apiLimit {
			break
		}

		id := entities.ID(music)

		if _, exists := migratedItems[id]; exists {
			continue
		}

		itemAddress, err := s.searcher.SearchItem(ctx, music)
		if err != nil {
			return err
		}

		pendingItems[id] = itemAddress
		s.migrationState.AddItem(music, itemAddress)

		s.searchedItems++
	}
	return nil
}

func (s *playlistImporter) insertAll(ctx context.Context, collectionID string, itemIDs map[string]string) error {
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
