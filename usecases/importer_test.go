package usecases

import (
	"context"
	"spotify_migration/entities"
	"spotify_migration/entities/data"
	"spotify_migration/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newImporterTestNewImporter(
	searcher ITargetSearch, collection ITargetCollection, targetWriter ITargetWriter, migrationState entities.IMigrationStateRepository,
) *youtubeImporter {

	return &youtubeImporter{
		searcher:       searcher,
		collection:     collection,
		apiLimit:       API_LIMIT,
		targetWriter:   targetWriter,
		migrationState: migrationState,
	}
}

func TestNewImporter(t *testing.T) {
	mockSearcher := mocks.NewITargetSearch(t)
	mockCollection := mocks.NewITargetCollection(t)
	mockTargetWriter := mocks.NewITargetWriter(t)
	mockMigrationState := mocks.NewIMigrationStateRepository(t)

	importer := NewImporter(mockSearcher, mockCollection, mockTargetWriter, mockMigrationState)

	assert.NotNil(t, importer)
	assert.IsType(t, &youtubeImporter{}, importer)

	youtubeImporter := importer.(*youtubeImporter)
	assert.Equal(t, API_LIMIT, youtubeImporter.apiLimit)
	assert.Equal(t, 0, youtubeImporter.searchedItems)
	assert.NotNil(t, youtubeImporter.searcher)
	assert.NotNil(t, youtubeImporter.collection)
	assert.NotNil(t, youtubeImporter.targetWriter)
	assert.NotNil(t, youtubeImporter.migrationState)
}

func TestYoutubeImporter_Import_Success_ExistingCollection(t *testing.T) {
	ctx := context.Background()
	collection := &data.Collection{
		Name: "Test Playlist",
		Musics: []*data.Music{
			{Title: "Song 1", Artist: "Artist 1", Album: "Album 1"},
			{Title: "Song 2", Artist: "Artist 2", Album: "Album 2"},
		},
	}

	existingCollectionID := "existing_collection_123"
	pendingItems := map[string]string{}

	mockSearcher := mocks.NewITargetSearch(t)
	mockCollection := mocks.NewITargetCollection(t)
	mockTargetWriter := mocks.NewITargetWriter(t)
	mockMigrationState := mocks.NewIMigrationStateRepository(t)

	mockCollection.On("CheckIfCollectionExists", ctx, collection.Name).Return(existingCollectionID, nil)
	mockMigrationState.On("Read").Return(true, nil)
	mockMigrationState.On("GetPendingItems").Return(pendingItems)

	mockSearcher.On("SearchItem", ctx, collection.Musics[0]).Return("youtube_id_1", nil)
	mockSearcher.On("SearchItem", ctx, collection.Musics[1]).Return("youtube_id_2", nil)

	mockMigrationState.On("AddItem", collection.Musics[0], "youtube_id_1").Return()
	mockMigrationState.On("AddItem", collection.Musics[1], "youtube_id_2").Return()

	song1ID := entities.ID(collection.Musics[0])
	song2ID := entities.ID(collection.Musics[1])
	mockTargetWriter.On("AddItemToPlaylist", ctx, existingCollectionID, "youtube_id_1").Return(nil)
	mockTargetWriter.On("AddItemToPlaylist", ctx, existingCollectionID, "youtube_id_2").Return(nil)
	mockMigrationState.On("UpdateItemToMigrated", song1ID).Return()
	mockMigrationState.On("UpdateItemToMigrated", song2ID).Return()
	mockMigrationState.On("Save").Return(nil)

	importer := NewImporter(mockSearcher, mockCollection, mockTargetWriter, mockMigrationState)

	result, err := importer.Import(ctx, collection)

	assert.NoError(t, err)
	assert.True(t, result)
}

func TestYoutubeImporter_getCollectionID_ExistingCollection(t *testing.T) {
	ctx := context.Background()
	collection := &data.Collection{Name: "Test Playlist"}
	existingID := "existing_123"

	mockSearcher := mocks.NewITargetSearch(t)
	mockCollection := mocks.NewITargetCollection(t)
	mockTargetWriter := mocks.NewITargetWriter(t)
	mockMigrationState := mocks.NewIMigrationStateRepository(t)

	mockCollection.On("CheckIfCollectionExists", ctx, collection.Name).Return(existingID, nil)

	importer := newImporterTestNewImporter(mockSearcher, mockCollection, mockTargetWriter, mockMigrationState)

	collectionID, err := importer.getCollectionID(ctx, collection)

	assert.NoError(t, err)
	assert.Equal(t, existingID, collectionID)
}

func TestYoutubeImporter_getCollectionID_NewCollection(t *testing.T) {
	ctx := context.Background()
	collection := &data.Collection{Name: "New Playlist"}
	newID := "new_456"

	mockSearcher := mocks.NewITargetSearch(t)
	mockCollection := mocks.NewITargetCollection(t)
	mockTargetWriter := mocks.NewITargetWriter(t)
	mockMigrationState := mocks.NewIMigrationStateRepository(t)

	mockCollection.On("CheckIfCollectionExists", ctx, collection.Name).Return("", nil)
	mockCollection.On("CreateCollection", ctx, collection.Name).Return(newID, nil)

	importer := newImporterTestNewImporter(mockSearcher, mockCollection, mockTargetWriter, mockMigrationState)

	collectionID, err := importer.getCollectionID(ctx, collection)

	assert.NoError(t, err)
	assert.Equal(t, newID, collectionID)
}

func TestYoutubeImporter_getCollectionID_CheckExistError(t *testing.T) {
	ctx := context.Background()
	collection := &data.Collection{Name: "Test Playlist"}
	expectedError := assert.AnError

	mockSearcher := mocks.NewITargetSearch(t)
	mockCollection := mocks.NewITargetCollection(t)
	mockTargetWriter := mocks.NewITargetWriter(t)
	mockMigrationState := mocks.NewIMigrationStateRepository(t)

	mockCollection.On("CheckIfCollectionExists", ctx, collection.Name).Return("", expectedError)

	importer := newImporterTestNewImporter(mockSearcher, mockCollection, mockTargetWriter, mockMigrationState)

	collectionID, err := importer.getCollectionID(ctx, collection)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, collectionID)
}

func TestYoutubeImporter_getCollectionID_CreateCollectionError(t *testing.T) {
	ctx := context.Background()
	collection := &data.Collection{Name: "New Playlist"}
	expectedError := assert.AnError

	mockSearcher := mocks.NewITargetSearch(t)
	mockCollection := mocks.NewITargetCollection(t)
	mockTargetWriter := mocks.NewITargetWriter(t)
	mockMigrationState := mocks.NewIMigrationStateRepository(t)

	mockCollection.On("CheckIfCollectionExists", ctx, collection.Name).Return("", nil)
	mockCollection.On("CreateCollection", ctx, collection.Name).Return("", expectedError)

	importer := newImporterTestNewImporter(mockSearcher, mockCollection, mockTargetWriter, mockMigrationState)

	collectionID, err := importer.getCollectionID(ctx, collection)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, collectionID)
}

func TestYoutubeImporter_retrievePendingItems_StateExists(t *testing.T) {
	expectedItems := map[string]string{
		"song1": "youtube_id_1",
		"song2": "youtube_id_2",
	}

	mockSearcher := mocks.NewITargetSearch(t)
	mockCollection := mocks.NewITargetCollection(t)
	mockTargetWriter := mocks.NewITargetWriter(t)
	mockMigrationState := mocks.NewIMigrationStateRepository(t)

	mockMigrationState.On("Read").Return(true, nil)
	mockMigrationState.On("GetPendingItems").Return(expectedItems)

	importer := newImporterTestNewImporter(mockSearcher, mockCollection, mockTargetWriter, mockMigrationState)

	items, err := importer.retrievePendingItems()

	assert.NoError(t, err)
	assert.Equal(t, expectedItems, items)
}

func TestYoutubeImporter_retrievePendingItems_StateDoesNotExist(t *testing.T) {
	mockSearcher := mocks.NewITargetSearch(t)
	mockCollection := mocks.NewITargetCollection(t)
	mockTargetWriter := mocks.NewITargetWriter(t)
	mockMigrationState := mocks.NewIMigrationStateRepository(t)

	mockMigrationState.On("Read").Return(false, nil)

	importer := newImporterTestNewImporter(mockSearcher, mockCollection, mockTargetWriter, mockMigrationState)

	items, err := importer.retrievePendingItems()

	assert.NoError(t, err)
	assert.NotNil(t, items)
	assert.Len(t, items, 0)
}

func TestYoutubeImporter_retrievePendingItems_ReadError(t *testing.T) {
	expectedError := assert.AnError

	mockSearcher := mocks.NewITargetSearch(t)
	mockCollection := mocks.NewITargetCollection(t)
	mockTargetWriter := mocks.NewITargetWriter(t)
	mockMigrationState := mocks.NewIMigrationStateRepository(t)

	mockMigrationState.On("Read").Return(false, expectedError)

	importer := newImporterTestNewImporter(mockSearcher, mockCollection, mockTargetWriter, mockMigrationState)

	items, err := importer.retrievePendingItems()

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, items)
}

func TestYoutubeImporter_getNewItems_Success(t *testing.T) {
	ctx := context.Background()
	collection := &data.Collection{
		Name: "Test Playlist",
		Musics: []*data.Music{
			{Title: "Song 1", Artist: "Artist 1", Album: "Album 1"},
			{Title: "Song 2", Artist: "Artist 2", Album: "Album 2"},
		},
	}
	itemIDs := map[string]string{}

	mockSearcher := mocks.NewITargetSearch(t)
	mockCollection := mocks.NewITargetCollection(t)
	mockTargetWriter := mocks.NewITargetWriter(t)
	mockMigrationState := mocks.NewIMigrationStateRepository(t)

	mockSearcher.On("SearchItem", ctx, collection.Musics[0]).Return("youtube_id_1", nil)
	mockSearcher.On("SearchItem", ctx, collection.Musics[1]).Return("youtube_id_2", nil)
	mockMigrationState.On("AddItem", collection.Musics[0], "youtube_id_1").Return()
	mockMigrationState.On("AddItem", collection.Musics[1], "youtube_id_2").Return()
	mockMigrationState.On("Save").Return(nil)

	importer := newImporterTestNewImporter(mockSearcher, mockCollection, mockTargetWriter, mockMigrationState)

	err := importer.getNewItems(ctx, collection, itemIDs)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(itemIDs))
	assert.Contains(t, itemIDs, entities.ID(collection.Musics[0]))
	assert.Contains(t, itemIDs, entities.ID(collection.Musics[1]))
	assert.Equal(t, 2, importer.searchedItems)
}

func TestYoutubeImporter_getNewItems_WithExistingItems(t *testing.T) {
	ctx := context.Background()
	collection := &data.Collection{
		Name: "Test Playlist",
		Musics: []*data.Music{
			{Title: "Song 1", Artist: "Artist 1", Album: "Album 1"},
			{Title: "Song 2", Artist: "Artist 2", Album: "Album 2"},
		},
	}

	existingItemID := entities.ID(collection.Musics[0])
	itemIDs := map[string]string{
		existingItemID: "existing_youtube_id",
	}

	mockSearcher := mocks.NewITargetSearch(t)
	mockCollection := mocks.NewITargetCollection(t)
	mockTargetWriter := mocks.NewITargetWriter(t)
	mockMigrationState := mocks.NewIMigrationStateRepository(t)

	mockSearcher.On("SearchItem", ctx, collection.Musics[1]).Return("youtube_id_2", nil)
	mockMigrationState.On("AddItem", collection.Musics[1], "youtube_id_2").Return()
	mockMigrationState.On("Save").Return(nil)

	importer := newImporterTestNewImporter(mockSearcher, mockCollection, mockTargetWriter, mockMigrationState)

	err := importer.getNewItems(ctx, collection, itemIDs)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(itemIDs))
	assert.Contains(t, itemIDs, existingItemID)
	assert.Contains(t, itemIDs, entities.ID(collection.Musics[1]))
}

func TestYoutubeImporter_getNewItems_APILimitReached(t *testing.T) {
	ctx := context.Background()
	collection := &data.Collection{
		Name: "Test Playlist",
		Musics: []*data.Music{
			{Title: "Song 1", Artist: "Artist 1", Album: "Album 1"},
			{Title: "Song 2", Artist: "Artist 2", Album: "Album 2"},
		},
	}
	itemIDs := map[string]string{}

	mockSearcher := mocks.NewITargetSearch(t)
	mockCollection := mocks.NewITargetCollection(t)
	mockTargetWriter := mocks.NewITargetWriter(t)
	mockMigrationState := mocks.NewIMigrationStateRepository(t)

	mockSearcher.On("SearchItem", ctx, collection.Musics[0]).Return("youtube_id_1", nil)
	mockMigrationState.On("AddItem", collection.Musics[0], "youtube_id_1").Return()
	mockMigrationState.On("Save").Return(nil)

	importer := newImporterTestNewImporter(mockSearcher, mockCollection, mockTargetWriter, mockMigrationState)
	importer.apiLimit = 1

	err := importer.getNewItems(ctx, collection, itemIDs)

	assert.NoError(t, err)
	assert.Equal(t, 1, importer.searchedItems)
}

func TestYoutubeImporter_getNewItems_SearchError(t *testing.T) {
	ctx := context.Background()
	collection := &data.Collection{
		Name: "Test Playlist",
		Musics: []*data.Music{
			{Title: "Song 1", Artist: "Artist 1", Album: "Album 1"},
		},
	}
	itemIDs := map[string]string{}
	expectedError := assert.AnError

	mockSearcher := mocks.NewITargetSearch(t)
	mockCollection := mocks.NewITargetCollection(t)
	mockTargetWriter := mocks.NewITargetWriter(t)
	mockMigrationState := mocks.NewIMigrationStateRepository(t)

	mockSearcher.On("SearchItem", ctx, collection.Musics[0]).Return("", expectedError)
	mockMigrationState.On("Save").Return(nil)

	importer := newImporterTestNewImporter(mockSearcher, mockCollection, mockTargetWriter, mockMigrationState)

	err := importer.getNewItems(ctx, collection, itemIDs)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 0, len(itemIDs))
}

func TestYoutubeImporter_insertAll_Success(t *testing.T) {
	ctx := context.Background()
	collectionID := "collection_123"
	itemIDs := map[string]string{
		"item1": "youtube_id_1",
		"item2": "youtube_id_2",
	}

	mockSearcher := mocks.NewITargetSearch(t)
	mockCollection := mocks.NewITargetCollection(t)
	mockTargetWriter := mocks.NewITargetWriter(t)
	mockMigrationState := mocks.NewIMigrationStateRepository(t)

	mockTargetWriter.On("AddItemToPlaylist", ctx, collectionID, "youtube_id_1").Return(nil)
	mockTargetWriter.On("AddItemToPlaylist", ctx, collectionID, "youtube_id_2").Return(nil)
	mockMigrationState.On("UpdateItemToMigrated", "item1").Return()
	mockMigrationState.On("UpdateItemToMigrated", "item2").Return()
	mockMigrationState.On("Save").Return(nil)

	importer := newImporterTestNewImporter(mockSearcher, mockCollection, mockTargetWriter, mockMigrationState)

	err := importer.insertAll(ctx, collectionID, itemIDs)

	assert.NoError(t, err)
}

func TestYoutubeImporter_insertAll_EmptyCollectionID(t *testing.T) {
	ctx := context.Background()
	itemIDs := map[string]string{
		"item1": "youtube_id_1",
	}

	mockSearcher := mocks.NewITargetSearch(t)
	mockCollection := mocks.NewITargetCollection(t)
	mockTargetWriter := mocks.NewITargetWriter(t)
	mockMigrationState := mocks.NewIMigrationStateRepository(t)

	importer := newImporterTestNewImporter(mockSearcher, mockCollection, mockTargetWriter, mockMigrationState)

	err := importer.insertAll(ctx, "", itemIDs)

	assert.NoError(t, err)
}

func TestYoutubeImporter_insertAll_EmptyItems(t *testing.T) {
	ctx := context.Background()
	collectionID := "collection_123"
	itemIDs := map[string]string{}

	mockSearcher := mocks.NewITargetSearch(t)
	mockCollection := mocks.NewITargetCollection(t)
	mockTargetWriter := mocks.NewITargetWriter(t)
	mockMigrationState := mocks.NewIMigrationStateRepository(t)

	importer := newImporterTestNewImporter(mockSearcher, mockCollection, mockTargetWriter, mockMigrationState)

	err := importer.insertAll(ctx, collectionID, itemIDs)

	assert.NoError(t, err)
}

func TestYoutubeImporter_insertAll_AddItemError(t *testing.T) {
	ctx := context.Background()
	collectionID := "collection_123"
	itemIDs := map[string]string{
		"item1": "youtube_id_1",
	}
	expectedError := assert.AnError

	mockSearcher := mocks.NewITargetSearch(t)
	mockCollection := mocks.NewITargetCollection(t)
	mockTargetWriter := mocks.NewITargetWriter(t)
	mockMigrationState := mocks.NewIMigrationStateRepository(t)

	mockTargetWriter.On("AddItemToPlaylist", ctx, collectionID, "youtube_id_1").Return(expectedError)
	mockMigrationState.On("Save").Return(nil)

	importer := newImporterTestNewImporter(mockSearcher, mockCollection, mockTargetWriter, mockMigrationState)

	err := importer.insertAll(ctx, collectionID, itemIDs)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestYoutubeImporter_insertAll_APILimitReached(t *testing.T) {
	ctx := context.Background()
	collectionID := "collection_123"

	itemIDs := make(map[string]string)
	for i := 0; i < API_LIMIT+10; i++ {
		itemIDs[string(rune('a'+i%26))+string(rune('0'+i%10))] = "youtube_id_" + string(rune('0'+i%10))
	}

	mockSearcher := mocks.NewITargetSearch(t)
	mockCollection := mocks.NewITargetCollection(t)
	mockTargetWriter := mocks.NewITargetWriter(t)
	mockMigrationState := mocks.NewIMigrationStateRepository(t)

	for i := 0; i < API_LIMIT; i++ {
		mockTargetWriter.On("AddItemToPlaylist", ctx, collectionID, mock.AnythingOfType("string")).Return(nil).Once()
		mockMigrationState.On("UpdateItemToMigrated", mock.AnythingOfType("string")).Return().Once()
	}
	mockMigrationState.On("Save").Return(nil)

	importer := newImporterTestNewImporter(mockSearcher, mockCollection, mockTargetWriter, mockMigrationState)

	err := importer.insertAll(ctx, collectionID, itemIDs)

	assert.NoError(t, err)
}

func BenchmarkYoutubeImporter_Import(b *testing.B) {
	ctx := context.Background()
	collection := &data.Collection{
		Name: "Benchmark Playlist",
		Musics: []*data.Music{
			{Title: "Benchmark Song", Artist: "Benchmark Artist", Album: "Benchmark Album"},
		},
	}

	pendingItems := map[string]string{}

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		mockSearcher := &mocks.ITargetSearch{}
		mockCollection := &mocks.ITargetCollection{}
		mockTargetWriter := &mocks.ITargetWriter{}
		mockMigrationState := &mocks.IMigrationStateRepository{}

		mockCollection.On("CheckIfCollectionExists", mock.Anything, mock.Anything).Return("collection_123", nil)
		mockMigrationState.On("Read", mock.Anything).Return(true, nil)
		mockMigrationState.On("GetPendingItems", mock.Anything).Return(pendingItems)
		mockSearcher.On("SearchItem", mock.Anything, mock.Anything).Return("youtube_id", nil)
		mockMigrationState.On("AddItem", mock.Anything, mock.Anything).Return()
		mockTargetWriter.On("AddItemToPlaylist", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		mockMigrationState.On("UpdateItemToMigrated", mock.Anything).Return()
		mockMigrationState.On("Save", mock.Anything).Return(nil)

		importer := NewImporter(mockSearcher, mockCollection, mockTargetWriter, mockMigrationState)
		b.StartTimer()

		_, err := importer.Import(ctx, collection)
		if err != nil {
			b.Fatal(err)
		}
	}
}
