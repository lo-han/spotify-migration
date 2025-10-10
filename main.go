package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"spotify_migration/adapters/extractor"
	"spotify_migration/adapters/importer"
	"spotify_migration/domain"
	"spotify_migration/ports"
	"spotify_migration/usecases"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please provide a resource kind and name")
		return
	}
	resourceKind := os.Args[1]
	resourceName := os.Args[2]

	ctx := context.Background()

	var spotify ports.IExtractor

	auth, token := spotifyAuth(ctx)
	youtubeService := youtubeService(ctx)

	youtube := importer.NewYoutubeImporter(youtubeService)

	switch resourceKind {
	case domain.PlaylistKind:
		spotify = extractor.NewSpotifyPlaylistExtractor(ctx, auth, token)

	case domain.AlbumKind:
		fmt.Println("Album extraction not supported")
		return
	}

	migration := usecases.NewMigration(spotify, youtube)

	if ok, err := migration.Migrate(ctx, resourceName); err != nil {
		fmt.Println("Error migrating resource:", err)
	} else if ok {
		fmt.Printf("%s migrated successfully: %s\n", resourceKind, resourceName)
	}
}

func spotifyAuth(ctx context.Context) (*spotifyauth.Authenticator, *oauth2.Token) {
	auth := spotifyauth.New(
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate),
	)

	random := make([]byte, 16)

	_, err := rand.Read(random)
	if err != nil {
		panic("could not generate code")
	}

	token, err := auth.Exchange(ctx, string(random))
	if err != nil {
		panic("could not exchange code")
	}

	return auth, token
}

func youtubeService(ctx context.Context) *youtube.Service {
	youtubeService, err := youtube.NewService(ctx, option.WithCredentialsFile("keyfile.json"))
	if err != nil {
		panic("could not get youtube service")
	}
	return youtubeService
}
