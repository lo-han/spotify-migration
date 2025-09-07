package adapters

import "spotify_migration/domain"

func NewSpotifyAlbumExtractor() *SpotifyAlbumExtractor {
	return &SpotifyAlbumExtractor{}
}

type SpotifyAlbumExtractor struct {
}

func (s *SpotifyAlbumExtractor) Extract(resourceName string) (*domain.Collection, error) {
	albumID, err := s.getAlbumID(resourceName)
	if err != nil {
		return nil, err
	}

	albumItems, err := s.getAlbumItems(albumID)
	if err != nil {
		return nil, err
	}

	return albumItems, nil
}

func (s *SpotifyAlbumExtractor) getAlbumID(resourceName string) (string, error) {
	// Simulate fetching album ID by resource name
	if resourceName == "" {
		return "", nil
	}
	return "album_id", nil
}

func (s *SpotifyAlbumExtractor) getAlbumItems(id string) (*domain.Collection, error) {
	// Simulate fetching album items by album ID
	if id == "" {
		return nil, nil
	}
	return &domain.Collection{
		Name: "My Album",
		Musics: []*domain.Music{
			{Title: "Song 1", Artist: "Artist A", Album: "My Album"},
			{Title: "Song 2", Artist: "Artist A", Album: "My Album"},
		},
	}, nil
}
