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
}

// FetchMetadata makes request with passed parameters to wallhaven.cc api and returns response as json
func FetchMetadata(args *cli.Context) (*[]WallpaperMetadata, error) {
	pages := args.Int("pages")

	var metadata []WallpaperMetadata

	for i := 1; i <= pages; i++ {
		url := applyParameters(args, i)
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("Failed to download metadata: %v", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var results wallhavenResponse
		json.Unmarshal(body, &results)
		metadata = append(metadata, results.Data...)
	}

	return &metadata, nil
}

func applyParameters(args *cli.Context, pageNumber int) string {
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
	query := v.Encode()

	url := url.URL{
		Scheme:   "https",
		Host:     "wallhaven.cc",
		Path:     "api/v1/search",
		RawQuery: query,
	}
	return url.String()
}
