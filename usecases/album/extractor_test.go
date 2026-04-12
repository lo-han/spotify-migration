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
	mockGetAlbums := mocks.NewIGetAlbums(t)

	// Act
	extractorUsecase := NewExtractor(mockGetAlbums)

	// Assert
	assert.NotNil(t, extractorUsecase)
	assert.IsType(t, &albumExtractor{}, extractorUsecase)
}

func TestAlbumExtractor_Extract_Success(t *testing.T) {
	ctx := context.Background()
	resourceName := "My Saved Albums"

	expectedAlbums := []*data.Album{
		{
			Title:   "Abbey Road",
			Artists: []string{"The Beatles"},
		},
		{
			Title:   "Dark Side of the Moon",
			Artists: []string{"Pink Floyd"},
		},
	}

	mockGetAlbums := mocks.NewIGetAlbums(t)
	mockGetAlbums.On("GetAlbums", ctx).Return(expectedAlbums, nil)

	albumExtractor := NewExtractor(mockGetAlbums)

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
	assert.Equal(t, "The Beatles", albums[0].Artists)

	assert.Equal(t, "Dark Side of the Moon", albums[1].Title)
	assert.Equal(t, "Pink Floyd", albums[1].Artists)

	mockGetAlbums.AssertExpectations(t)
}

func TestAlbumExtractor_Extract_EmptyAlbums(t *testing.T) {
	ctx := context.Background()
	resourceName := "Empty Albums"

	expectedAlbums := []*data.Album{}

	mockGetAlbums := mocks.NewIGetAlbums(t)
	mockGetAlbums.On("GetAlbums", ctx).Return(expectedAlbums, nil)

	albumExtractor := NewExtractor(mockGetAlbums)

	result, err := albumExtractor.Extract(ctx, resourceName)

	albums, ok := result.([]*data.Album)
	if !ok {
		panic("invalid album data")
	}

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedAlbums, albums)
	assert.Len(t, albums, 0)

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

	expectedAlbums := []*data.Album{
		{
			Title:   "Thriller",
			Artists: []string{"Michael Jackson"},
		},
	}

	mockGetAlbums := mocks.NewIGetAlbums(t)
	mockGetAlbums.On("GetAlbums", ctx).Return(expectedAlbums, nil)

	albumExtractor := NewExtractor(mockGetAlbums)

	result, err := albumExtractor.Extract(ctx, resourceName)

	albums, ok := result.([]*data.Album)
	if !ok {
		panic("invalid album data")
	}

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, albums, 1)
	assert.Equal(t, "Thriller", albums[0].Title)
	assert.Equal(t, "Michael Jackson", albums[0].Artists)

	mockGetAlbums.AssertExpectations(t)
}
