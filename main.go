package main

import (
	"context"
    "fmt"
	"log"
	"os"
    "net/http"

    "github.com/gorilla/mux"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
	"github.com/zmb3/spotify/v2"
)


func newRouter() *mux.Router { 
    r := mux.NewRouter()
    r.HandleFunc("/hello", handler).Methods("GET")
    return r
}

func main() {
    // Declaring a new router
    r := newRouter()

    fmt.Println("HELLO")
    r.HandleFunc("/", handler)
    r.HandleFunc("/playlist-search", handlePlaylistSearch)

    http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World!")
}

func handlePlaylistSearch(w http.ResponseWriter, r *http.Request) {
    playlistSearch()
}

func playlistSearch() {
	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)
	// search for playlists and albums containing "holiday"
	results, err := client.Search(ctx, "holiday", spotify.SearchTypePlaylist|spotify.SearchTypeAlbum)
	if err != nil {
		log.Fatal(err)
	}

	// handle album results
	if results.Albums != nil {
		fmt.Println("Albums:")
		for _, item := range results.Albums.Albums {
			fmt.Println("   ", item.Name)
		}
	}
	// handle playlist results
	if results.Playlists != nil {
		fmt.Println("Playlists:")
		for _, item := range results.Playlists.Playlists {
			fmt.Println("   ", item.Name)
		}
	}
}


