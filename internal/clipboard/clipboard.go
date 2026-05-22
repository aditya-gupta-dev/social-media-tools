package clipboard

import (
	"strings"

	"github.com/atotto/clipboard"
)

func GetURL() string {
	text, err := clipboard.ReadAll()
	if err != nil {
		return ""
	}
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "http://") || strings.HasPrefix(text, "https://") {
		return text
	}
	return ""
}
