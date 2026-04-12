package album

import (
	"context"
	"errors"
	"spotify_migration/entities/data"
	"spotify_migration/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewImporter(t *testing.T) {
	// Arrange
	mockSaveAlbuns := mocks.NewISaveAlbuns(t)

	// Act
	importerUsecase := NewImporter(mockSaveAlbuns)

	// Assert
	assert.NotNil(t, importerUsecase)
	assert.IsType(t, &playlistImporter{}, importerUsecase)
}

func TestAlbumImporter_Import_Success(t *testing.T) {
	ctx := context.Background()
	albums := []*data.Album{
		{
			Title:  "Abbey Road",
			Artist: "The Beatles",
		},
		{
			Title:  "Dark Side of the Moon",
			Artist: "Pink Floyd",
		},
	}

	mockSaveAlbuns := mocks.NewISaveAlbuns(t)
	mockSaveAlbuns.On("SaveAlbuns", ctx, albums).Return(nil)

	albumImporter := NewImporter(mockSaveAlbuns)

	result, err := albumImporter.Import(ctx, albums)

	assert.NoError(t, err)
	assert.True(t, result)

	mockSaveAlbuns.AssertExpectations(t)
}

func TestAlbumImporter_Import_EmptyAlbums(t *testing.T) {
	ctx := context.Background()
	albums := []*data.Album{}

	mockSaveAlbuns := mocks.NewISaveAlbuns(t)

	albumImporter := NewImporter(mockSaveAlbuns)

	result, err := albumImporter.Import(ctx, albums)

	assert.NoError(t, err)
	assert.False(t, result)

	mockSaveAlbuns.AssertExpectations(t)
}

func TestAlbumImporter_Import_NilCollection(t *testing.T) {
	mockSaveAlbuns := mocks.NewISaveAlbuns(t)

	albumImporter := NewImporter(mockSaveAlbuns)

	result, err := albumImporter.Import(context.Background(), nil)

	assert.NoError(t, err)
	assert.False(t, result)

	mockSaveAlbuns.AssertExpectations(t)
}

func TestAlbumImporter_Import_InvalidCollectionType(t *testing.T) {
	ctx := context.Background()
	invalidCollection := "not an album slice"

	mockSaveAlbuns := mocks.NewISaveAlbuns(t)

	albumImporter := NewImporter(mockSaveAlbuns)

	result, err := albumImporter.Import(ctx, invalidCollection)

	assert.Error(t, err)
	assert.False(t, result)
	assert.Equal(t, "invalid album data", err.Error())

	mockSaveAlbuns.AssertExpectations(t)
}

func TestAlbumImporter_Import_SaveAlbunsError(t *testing.T) {
	ctx := context.Background()
	albums := []*data.Album{
		{
			Title:  "Thriller",
			Artist: "Michael Jackson",
		},
	}
	expectedError := errors.New("failed to save albums")

	mockSaveAlbuns := mocks.NewISaveAlbuns(t)
	mockSaveAlbuns.On("SaveAlbuns", ctx, albums).Return(expectedError)

	albumImporter := NewImporter(mockSaveAlbuns)

	result, err := albumImporter.Import(ctx, albums)

	assert.Error(t, err)
	assert.False(t, result)
	assert.Equal(t, expectedError, err)

	mockSaveAlbuns.AssertExpectations(t)
}

func TestAlbumImporter_Import_SingleAlbum(t *testing.T) {
	ctx := context.Background()
	albums := []*data.Album{
		{
			Title:  "Nevermind",
			Artist: "Nirvana",
		},
	}

	mockSaveAlbuns := mocks.NewISaveAlbuns(t)
	mockSaveAlbuns.On("SaveAlbuns", ctx, albums).Return(nil)

	albumImporter := NewImporter(mockSaveAlbuns)

	result, err := albumImporter.Import(ctx, albums)

	assert.NoError(t, err)
	assert.True(t, result)

	mockSaveAlbuns.AssertExpectations(t)
}

func TestAlbumImporter_Import_LargeAlbumCollection(t *testing.T) {
	ctx := context.Background()

	// Create a larger collection of albums
	albums := make([]*data.Album, 100)
	for i := 0; i < 100; i++ {
		albums[i] = &data.Album{
			Title:  "Album " + string(rune(i)),
			Artist: "Artist " + string(rune(i)),
		}
	}

	mockSaveAlbuns := mocks.NewISaveAlbuns(t)
	mockSaveAlbuns.On("SaveAlbuns", ctx, albums).Return(nil)

	albumImporter := NewImporter(mockSaveAlbuns)

	result, err := albumImporter.Import(ctx, albums)

	assert.NoError(t, err)
	assert.True(t, result)

	mockSaveAlbuns.AssertExpectations(t)
}
