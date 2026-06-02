package youtube

import (
	"path/filepath"
	"runtime"
	"testing"
)

func testDataDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "testdata")
}

func withMockClient(t *testing.T) func() {
	t.Helper()
	orig := httpClient
	SetHTTPClient(NewMockClient(testDataDir()))
	return func() { httpClient = orig }
}

// --- PlaylistStatistics via mock ---

func TestPlaylistStatistics_mock_count(t *testing.T) {
	defer withMockClient(t)()

	vs := PlaylistStatistics("PLiVdPopzGBsV7TgjAw9GH43Ck9QCxrw5w", "dummy-key", "", false)
	if len(vs) == 0 {
		t.Fatal("expected videos from mock playlist, got none")
	}
	// fixture has 7 items
	if len(vs) != 7 {
		t.Errorf("expected 7 videos, got %d", len(vs))
	}
}

func TestPlaylistStatistics_mock_fields(t *testing.T) {
	defer withMockClient(t)()

	vs := PlaylistStatistics("PLiVdPopzGBsV7TgjAw9GH43Ck9QCxrw5w", "dummy-key", "", false)
	for _, v := range vs {
		if v.Title == "" {
			t.Error("video has empty title")
		}
		if v.URL == "" {
			t.Error("video has empty URL")
		}
		if v.ViewCount < 0 {
			t.Errorf("video %q has negative ViewCount", v.Title)
		}
	}
}

func TestPlaylistStatistics_mock_metrics(t *testing.T) {
	defer withMockClient(t)()

	vs := PlaylistStatistics("PLiVdPopzGBsV7TgjAw9GH43Ck9QCxrw5w", "dummy-key", "", false)
	for _, v := range vs {
		// TotalReaction must equal likes + dislikes + comments
		want := v.LikeCount + v.DislikeCount + v.CommentCount
		if v.TotalReaction != want {
			t.Errorf("video %q: TotalReaction=%d, want %d", v.Title, v.TotalReaction, want)
		}
		// GlobalBuzzIndex must equal views * totalReaction
		wantGBI := v.ViewCount * v.TotalReaction
		if v.GlobalBuzzIndex != wantGBI {
			t.Errorf("video %q: GlobalBuzzIndex=%d, want %d", v.Title, v.GlobalBuzzIndex, wantGBI)
		}
		// No NaN / Inf in float fields
		if v.PositiveInterestingness != v.PositiveInterestingness {
			t.Errorf("video %q: PositiveInterestingness is NaN", v.Title)
		}
		if v.TotalInterestingness != v.TotalInterestingness {
			t.Errorf("video %q: TotalInterestingness is NaN", v.Title)
		}
	}
}

func TestPlaylistStatistics_mock_safeDiv(t *testing.T) {
	defer withMockClient(t)()

	// All fixture videos have views > 0, but safeDiv must never produce Inf
	vs := PlaylistStatistics("PLiVdPopzGBsV7TgjAw9GH43Ck9QCxrw5w", "dummy-key", "", false)
	for _, v := range vs {
		if v.PositiveInterestingness > 1e15 || v.PositiveInterestingness < -1e15 {
			t.Errorf("video %q: PositiveInterestingness looks like Inf: %f", v.Title, v.PositiveInterestingness)
		}
	}
}

func TestPlaylistStatistics_mock_allStrategies(t *testing.T) {
	defer withMockClient(t)()

	vs := PlaylistStatistics("PLiVdPopzGBsV7TgjAw9GH43Ck9QCxrw5w", "dummy-key", "", false)
	ApplyAllStrategies(vs)

	for _, v := range vs {
		if len(v.AllScores) != len(StrategyOrder) {
			t.Errorf("video %q: expected %d strategy scores, got %d", v.Title, len(StrategyOrder), len(v.AllScores))
		}
		for _, slug := range StrategyOrder {
			score := v.AllScores[slug]
			if score != score {
				t.Errorf("video %q strategy %q: NaN score", v.Title, slug)
			}
			if score < 0 {
				t.Errorf("video %q strategy %q: negative score %f", v.Title, slug, score)
			}
		}
	}
}

func TestPlaylistStatistics_mock_sortPreservesData(t *testing.T) {
	defer withMockClient(t)()

	vs := PlaylistStatistics("PLiVdPopzGBsV7TgjAw9GH43Ck9QCxrw5w", "dummy-key", "", false)
	before := snapshot(vs)
	SortBy(vs, "likes")
	after := snapshot(vs)
	assertFieldsUnchanged(t, "SortBy/likes on mock data", before, after)
}
