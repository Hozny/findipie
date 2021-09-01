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

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
	"github.com/zmb3/spotify/v2"
)

var MAX_REQUEST_SIZE = 1048576
var SPOTIFY_ID = "bf253aa4136a4f7f9caa0c9dabfb165c"
var SPOTIFY_SECRET = "8cbaaea4aaeb451d8f4e3cbf9318e427"


func getRouter(ctx context.Context, spotifyClient *spotify.Client) *mux.Router{ 
    r := mux.NewRouter()
    r.HandleFunc("/", home)
    r.HandleFunc("/search-users", searchUsers(ctx, spotifyClient))
    return r
}

func main() {
	ctx := context.Background()
	config := &clientcredentials.Config{
		// ClientID:     os.Getenv("SPOTIFY_ID"),
		// ClientSecret: os.Getenv("SPOTIFY_SECRET"),
        ClientID:     SPOTIFY_ID,
		ClientSecret: SPOTIFY_SECRET,
		TokenURL:     spotifyauth.TokenURL,
	}
    fmt.Println(os.Getenv("SPOTIFY_ID"), os.Getenv("SPOTIFY_SECRET"))
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
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

type SearchUsersRequest struct {
    Username string `json:username`
}

func searchUsers(ctx context.Context, spotifyClient *spotify.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        reqBody, err := ioutil.ReadAll(io.LimitReader(r.Body, 12039123))
        if err != nil {
            fmt.Println("It broke")
        }

        var searchRequest SearchUsersRequest
        json.Unmarshal(reqBody, &searchRequest)
        if err != nil {
            fmt.Println("Invalid request %+v", reqBody)
            // TODO: return error code
            return
        }
        username := string(searchRequest.Username)
        fmt.Printf("Received search user request with username: %s\n", username)


        spotifyPlaylistSearch(ctx, username)
        // fmt.Fprintf(w, "%+v", string(reqBody)) 
        // fmt.Printf("%+v", string(reqBody)) 
    }
}

func spotifyPlaylistSearch(ctx context.Context, username string) {
	config := &clientcredentials.Config{
		// ClientID:     os.Getenv("SPOTIFY_ID"),
		// ClientSecret: os.Getenv("SPOTIFY_SECRET"),
        ClientID:     SPOTIFY_ID,
		ClientSecret: SPOTIFY_SECRET,
		TokenURL:     spotifyauth.TokenURL,
	}
    fmt.Println(os.Getenv("SPOTIFY_ID"), os.Getenv("SPOTIFY_SECRET"))
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)

	// search for playlists and albums containing "holiday"
	results, err := client.Search(ctx, username, spotify.SearchTypePlaylist)
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
			fmt.Println("   ", item.Name, item.Owner)
		}
	}
    // spotify.get(https://api.spotify.com/v1/users/{user_id}/playlists)
}
