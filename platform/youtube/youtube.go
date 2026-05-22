package youtube

import (
	"regexp"
)

type YouTube struct{}

func (y *YouTube) Match(url string) bool {
	// Simple regex for YouTube URLs
	re := regexp.MustCompile(`(https?://)?(www\.)?(youtube\.com|youtu\.be)/.+`)
	return re.MatchString(url)
}

func (y *YouTube) GetName() string {
	return "YouTube"
}
