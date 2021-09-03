package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

var MAX_REQUEST_SIZE = 1048576

func getRouter(ctx context.Context, spotifyClient *spotify.Client) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/search-playlists", SearchPlaylists(ctx, spotifyClient))
	return r
}

func main() {
	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
		// TODO: handle this error
		return
	}
	httpClient := spotifyauth.New().Client(ctx, token)
	spotifyClient := spotify.New(httpClient)

	router := getRouter(ctx, spotifyClient)

	http.Handle("/", router)
	http.ListenAndServe("127.0.0.1:8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HELLO")
}
