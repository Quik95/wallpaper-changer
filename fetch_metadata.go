package wallpaperchanger

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/levigross/grequests"
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
func FetchMetadata(args *cli.Context) *[]WallpaperMetadata {
	url := applyParameters(args)
	resp, err := grequests.Get(url, nil)
	if err != nil {
		log.Fatalln("Failed to make request: ", err)
	}

	var metadata wallhavenResponse
	json.Unmarshal(resp.Bytes(), &metadata)

	return &metadata.Data
}

func applyParameters(args *cli.Context) string {
	v := url.Values{}
	v.Set("categories", args.String("categories"))
	v.Set("purity", args.String("purity"))
	v.Set("sorting", args.String("sorting"))
	v.Set("order", args.String("order"))
	v.Set("topRange", args.String("top-range"))
	v.Set("atleast", args.String("atleast"))
	v.Set("resolutions", args.String("resolutions"))
	v.Set("ratios", args.String("ratios"))
	query := v.Encode()

	url := url.URL{
		Scheme:   "https",
		Host:     "wallhaven.cc",
		Path:     "api/v1/search",
		RawQuery: query,
	}
	return url.String()
}
