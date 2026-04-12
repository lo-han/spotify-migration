package main

import (
	"context"
	"log"
	"os"
	"spotify_migration/adapters"
	"spotify_migration/entities"
	"spotify_migration/entities/data"
	"spotify_migration/usecases/album"
	"spotify_migration/usecases/playlist"

	"github.com/lo-han/spotify/v2"
)

func main() {
	if len(os.Args) < 3 {
		log.Println("Please provide a resource kind and name")
		return
	}
	resourceKind := os.Args[1]
	resourceName := os.Args[2]

	var spotifyExtractor entities.IExtractorUsecase
	var youtubeImporter entities.IImporterUsecase

	ctx := context.Background()
	auth, token, youtubeService := Authorize(ctx)

	switch resourceKind {
	case data.PlaylistKind:
		spotifyExtractor = playlist.NewExtractor(adapters.NewSpotifyGetter(spotify.New(auth.Client(ctx, token))))

		youtubeImporter = playlist.NewImporter(
			adapters.NewYoutubeSearch(youtubeService),
			adapters.NewYoutubeCollection(youtubeService),
			adapters.NewYoutubeCollectionWriter(youtubeService),
			adapters.NewMigrationState(resourceName),
		)

	case data.AlbumKind:
		spotifyExtractor = album.NewExtractor(adapters.NewSpotifyAlbumSearch(spotify.New(auth.Client(ctx, token))))
		youtubeImporter = album.NewImporter(adapters.NewYoutubeAlbumSave(youtubeService))

	default:
		log.Printf("Unsupported resource kind: %s\n", resourceKind)
		return
	}

	migration := entities.NewMigration(spotifyExtractor, youtubeImporter)

	if ok, err := migration.Migrate(ctx, resourceName); err != nil {
		log.Println("Error migrating resource:", err)
	} else if ok {
		log.Printf("%s migrated successfully: %s\n", resourceKind, resourceName)
	}
}
