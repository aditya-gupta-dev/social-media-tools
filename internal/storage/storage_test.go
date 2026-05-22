package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetDownloadPath(t *testing.T) {
	platform := "YouTube"
	
	t.Run("Audio Path", func(t *testing.T) {
		path := GetDownloadPath("audio", platform, false)
		if !strings.Contains(path, "Music") || !strings.Contains(path, AppName) || !strings.Contains(path, platform) {
			t.Errorf("Unexpected audio path: %s", path)
		}
	})

	t.Run("Video Path", func(t *testing.T) {
		path := GetDownloadPath("video", platform, false)
		if !strings.Contains(path, "Videos") || !strings.Contains(path, AppName) || !strings.Contains(path, platform) {
			t.Errorf("Unexpected video path: %s", path)
		}
	})

	t.Run("Thumbnail Path", func(t *testing.T) {
		path := GetDownloadPath("thumbnail", platform, false)
		if !strings.Contains(path, "Pictures") || !strings.Contains(path, AppName) || !strings.Contains(path, platform) {
			t.Errorf("Unexpected thumbnail path: %s", path)
		}
	})

	t.Run("PWD Path", func(t *testing.T) {
		pwd, _ := os.Getwd()
		path := GetDownloadPath("video", platform, true)
		expected := filepath.Join(pwd, platform)
		if path != expected {
			t.Errorf("Expected %s, got %s", expected, path)
		}
	})
}
