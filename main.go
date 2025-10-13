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
	"strings"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const (
	COMM_PROTOCOl             = "https://"
	REDIRECT_HOST             = "127.0.0.1:8080"
	GET_SPOTIFY_CODE_RESOURCE = "/get-spotify-code"
	GET_YOUTUBE_CODE_RESOURCE = "/get-youtube-code"
	REDIRECT_SPOTIFY_URL      = COMM_PROTOCOl + REDIRECT_HOST + GET_SPOTIFY_CODE_RESOURCE
	REDIRECT_YOUTUBE_URL      = COMM_PROTOCOl + REDIRECT_HOST + GET_YOUTUBE_CODE_RESOURCE
)

var (
	YOUTUBE_SCOPES = []string{
		"https://www.googleapis.com/auth/youtube",
		"https://www.googleapis.com/auth/youtube.channel-memberships.creator",
		"https://www.googleapis.com/auth/youtube.force-ssl",
		"https://www.googleapis.com/auth/youtube.readonly",
		"https://www.googleapis.com/auth/youtube.upload",
		"https://www.googleapis.com/auth/youtubepartner",
		"https://www.googleapis.com/auth/youtubepartner-channel-audit",
	}
)

const (
	CERT_FILE = "server.crt"
	KEY_FILE  = "server.key"
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
	var spotify ports.IExtractor

	ctx := context.Background()
	waitSpotifyCode := make(chan string)
	defer close(waitSpotifyCode)
	waitYoutubeCode := make(chan string)
	defer close(waitYoutubeCode)

	startRedirectServer(map[string]chan<- string{
		GET_SPOTIFY_CODE_RESOURCE: waitSpotifyCode,
		GET_YOUTUBE_CODE_RESOURCE: waitYoutubeCode,
	})

	log.Println("Authorize Spotify: https://accounts.spotify.com/authorize?client_id=" + os.Getenv("SPOTIFY_ID") + "&response_type=code&redirect_uri=" + REDIRECT_SPOTIFY_URL)
	log.Println("Authorize YouTube: https://accounts.google.com/o/oauth2/auth?response_type=code&client_id=" + os.Getenv("YOUTUBE_ID") + "&redirect_uri=" + REDIRECT_YOUTUBE_URL + "&access_type=offline&scope=" + strings.Join(YOUTUBE_SCOPES, "+"))

	auth, token := spotifyAuth(ctx, waitSpotifyCode)
	youtubeService := youtubeService(ctx, waitYoutubeCode)

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

func spotifyAuth(ctx context.Context, waitSpotifyCode <-chan string) (*spotifyauth.Authenticator, *oauth2.Token) {
	auth := spotifyauth.New(
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate),
		spotifyauth.WithRedirectURL(REDIRECT_SPOTIFY_URL),
	)

	code := <-waitSpotifyCode

	token, err := auth.Exchange(ctx, code)
	if err != nil {
		log.Fatal("could not exchange spotify code")
	}

	return auth, token
}

func youtubeService(ctx context.Context, waitYoutubeCode <-chan string) *youtube.Service {
	config := &oauth2.Config{
		ClientID:     os.Getenv("YOUTUBE_ID"),
		ClientSecret: os.Getenv("YOUTUBE_SECRET"),
		RedirectURL:  REDIRECT_YOUTUBE_URL,
		Endpoint:     google.Endpoint,
	}

	code := <-waitYoutubeCode

	token, err := config.Exchange(ctx, code)
	if err != nil {
		log.Fatal("could not exchange youtube code")
	}

	youtubeService, err := youtube.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))

	if err != nil {
		log.Fatal("could not get youtube service")
	}
	return youtubeService
}

func startRedirectServer(channelByResources map[string]chan<- string) {
	waitSpotifyCode := channelByResources[GET_SPOTIFY_CODE_RESOURCE]
	waitYoutubeCode := channelByResources[GET_YOUTUBE_CODE_RESOURCE]

	http.HandleFunc(GET_SPOTIFY_CODE_RESOURCE, func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		waitSpotifyCode <- values.Get("code")
		log.Println("spotify code received!")
	})

	http.HandleFunc(GET_YOUTUBE_CODE_RESOURCE, func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		waitYoutubeCode <- values.Get("code")
		log.Println("youtube code received!")
	})

	server := &http.Server{Addr: REDIRECT_HOST}

	go server.ListenAndServeTLS(CERT_FILE, KEY_FILE)
}
