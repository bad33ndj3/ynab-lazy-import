package downloaddirectory

import (
	"fmt"
	"os"
	"os/user"
)

var errFailedToGetPath error = fmt.Errorf("failed to get path")

func DownloadDirectory() (*string, error) {
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
