package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"spotify_migration/adapters/extractor"
	"spotify_migration/adapters/importer"
	"spotify_migration/domain"
	"spotify_migration/ports"
	"spotify_migration/usecases"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
)

const (
	COMM_PROTOCOl     = "https://"
	REDIRECT_HOST     = "127.0.0.1:8080"
	CODE_GET_RESOURCE = "/get-code"
	REDIRECT_URL      = COMM_PROTOCOl + REDIRECT_HOST + CODE_GET_RESOURCE
)

func main() {
	if len(os.Args) < 3 {
		log.Println("Please provide a resource kind and name")
		return
	}
	resourceKind := os.Args[1]
	resourceName := os.Args[2]

	if resourceKind != domain.PlaylistKind {
		log.Println("selected migration not supported")
		return
	}

	ctx := context.Background()

	var spotify ports.IExtractor

	auth, token := spotifyAuth(ctx)
	youtubeService := youtubeService(ctx)

	youtube := importer.NewYoutubeImporter(youtubeService)

	switch resourceKind {
	case domain.PlaylistKind:
		spotify = extractor.NewSpotifyPlaylistExtractor(ctx, auth, token)
	}

	migration := usecases.NewMigration(spotify, youtube)

	if ok, err := migration.Migrate(ctx, resourceName); err != nil {
		log.Println("Error migrating resource:", err)
	} else if ok {
		log.Printf("%s migrated successfully: %s\n", resourceKind, resourceName)
	}
}

func spotifyAuth(ctx context.Context) (*spotifyauth.Authenticator, *oauth2.Token) {
	waitToReceiveCode := make(chan string)
	var code string

	http.HandleFunc(CODE_GET_RESOURCE, func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		waitToReceiveCode <- values.Get("code")
	})

	go http.ListenAndServeTLS(REDIRECT_HOST, "server.crt", "server.key", nil)

	auth := spotifyauth.New(
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate),
		spotifyauth.WithRedirectURL(REDIRECT_URL),
	)

	log.Println("https://accounts.spotify.com/authorize?client_id=" + os.Getenv("SPOTIFY_ID") + "&response_type=code&redirect_uri=" + REDIRECT_URL)

	code = <-waitToReceiveCode

	token, err := auth.Exchange(ctx, code)
	if err != nil {
		log.Fatal("could not exchange code")
	}

	return auth, token
}

func youtubeService(ctx context.Context) *youtube.Service {
	youtubeService, err := youtube.NewService(ctx)
	if err != nil {
		log.Fatal("could not get youtube service")
	}
	return youtubeService
}
