package wallpaperchanger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// SearchConfig describes user defined wallpaper search options
type SearchConfig struct {
	Categories  string
	Purity      string
	Sorting     string
	Order       string
	TopRange    string
	Atleast     string
	Resolutions string
	Ratios      string
	Pages       int
	Seed        string
	Query       string
	APIKey      string
}

// WallpaperMetadata represents information needed to download a single wallpaper
type WallpaperMetadata struct {
	Path     string
	ID       string
	FileType string `json:"file_type"`
}

type wallhavenResponse struct {
	Data []WallpaperMetadata
	Meta struct{ Seed string }
}

// FetchMetadata makes request with passed parameters to wallhaven.cc api and returns response as json
func FetchMetadata(wallpaperConf SearchConfig) (*[]WallpaperMetadata, error) {
	pages := wallpaperConf.Pages
	seed := wallpaperConf.Seed

	metadata := make(chan []WallpaperMetadata, pages)
	failures := make(chan error, 0)

	// Fetch first page of results manually to get a seed
	url := applyParameters(wallpaperConf)
	pageOne, err := fetchPage(url, 1, &seed)
	if err != nil {
		return nil, err
	}
	metadata <- pageOne

	// Fetch rest of pages using previously fetched seed
	for i := 2; i <= pages; i++ {
		go func(i int, seed string) {
			resp, err := fetchPage(url, i, &seed)
			if err != nil {
				failures <- err
				return
			}

			metadata <- resp
		}(i, seed)
	}

	select {
	case err := <-failures:
		return nil, err
	default:
		{
			var wallpapers []WallpaperMetadata
			for i := 0; i < pages; i++ {
				wallpapers = append(wallpapers, <-metadata...)
			}
			close(metadata)
			return &wallpapers, nil
		}
	}
}

// fetchPage fetches given page of wallpapers and sets the seed to the one returned by server
func fetchPage(fetchURL url.URL, page int, seed *string) ([]WallpaperMetadata, error) {
	query, err := url.ParseQuery(fetchURL.RawQuery)
	if err != nil {
		return nil, err
	}
	query.Set("page", fmt.Sprint(page))
	query.Set("seed", *seed)
	fetchURL.RawQuery = query.Encode()

	resp, err := http.Get(fetchURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var metadata wallhavenResponse
	json.Unmarshal(body, &metadata)

	// has no effect when user provided a seed
	// but if it wasn't provided changes it to the one
	// returned by server
	// seed doesn't change for subsequent results
	seed = &metadata.Meta.Seed

	return metadata.Data, nil
}

func applyParameters(cfg SearchConfig) url.URL {
	v := url.Values{}
	v.Set("categories", cfg.Categories)
	v.Set("purity", cfg.Purity)
	v.Set("sorting", cfg.Sorting)
	v.Set("order", cfg.Order)
	v.Set("topRange", cfg.TopRange)
	v.Set("atleast", cfg.Atleast)
	v.Set("resolutions", cfg.Resolutions)
	v.Set("ratios", cfg.Ratios)
	v.Set("apikey", cfg.APIKey)
	v.Set("q", cfg.Query)

	query := v.Encode()

	url := url.URL{
		Scheme:   "https",
		Host:     "wallhaven.cc",
		Path:     "api/v1/search",
		RawQuery: query,
	}
	return url
}
