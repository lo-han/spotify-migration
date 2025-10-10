package main

import (
	"fmt"
	"os"
	"spotify_migration/adapters/extractor"
	"spotify_migration/adapters/importer"
	"spotify_migration/domain"
	"spotify_migration/ports"
	"spotify_migration/usecases"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please provide a resource kind and name")
		return
	}
	resourceKind := os.Args[1]
	resourceName := os.Args[2]

	var spotify ports.IExtractor

	youtube := importer.NewYoutubeImporter()

	switch resourceKind {
	case domain.PlaylistKind:
		spotify = extractor.NewSpotifyPlaylistExtractor()

	case domain.AlbumKind:
		spotify = extractor.NewSpotifyAlbumExtractor()
	}

	migration := usecases.NewMigration(spotify, youtube)

	if ok, err := migration.Migrate(resourceName); err != nil {
		fmt.Println("Error migrating resource:", err)
	} else if ok {
		fmt.Printf("%s migrated successfully: %s\n", resourceKind, resourceName)
	}
}
