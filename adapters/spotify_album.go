package adapters

import "spotify_migration/domain"

func NewSpotifyPlaylistExtractor() *SpotifyPlaylistExtractor {
	return &SpotifyPlaylistExtractor{}
}

type SpotifyPlaylistExtractor struct {
}

func (s *SpotifyPlaylistExtractor) Extract(resourceName string) (*domain.Collection, error) {
	playlistID, err := s.getPlaylistID(resourceName)
	if err != nil {
		return nil, err
	}

	playlistItems, err := s.getPlaylistItems(playlistID)
	if err != nil {
		return nil, err
	}

	return playlistItems, nil
}

func (s *SpotifyPlaylistExtractor) getPlaylistID(resourceName string) (string, error) {
	// Simulate fetching playlist ID by resource name
	if resourceName == "" {
		return "", nil
	}
	return "playlist_id", nil
}

func (s *SpotifyPlaylistExtractor) getPlaylistItems(id string) (*domain.Collection, error) {
	// Simulate fetching playlist items by playlist ID
	if id == "" {
		return nil, nil
	}
	return &domain.Collection{
		Name: "My Playlist",
		Musics: []*domain.Music{
			{Title: "Song 1"},
			{Title: "Song 2"},
		},
	}, nil
}
