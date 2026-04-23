package service

import (
	"path/filepath"
	"testing"

	"github.com/LuuDinhTheTai/tzone/infrastructure/configuration"
)

func TestNewAIChatServiceLoadsCatalog(t *testing.T) {
	dataPath := filepath.Clean(filepath.Join("..", "..", "phoneExample.json"))
	svc, err := NewAIChatService(configuration.AIConfig{
		PhoneDataPath: dataPath,
	})
	if err != nil {
		t.Fatalf("expected catalog load success, got error: %v", err)
	}
	if len(svc.catalog) == 0 {
		t.Fatal("expected non-empty catalog")
	}
}

func TestPickCandidatesFindsRelevantDevice(t *testing.T) {
	dataPath := filepath.Clean(filepath.Join("..", "..", "phoneExample.json"))
	svc, err := NewAIChatService(configuration.AIConfig{
		PhoneDataPath: dataPath,
	})
	if err != nil {
		t.Fatalf("failed to init service: %v", err)
	}

	candidates := svc.pickCandidates("toi can iphone", 10)
	if len(candidates) == 0 {
		t.Fatal("expected candidates for iphone query")
	}
}

func TestBuildCardsSkipsUnknownIDs(t *testing.T) {
	dataPath := filepath.Clean(filepath.Join("..", "..", "phoneExample.json"))
	svc, err := NewAIChatService(configuration.AIConfig{
		PhoneDataPath: dataPath,
	})
	if err != nil {
		t.Fatalf("failed to init service: %v", err)
	}

	firstID := svc.catalog[0].ID
	cards := svc.buildCards([]string{"not-found", firstID, firstID}, 4)
	if len(cards) != 1 {
		t.Fatalf("expected 1 card, got %d", len(cards))
	}
	if cards[0].ID != firstID {
		t.Fatalf("expected card ID %s, got %s", firstID, cards[0].ID)
	}
}

func TestNormalizeYouTubeURL(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want bool
	}{
		{name: "youtube watch", raw: "https://www.youtube.com/watch?v=abc", want: true},
		{name: "youtu short", raw: "https://youtu.be/abc", want: true},
		{name: "m youtube", raw: "https://m.youtube.com/watch?v=abc", want: true},
		{name: "non youtube", raw: "https://vimeo.com/123", want: false},
		{name: "invalid", raw: "not-a-url", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeYouTubeURL(tt.raw)
			if (got != "") != tt.want {
				t.Fatalf("normalizeYouTubeURL(%q)=%q, want valid=%v", tt.raw, got, tt.want)
			}
		})
	}
}

func TestSanitizeYouTubeVideos(t *testing.T) {
	raw := []geminiVideoReplyRow{
		{Title: "Review 1", URL: "https://www.youtube.com/watch?v=abc"},
		{Title: "Review 1 duplicate", URL: "https://www.youtube.com/watch?v=abc"},
		{Title: "Bad host", URL: "https://example.com/video"},
		{Title: "", URL: "https://youtu.be/xyz"},
	}

	videos := sanitizeYouTubeVideos(raw, 3)
	if len(videos) != 2 {
		t.Fatalf("expected 2 valid videos, got %d", len(videos))
	}
	if videos[1].Title == "" {
		t.Fatal("expected fallback title for empty title")
	}
}

func TestBuildSingleYouTubeFallbackVideo(t *testing.T) {
	videos := buildSingleYouTubeFallbackVideo("Samsung S24")
	if len(videos) != 1 {
		t.Fatalf("expected 1 fallback video, got %d", len(videos))
	}
	if videos[0].URL == "" || videos[0].Title == "" {
		t.Fatal("expected fallback video to have title and url")
	}
}
