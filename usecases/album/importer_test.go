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
	mockSaveAlbums := mocks.NewISaveAlbums(t)

	// Act
	importerUsecase := NewImporter(mockSaveAlbums)

	// Assert
	assert.NotNil(t, importerUsecase)
	assert.IsType(t, &playlistImporter{}, importerUsecase)
}

func TestAlbumImporter_Import_Success(t *testing.T) {
	ctx := context.Background()
	albums := []*data.Album{
		{
			Title:   "Abbey Road",
			Artists: []string{"The Beatles"},
		},
		{
			Title:   "Dark Side of the Moon",
			Artists: []string{"Pink Floyd"},
		},
	}

	mockSaveAlbums := mocks.NewISaveAlbums(t)
	mockSaveAlbums.On("SaveAlbums", ctx, albums).Return(nil)

	albumImporter := NewImporter(mockSaveAlbums)

	result, err := albumImporter.Import(ctx, albums)

	assert.NoError(t, err)
	assert.True(t, result)

	mockSaveAlbums.AssertExpectations(t)
}

func TestAlbumImporter_Import_EmptyAlbums(t *testing.T) {
	ctx := context.Background()
	albums := []*data.Album{}

	mockSaveAlbums := mocks.NewISaveAlbums(t)

	albumImporter := NewImporter(mockSaveAlbums)

	result, err := albumImporter.Import(ctx, albums)

	assert.NoError(t, err)
	assert.False(t, result)

	mockSaveAlbums.AssertExpectations(t)
}

func TestAlbumImporter_Import_NilCollection(t *testing.T) {
	mockSaveAlbums := mocks.NewISaveAlbums(t)

	albumImporter := NewImporter(mockSaveAlbums)

	result, err := albumImporter.Import(context.Background(), nil)

	assert.NoError(t, err)
	assert.False(t, result)

	mockSaveAlbums.AssertExpectations(t)
}

func TestAlbumImporter_Import_InvalidCollectionType(t *testing.T) {
	ctx := context.Background()
	invalidCollection := "not an album slice"

	mockSaveAlbums := mocks.NewISaveAlbums(t)

	albumImporter := NewImporter(mockSaveAlbums)

	result, err := albumImporter.Import(ctx, invalidCollection)

	assert.Error(t, err)
	assert.False(t, result)
	assert.Equal(t, "invalid album data", err.Error())

	mockSaveAlbums.AssertExpectations(t)
}

func TestAlbumImporter_Import_SaveAlbumsError(t *testing.T) {
	ctx := context.Background()
	albums := []*data.Album{
		{
			Title:   "Thriller",
			Artists: []string{"Michael Jackson"},
		},
	}
	expectedError := errors.New("failed to save albums")

	mockSaveAlbums := mocks.NewISaveAlbums(t)
	mockSaveAlbums.On("SaveAlbums", ctx, albums).Return(expectedError)

	albumImporter := NewImporter(mockSaveAlbums)

	result, err := albumImporter.Import(ctx, albums)

	assert.Error(t, err)
	assert.False(t, result)
	assert.Equal(t, expectedError, err)

	mockSaveAlbums.AssertExpectations(t)
}

func TestAlbumImporter_Import_SingleAlbum(t *testing.T) {
	ctx := context.Background()
	albums := []*data.Album{
		{
			Title:   "Nevermind",
			Artists: []string{"Nirvana"},
		},
	}

	mockSaveAlbums := mocks.NewISaveAlbums(t)
	mockSaveAlbums.On("SaveAlbums", ctx, albums).Return(nil)

	albumImporter := NewImporter(mockSaveAlbums)

	result, err := albumImporter.Import(ctx, albums)

	assert.NoError(t, err)
	assert.True(t, result)

	mockSaveAlbums.AssertExpectations(t)
}

func TestAlbumImporter_Import_LargeAlbumCollection(t *testing.T) {
	ctx := context.Background()

	// Create a larger collection of albums
	albums := make([]*data.Album, 100)
	for i := 0; i < 100; i++ {
		albums[i] = &data.Album{
			Title:   "Album " + string(rune(i)),
			Artists: []string{"Artists " + string(rune(i))},
		}
	}

	mockSaveAlbums := mocks.NewISaveAlbums(t)
	mockSaveAlbums.On("SaveAlbums", ctx, albums).Return(nil)

	albumImporter := NewImporter(mockSaveAlbums)

	result, err := albumImporter.Import(ctx, albums)

	assert.NoError(t, err)
	assert.True(t, result)

	mockSaveAlbums.AssertExpectations(t)
}
