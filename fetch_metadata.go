package wallpaperchanger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/urfave/cli/v2"
)

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
func FetchMetadata(args *cli.Context) (*[]WallpaperMetadata, error) {
	pages := args.Int("pages")

	metadata := make(chan []WallpaperMetadata, pages)
	failures := make(chan error, 0)
	var seed string

	url := applyParameters(args, 1, "")
	resp, err := fetch(url)
	if err != nil {
		return nil, err
	}
	var met []WallpaperMetadata
	parseJSON(resp, &met, &seed)
	metadata <- met

	for i := 2; i <= pages; i++ {
		go func(i int, seed string) {
			url := applyParameters(args, i, seed)
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

func applyParameters(args *cli.Context, pageNumber int, seed string) string {
	v := url.Values{}
	v.Set("categories", args.String("categories"))
	v.Set("purity", args.String("purity"))
	v.Set("sorting", args.String("sorting"))
	v.Set("order", args.String("order"))
	v.Set("topRange", args.String("top-range"))
	v.Set("atleast", args.String("atleast"))
	v.Set("resolutions", args.String("resolutions"))
	v.Set("ratios", args.String("ratios"))
	v.Set("apikey", args.String("api-key"))
	v.Set("page", fmt.Sprint(pageNumber))

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
