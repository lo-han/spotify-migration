package usecases

import (
	"context"
	"errors"
	"spotify_migration/entities/data"
	"spotify_migration/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPlaylistExtractor(t *testing.T) {
	// Arrange
	mockSourceGetter := mocks.NewISourceGetter(t)

	// Act
	extractor := NewPlaylistExtractor(mockSourceGetter)

	// Assert
	assert.NotNil(t, extractor)
	assert.IsType(t, &spotifyPlaylistExtractor{}, extractor)
}

func TestSpotifyPlaylistExtractor_Extract_Success(t *testing.T) {
	ctx := context.Background()
	resourceName := "My Test Playlist"
	playlistID := "spotify:playlist:1234567890"

	expectedCollection := &data.Collection{
		Name: "My Test Playlist",
		Musics: []*data.Music{
			{
				Title:  "Song 1",
				Artist: "Artist 1",
				Album:  "Album 1",
			},
			{
				Title:  "Song 2",
				Artist: "Artist 2",
				Album:  "Album 2",
			},
		},
	}

	mockSourceGetter := mocks.NewISourceGetter(t)
	mockSourceGetter.On("GetPlaylistID", ctx, resourceName).Return(playlistID, nil)
	mockSourceGetter.On("GetPlaylistItems", ctx, resourceName, playlistID).Return(expectedCollection, nil)

	extractor := NewPlaylistExtractor(mockSourceGetter)

	result, err := extractor.Extract(ctx, resourceName)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCollection, result)
	assert.Equal(t, "My Test Playlist", result.Name)
	assert.Len(t, result.Musics, 2)

	assert.Equal(t, "Song 1", result.Musics[0].Title)
	assert.Equal(t, "Artist 1", result.Musics[0].Artist)
	assert.Equal(t, "Album 1", result.Musics[0].Album)

	assert.Equal(t, "Song 2", result.Musics[1].Title)
	assert.Equal(t, "Artist 2", result.Musics[1].Artist)
	assert.Equal(t, "Album 2", result.Musics[1].Album)

	mockSourceGetter.AssertExpectations(t)
}

func TestSpotifyPlaylistExtractor_Extract_EmptyPlaylist(t *testing.T) {
	ctx := context.Background()
	resourceName := "Empty Playlist"
	playlistID := "spotify:playlist:empty123"

	expectedCollection := &data.Collection{
		Name:   "Empty Playlist",
		Musics: []*data.Music{},
	}

	mockSourceGetter := mocks.NewISourceGetter(t)
	mockSourceGetter.On("GetPlaylistID", ctx, resourceName).Return(playlistID, nil)
	mockSourceGetter.On("GetPlaylistItems", ctx, resourceName, playlistID).Return(expectedCollection, nil)

	extractor := NewPlaylistExtractor(mockSourceGetter)

	result, err := extractor.Extract(ctx, resourceName)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCollection, result)
	assert.Equal(t, "Empty Playlist", result.Name)
	assert.Len(t, result.Musics, 0)

	mockSourceGetter.AssertExpectations(t)
}

func TestSpotifyPlaylistExtractor_Extract_GetPlaylistIDError(t *testing.T) {
	ctx := context.Background()
	resourceName := "Nonexistent Playlist"
	expectedError := errors.New("playlist not found")

	mockSourceGetter := mocks.NewISourceGetter(t)
	mockSourceGetter.On("GetPlaylistID", ctx, resourceName).Return("", expectedError)

	extractor := NewPlaylistExtractor(mockSourceGetter)

	result, err := extractor.Extract(ctx, resourceName)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)

	mockSourceGetter.AssertExpectations(t)
}

func TestSpotifyPlaylistExtractor_Extract_GetPlaylistItemsError(t *testing.T) {
	ctx := context.Background()
	resourceName := "Problematic Playlist"
	playlistID := "spotify:playlist:problem123"
	expectedError := errors.New("failed to fetch playlist items")

	mockSourceGetter := mocks.NewISourceGetter(t)
	mockSourceGetter.On("GetPlaylistID", ctx, resourceName).Return(playlistID, nil)
	mockSourceGetter.On("GetPlaylistItems", ctx, resourceName, playlistID).Return(nil, expectedError)

	extractor := NewPlaylistExtractor(mockSourceGetter)

	result, err := extractor.Extract(ctx, resourceName)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)

	mockSourceGetter.AssertExpectations(t)
}
