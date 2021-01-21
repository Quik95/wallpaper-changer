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

	metadata := make(chan []WallpaperMetadata, pages)
	failures := make(chan error, 0)
	var seed string

	// Fetch first page of results manually to get a seed
	url := applyParameters(wallpaperConf, 1, "")
	resp, err := fetch(url)
	if err != nil {
		return nil, err
	}
	var met []WallpaperMetadata
	parseJSON(resp, &met, &seed)
	metadata <- met

	// Fetch rest of pages using previously fetched seed
	for i := 2; i <= pages; i++ {
		go func(i int, seed string) {
			url := applyParameters(wallpaperConf, i, seed)
			resp, err := fetch(url)
			if err != nil {
				failures <- err
				return
			}

			var met []WallpaperMetadata
			parseJSON(resp, &met, &seed)
			metadata <- met
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

func parseJSON(in []byte, out *([]WallpaperMetadata), seed *string) {
	var results wallhavenResponse
	json.Unmarshal(in, &results)
	*seed = results.Meta.Seed
	*out = append(*out, results.Data...)
}

func fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to download metadata: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func applyParameters(cfg SearchConfig, pageNumber int, seed string) string {
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
	v.Set("page", fmt.Sprint(pageNumber))
	v.Set("q", cfg.Query)

	if seed != "" {
		v.Set("seed", seed)
	}

	query := v.Encode()

	url := url.URL{
		Scheme:   "https",
		Host:     "wallhaven.cc",
		Path:     "api/v1/search",
		RawQuery: query,
	}
	return url.String()
}
