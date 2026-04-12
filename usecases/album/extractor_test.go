package album

import (
	"context"
	"errors"
	"spotify_migration/entities/data"
	"spotify_migration/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewExtractor(t *testing.T) {
	// Arrange
	mockGetAlbuns := mocks.NewIGetAlbuns(t)

	// Act
	extractorUsecase := NewExtractor(mockGetAlbuns)

	// Assert
	assert.NotNil(t, extractorUsecase)
	assert.IsType(t, &albumExtractor{}, extractorUsecase)
}

func TestAlbumExtractor_Extract_Success(t *testing.T) {
	ctx := context.Background()
	resourceName := "My Saved Albums"

	expectedAlbums := []*data.Album{
		{
			Title:  "Abbey Road",
			Artist: "The Beatles",
		},
		{
			Title:  "Dark Side of the Moon",
			Artist: "Pink Floyd",
		},
	}

	mockGetAlbuns := mocks.NewIGetAlbuns(t)
	mockGetAlbuns.On("GetAlbuns", ctx).Return(expectedAlbums, nil)

	albumExtractor := NewExtractor(mockGetAlbuns)

	result, err := albumExtractor.Extract(ctx, resourceName)

	albums, ok := result.([]*data.Album)
	if !ok {
		panic("invalid album data")
	}

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedAlbums, albums)
	assert.Len(t, albums, 2)

	assert.Equal(t, "Abbey Road", albums[0].Title)
	assert.Equal(t, "The Beatles", albums[0].Artist)

	assert.Equal(t, "Dark Side of the Moon", albums[1].Title)
	assert.Equal(t, "Pink Floyd", albums[1].Artist)

	mockGetAlbuns.AssertExpectations(t)
}

func TestAlbumExtractor_Extract_EmptyAlbums(t *testing.T) {
	ctx := context.Background()
	resourceName := "Empty Albums"

	expectedAlbums := []*data.Album{}

	mockGetAlbuns := mocks.NewIGetAlbuns(t)
	mockGetAlbuns.On("GetAlbuns", ctx).Return(expectedAlbums, nil)

	albumExtractor := NewExtractor(mockGetAlbuns)

	result, err := albumExtractor.Extract(ctx, resourceName)

	albums, ok := result.([]*data.Album)
	if !ok {
		panic("invalid album data")
	}

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedAlbums, albums)
	assert.Len(t, albums, 0)

	mockGetAlbuns.AssertExpectations(t)
}

func TestAlbumExtractor_Extract_GetAlbunsError(t *testing.T) {
	ctx := context.Background()
	resourceName := "My Albums"
	expectedError := errors.New("failed to fetch albums")

	mockGetAlbuns := mocks.NewIGetAlbuns(t)
	mockGetAlbuns.On("GetAlbuns", ctx).Return(nil, expectedError)

	albumExtractor := NewExtractor(mockGetAlbuns)

	result, err := albumExtractor.Extract(ctx, resourceName)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)

	mockGetAlbuns.AssertExpectations(t)
}

func TestAlbumExtractor_Extract_SingleAlbum(t *testing.T) {
	ctx := context.Background()
	resourceName := "My Albums"

	expectedAlbums := []*data.Album{
		{
			Title:  "Thriller",
			Artist: "Michael Jackson",
		},
	}

	mockGetAlbuns := mocks.NewIGetAlbuns(t)
	mockGetAlbuns.On("GetAlbuns", ctx).Return(expectedAlbums, nil)

	albumExtractor := NewExtractor(mockGetAlbuns)

	result, err := albumExtractor.Extract(ctx, resourceName)

	albums, ok := result.([]*data.Album)
	if !ok {
		panic("invalid album data")
	}

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, albums, 1)
	assert.Equal(t, "Thriller", albums[0].Title)
	assert.Equal(t, "Michael Jackson", albums[0].Artist)

	mockGetAlbuns.AssertExpectations(t)
}
