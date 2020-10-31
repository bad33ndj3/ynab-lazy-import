package dirutil

import (
	"fmt"
	"os"
	"path/filepath"
)

// FilePathWalkDir get all files with a specified extension from a directory.
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

// DownloadPath gets the default download path of the current user.
func DownloadPath() (string, error) {
	dir, err := GetUserDirDirectory("Downloads")
	if err != nil {
		return "", err
	}
	return dir, nil
}

// GetUserDirDirectory gets an extension from the current users directory if it exists.
func GetUserDirDirectory(directory string) (string, error) {
	userDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	downloadDir := fmt.Sprintf("%s/%s", userDir, directory)
	if _, err := os.Stat(downloadDir); err != nil {
		return "", err
	}

	return downloadDir, nil
}
