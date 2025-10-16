package usecases

import (
	"context"
	"log"
	"spotify_migration/domain"
	"spotify_migration/ports"
)

const (
	API_LIMIT = 65
)

func NewImporter(searcher ports.ITargetSearch, collection ports.ITargetCollection, targetWriter ports.ITargetWriter) domain.IImporterUsecase {
	return &youtubeImporter{
		searcher:     searcher,
		collection:   collection,
		apiLimit:     API_LIMIT,
		targetWriter: targetWriter,
	}
}

type youtubeImporter struct {
	searcher     ports.ITargetSearch
	collection   ports.ITargetCollection
	targetWriter ports.ITargetWriter
	apiLimit     int
}

func (s *youtubeImporter) Import(ctx context.Context, collection *domain.Collection, migrationState domain.IMigrationStateRepository) (bool, error) {
	if collection == nil {
		return false, nil
	}

	collectionID, err := s.checkIfCollectionExists(ctx, collection)
	if err != nil {
		return false, err
	}

	var itemIDs []string

	for _, music := range collection.Musics {
		itemID, err := s.searcher.SearchItem(ctx, music)
		if err != nil {
			return false, err
		}
		itemIDs = append(itemIDs, itemID)
	}
	log.Println("Found", len(itemIDs), "items to import in collection", collection.Name)
	log.Println("Importing items...")

	err = s.insertAll(ctx, migrationState, collectionID, itemIDs)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *youtubeImporter) checkIfCollectionExists(ctx context.Context, collection *domain.Collection) (string, error) {
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

func (u *youtubeImporter) insertAll(ctx context.Context, migrationState domain.IMigrationStateRepository, collectionID string, itemIDs []string) error {
	if collectionID == "" || len(itemIDs) == 0 {
		return nil
	}

	exists, err := migrationState.Read()
	if err != nil {
		return err
	}
	if exists {
		itemIDs = migrationState.GetPendingItems()
	}

	for index, itemID := range itemIDs {
		if index == u.apiLimit-1 {
			break
		}

		err := u.targetWriter.AddItemToPlaylist(ctx, collectionID, itemID)
		if err != nil {
			return err
		}
		migrationState.UpdateItemToMigrated(itemID)
	}

	migrationState.Save()

	return nil
}
