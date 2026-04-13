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

func TestNewExtractor(t *testing.T) {
	mockGetAlbums := mocks.NewIGetAlbums(t)

	extractorUsecase := NewExtractor(mockGetAlbums)

	assert.NotNil(t, extractorUsecase)
	assert.IsType(t, &albumExtractor{}, extractorUsecase)
}

func TestAlbumExtractor_Extract_Success(t *testing.T) {
	ctx := context.Background()
	resourceName := "My Saved Albums"

	expectedAlbums := &data.Collection{
		Name:   "Abbey Road",
		Artist: entities.PtrStr("The Beatles"),
	}
	expectedAlbums.Append(
		&data.Collection{
			Name:   "Dark Side of the Moon",
			Artist: entities.PtrStr("Pink Floyd"),
		},
	)

	mockGetAlbums := mocks.NewIGetAlbums(t)
	mockGetAlbums.On("GetAlbums", ctx).Return(expectedAlbums, nil)

	albumExtractor := NewExtractor(mockGetAlbums)

	result, err := albumExtractor.Extract(ctx, resourceName)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, "Abbey Road", result.Name)
	assert.Equal(t, "The Beatles", *result.Artist)

	result = result.Next()
	assert.Equal(t, "Dark Side of the Moon", result.Name)
	assert.Equal(t, "Pink Floyd", *result.Artist)

	mockGetAlbums.AssertExpectations(t)
}

func TestAlbumExtractor_Extract_EmptyAlbums(t *testing.T) {
	ctx := context.Background()
	resourceName := "Empty Albums"

	mockGetAlbums := mocks.NewIGetAlbums(t)
	mockGetAlbums.On("GetAlbums", ctx).Return(nil, nil)

	albumExtractor := NewExtractor(mockGetAlbums)

	result, err := albumExtractor.Extract(ctx, resourceName)

	assert.NoError(t, err)
	assert.Nil(t, result)

	mockGetAlbums.AssertExpectations(t)
}

func TestAlbumExtractor_Extract_GetAlbumsError(t *testing.T) {
	ctx := context.Background()
	resourceName := "My Albums"
	expectedError := errors.New("failed to fetch albums")

	mockGetAlbums := mocks.NewIGetAlbums(t)
	mockGetAlbums.On("GetAlbums", ctx).Return(nil, expectedError)

	albumExtractor := NewExtractor(mockGetAlbums)

	result, err := albumExtractor.Extract(ctx, resourceName)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)

	mockGetAlbums.AssertExpectations(t)
}

func TestAlbumExtractor_Extract_SingleAlbum(t *testing.T) {
	ctx := context.Background()
	resourceName := "My Albums"

	expectedAlbums := &data.Collection{
		Name:   "Thriller",
		Artist: entities.PtrStr("Michael Jackson"),
	}

	mockGetAlbums := mocks.NewIGetAlbums(t)
	mockGetAlbums.On("GetAlbums", ctx).Return(expectedAlbums, nil)

	albumExtractor := NewExtractor(mockGetAlbums)

	result, err := albumExtractor.Extract(ctx, resourceName)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Thriller", result.Name)
	assert.Equal(t, "Michael Jackson", *result.Artist)

	mockGetAlbums.AssertExpectations(t)
}
