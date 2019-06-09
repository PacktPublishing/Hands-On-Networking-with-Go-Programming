package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

var (
	baseURL = "localhost:9999"
	conf    = &oauth2.Config{
		RedirectURL:  "http://" + baseURL + "/SpotifyCallback",
		ClientID:     "xxxxxxxx", // See docs at https://developer.spotify.com/documentation/general/guides/app-settings/#register-your-app
		ClientSecret: "xxxxxxxx",
		Scopes:       []string{"user-library-read"}, // access to read what I've saved within Spotify
		Endpoint:     spotify.Endpoint,
	}
	state = uuid.Must(uuid.NewV4()).String()
)

func main() {
	// Start a Web server to collect the code and push it to the channel.
	tokensChan := make(chan string, 1)
	s := http.Server{
		Addr: baseURL,
	}
	s.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if actualState := r.URL.Query().Get("state"); actualState != state {
			http.Error(w, "unexpected authentication value", http.StatusUnauthorized)
			return
		}
		tokensChan <- r.URL.Query().Get("code")
		close(tokensChan)
		w.Write([]byte("Spotify authentication complete, you can close this window."))
	})
	go func() {
		err := s.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatalf("error starting web server: %v", err)
		}
	}()

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	ctx := context.Background()
	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	code := <-tokensChan
	s.Close()
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}
	if !tok.Valid() {
		log.Fatal("The token is invalid")
	}

	client := conf.Client(ctx, tok)
	albums, err := getAllAlbums(client)
	if err != nil {
		log.Fatal(err)
	}
	for _, alb := range albums {
		fmt.Printf("%v, %v\n", getArtistName(alb), alb.Name)
	}
}

func getArtistName(a album) string {
	var names []string
	for _, ar := range a.Artists {
		names = append(names, ar.Name)
	}
	return strings.Join(names, ", ")
}

func getAllAlbums(c *http.Client) (albums []album, err error) {
	var url = "https://api.spotify.com/v1/me/albums"

	for {
		resp, err := c.Get(url)
		if err != nil {
			return albums, fmt.Errorf("error getting Spotify albums: %v", err)
		}
		bdy, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return albums, fmt.Errorf("error reading Spotify albums: %v", err)
		}
		var ar albumResponse
		err = json.Unmarshal(bdy, &ar)
		if err != nil {
			return albums, fmt.Errorf("error understanding Spotify albums: %v", err)
		}
		for _, itm := range ar.Items {
			albums = append(albums, itm.Album)
		}
		url = ar.Next
		if url == "" {
			break
		}
	}
	return
}

type albumResponse struct {
	Href  string `json:"href"`
	Items []struct {
		AddedAt time.Time `json:"added_at"`
		Album   album     `json:"album"`
	} `json:"items"`
	Limit    int         `json:"limit"`
	Next     string      `json:"next"`
	Offset   int         `json:"offset"`
	Previous interface{} `json:"previous"`
	Total    int         `json:"total"`
}

type album struct {
	AlbumType string `json:"album_type"`
	Artists   []struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href string `json:"href"`
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
		URI  string `json:"uri"`
	} `json:"artists"`
	AvailableMarkets []string `json:"available_markets"`
	Copyrights       []struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"copyrights"`
	ExternalIds struct {
		Upc string `json:"upc"`
	} `json:"external_ids"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Genres []interface{} `json:"genres"`
	Href   string        `json:"href"`
	ID     string        `json:"id"`
	Images []struct {
		Height int    `json:"height"`
		URL    string `json:"url"`
		Width  int    `json:"width"`
	} `json:"images"`
	Label                string `json:"label"`
	Name                 string `json:"name"`
	Popularity           int    `json:"popularity"`
	ReleaseDate          string `json:"release_date"`
	ReleaseDatePrecision string `json:"release_date_precision"`
	TotalTracks          int    `json:"total_tracks"`
	Tracks               struct {
		Href  string `json:"href"`
		Items []struct {
			Artists []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"artists"`
			AvailableMarkets []string `json:"available_markets"`
			DiscNumber       int      `json:"disc_number"`
			DurationMs       int      `json:"duration_ms"`
			Explicit         bool     `json:"explicit"`
			ExternalUrls     struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href        string `json:"href"`
			ID          string `json:"id"`
			IsLocal     bool   `json:"is_local"`
			Name        string `json:"name"`
			PreviewURL  string `json:"preview_url"`
			TrackNumber int    `json:"track_number"`
			Type        string `json:"type"`
			URI         string `json:"uri"`
		} `json:"items"`
		Limit    int         `json:"limit"`
		Next     interface{} `json:"next"`
		Offset   int         `json:"offset"`
		Previous interface{} `json:"previous"`
		Total    int         `json:"total"`
	} `json:"tracks"`
	Type string `json:"type"`
	URI  string `json:"uri"`
}
