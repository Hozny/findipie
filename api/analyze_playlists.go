package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/zmb3/spotify/v2"
)

type SsearchPlaylistsRequest struct {
	PlaylistName string `json:"playlistName"`
}

type SsearchPlaylistsResult struct {
	PlaylistName   string   `json:"username"`
	PlaylistID     string   `json:"playlistID"`
	OwnerName      string   `json:"ownerName"`
	OwnerID        string   `json:"ownerID"`
	PlaylistUrl    string   `json:"playlistURL"`
	PlaylistImages []string `json:"playlistImages"`
}

type SsearchPlaylistsResponse struct {
	ResultsLength int                     `json:"resultsLength"`
	Playlists     []SearchPlaylistsResult `json:"playlists"`
}

func SdarchPlaylists(ctx context.Context, spotifyClient *spotify.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := ioutil.ReadAll(io.LimitReader(r.Body, 12039123))
		if err != nil {
			fmt.Println("It broke")
		}

		var searchRequest SearchPlaylistsRequest
		json.Unmarshal(reqBody, &searchRequest)
		if err != nil {
			fmt.Printf("Invalid request %+v", reqBody)
			// TODO: handle error
			return
		}
		playlistName := string(searchRequest.PlaylistName)
		fmt.Printf("Received search playlist request with playlist name: %s\n", playlistName)

		results, err := spotifyClient.Search(ctx, playlistName, spotify.SearchTypePlaylist)
		if err != nil {
			log.Fatal(err)
		}

		// handle playlist results
		var searchResults []SearchPlaylistsResult
		for _, item := range results.Playlists.Playlists {
			var images []string
			for _, image := range item.Images {
				images = append(images, image.URL)
			}
			searchResults = append(searchResults, SearchPlaylistsResult{
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


