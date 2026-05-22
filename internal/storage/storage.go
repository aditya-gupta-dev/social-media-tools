package storage

import (
	"os"
	"path/filepath"
)

const AppName = "social-media-tools"

func GetDownloadPath(mediaType string, platformName string, usePWD bool) string {
	if usePWD {
		pwd, _ := os.Getwd()
		return filepath.Join(pwd, platformName)
	}

	home, _ := os.UserHomeDir()
	var folder string
	switch mediaType {
	case "audio":
		folder = filepath.Join(home, "Music", AppName, platformName)
	case "video":
		folder = filepath.Join(home, "Videos", AppName, platformName)
	case "thumbnail":
		folder = filepath.Join(home, "Pictures", AppName, platformName)
	default:
		folder = filepath.Join(home, "Downloads", AppName, platformName)
	}

	_ = os.MkdirAll(folder, 0755)
	return folder
}
