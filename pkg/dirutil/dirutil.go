package dirutil

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

var errFailedToGetPath error = fmt.Errorf("failed to get path")

func FilePathWalkDir(root, ext string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ext {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func DownloadPath() (*string, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("failed to get user for download path: %w", err)
	}
	downloadDir := fmt.Sprintf("%s/%s", usr.HomeDir, "Downloads")
	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to get download path: %w", errFailedToGetPath)
	}
	return &downloadDir, nil
}
