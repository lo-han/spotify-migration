package album

import (
	"context"
	"errors"
	"spotify_migration/entities"
	"spotify_migration/entities/data"
	"spotify_migration/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewImporter(t *testing.T) {
	mockSaveAlbums := mocks.NewISaveAlbums(t)

	importerUsecase := NewImporter(mockSaveAlbums)

	assert.NotNil(t, importerUsecase)
	assert.IsType(t, &playlistImporter{}, importerUsecase)
}

func TestAlbumImporter_Import_Success(t *testing.T) {
	ctx := context.Background()
	albums := &data.Collection{
		Name:   "Abbey Road",
		Artist: entities.PtrStr("The Beatles"),
	}
	albums.Append(&data.Collection{
		Name:   "Dark Side of the Moon",
		Artist: entities.PtrStr("Pink Floyd"),
	})

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

	mockSaveAlbums := mocks.NewISaveAlbums(t)

	albumImporter := NewImporter(mockSaveAlbums)

	result, err := albumImporter.Import(ctx, nil)

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

func TestAlbumImporter_Import_SaveAlbumsError(t *testing.T) {
	ctx := context.Background()
	albums := &data.Collection{
		Name:   "Thriller",
		Artist: entities.PtrStr("Michael Jackson"),
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
	albums := &data.Collection{
		Name:   "Nevermind",
		Artist: entities.PtrStr("Nirvana"),
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
	albums := &data.Collection{}
	for i := 0; i < 100; i++ {
		albums.Append(&data.Collection{
			Name:   "Collection " + string(rune(i)),
			Artist: entities.PtrStr("Artist " + string(rune(i))),
		})
	}

	mockSaveAlbums := mocks.NewISaveAlbums(t)
	mockSaveAlbums.On("SaveAlbums", ctx, albums).Return(nil)

	albumImporter := NewImporter(mockSaveAlbums)

	result, err := albumImporter.Import(ctx, albums)

	assert.NoError(t, err)
	assert.True(t, result)

	mockSaveAlbums.AssertExpectations(t)
}
