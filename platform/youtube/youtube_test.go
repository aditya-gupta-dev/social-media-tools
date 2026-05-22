package youtube

import "testing"

func TestYouTubeMatch(t *testing.T) {
	yt := &YouTube{}
	tests := []struct {
		url   string
		match bool
	}{
		{"https://www.youtube.com/watch?v=dQw4w9WgXcQ", true},
		{"https://youtu.be/dQw4w9WgXcQ", true},
		{"https://www.youtube.com/shorts/XYZ", true},
		{"https://google.com", false},
		{"not a url", false},
	}

	for _, tt := range tests {
		if got := yt.Match(tt.url); got != tt.match {
			t.Errorf("Match(%q) = %v; want %v", tt.url, got, tt.match)
		}
	}
}
