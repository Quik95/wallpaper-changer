package wallpaperchanger

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadWallpaper downloads wallpaper and returns it's save path
func DownloadWallpaper(wlp *WallpaperMetadata, basePath string) (string, error) {
	// determine file type
	ext := ".jpg"
	if wlp.FileType == "image/png" {
		ext = ".png"
	}
	savePath, err := filepath.Abs(filepath.Join(basePath, "Wallhaven-"+wlp.ID+ext))
	if err != nil {
		return "", err
	}

	//skip download if file already exists
	if _, err := os.Stat(savePath); err == nil {
		return savePath, nil
	}

	resp, err := http.Get(wlp.Path)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	out, err := os.Create(savePath)
	if err != nil {
		return "", nil
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return savePath, nil
}
