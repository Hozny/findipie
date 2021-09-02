package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	r.HandleFunc("/search-playlists", searchPlaylists(ctx, spotifyClient))
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

type SearchPlaylistRequest struct {
	PlaylistName string `json:"playlistName"`
}

type SearchPlaylistResult struct {
	PlaylistName   string   `json:"username"`
	PlaylistID     string   `json:"playlistID"`
	OwnerName      string   `json:"ownerName"`
	OwnerID        string   `json:"ownerID"`
	PlaylistUrl    string   `json:"playlistURL"`
	PlaylistImages []string `json:"playlistImages"`
}

type SearchPlaylistsResponse struct {
	ResultsLength int                    `json:"resultsLength"`
	Playlists     []SearchPlaylistResult `json:"playlists"`
}

func searchPlaylists(ctx context.Context, spotifyClient *spotify.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := ioutil.ReadAll(io.LimitReader(r.Body, 12039123))
		if err != nil {
			fmt.Println("It broke")
		}

		var searchRequest SearchPlaylistRequest
		json.Unmarshal(reqBody, &searchRequest)
		if err != nil {
			fmt.Printf("Invalid request %+v", reqBody)
			// TODO: return error code
			return
		}
		playlistName := string(searchRequest.PlaylistName)
		fmt.Printf("Received search playlist request with playlist name: %s\n", playlistName)

		results, err := spotifyClient.Search(ctx, playlistName, spotify.SearchTypePlaylist)
		if err != nil {
			log.Fatal(err)
		}

		// handle playlist results
		var searchResults []SearchPlaylistResult
		for _, item := range results.Playlists.Playlists {
			var images []string
			for _, image := range item.Images {
				images = append(images, image.URL)
			}
			searchResults = append(searchResults, SearchPlaylistResult{
				PlaylistName:   item.Name,
				OwnerName:      item.Owner.DisplayName,
				OwnerID:        item.Owner.ID,
				PlaylistUrl:    item.ExternalURLs["spotify"],
				PlaylistImages: images,
			})
		}
		response := SearchPlaylistsResponse{
			ResultsLength: len(searchResults),
			Playlists:     searchResults,
		}
		json.NewEncoder(w).Encode(response)
	}
}
