package youtube

import (
	"testing"
	"time"
)

func sampleVideos() []VideoStatistics {
	now := time.Now()
	return []VideoStatistics{
		{Title: "a", ViewCount: 1000, LikeCount: 10, DislikeCount: 5, CommentCount: 2,
			TotalReaction: 17, TotalInterestingness: 0.017, PositiveInterestingness: 0.005,
			PositiveNegativeCoefficient: 2.0, GlobalBuzzIndex: 17000, Duration: 120,
			PublishedAt: now.Add(-10 * 24 * time.Hour)},
		{Title: "b", ViewCount: 500, LikeCount: 80, DislikeCount: 1, CommentCount: 50,
			TotalReaction: 131, TotalInterestingness: 0.262, PositiveInterestingness: 0.158,
			PositiveNegativeCoefficient: 40.0, GlobalBuzzIndex: 65500, Duration: 600,
			PublishedAt: now.Add(-2 * 24 * time.Hour)},
		{Title: "c", ViewCount: 2000, LikeCount: 30, DislikeCount: 20, CommentCount: 5,
			TotalReaction: 55, TotalInterestingness: 0.0275, PositiveInterestingness: 0.005,
			PositiveNegativeCoefficient: 1.43, GlobalBuzzIndex: 110000, Duration: 900,
			PublishedAt: now.Add(-30 * 24 * time.Hour)},
	}
}

func TestSortBy_totalInterest(t *testing.T) {
	vs := sampleVideos()
	SortBy(vs, "total-interest")
	if vs[0].Title != "b" {
		t.Errorf("total-interest: expected b first, got %s", vs[0].Title)
	}
}

func TestSortBy_default(t *testing.T) {
	vs := sampleVideos()
	SortBy(vs, "")
	if vs[0].Title != "b" {
		t.Errorf("default: expected b first, got %s", vs[0].Title)
	}
}

func TestSortBy_positiveInterest(t *testing.T) {
	vs := sampleVideos()
	SortBy(vs, "positive-interest")
	if vs[0].Title != "b" {
		t.Errorf("positive-interest: expected b first, got %s", vs[0].Title)
	}
}

func TestSortBy_likes(t *testing.T) {
	vs := sampleVideos()
	SortBy(vs, "likes")
	if vs[0].Title != "b" {
		t.Errorf("likes: expected b first, got %s", vs[0].Title)
	}
}

func TestSortBy_duration(t *testing.T) {
	vs := sampleVideos()
	SortBy(vs, "duration")
	if vs[0].Title != "c" {
		t.Errorf("duration: expected c first (longest), got %s", vs[0].Title)
	}
	if vs[2].Title != "a" {
		t.Errorf("duration: expected a last (shortest), got %s", vs[2].Title)
	}
}

func TestSortBy_totalReaction(t *testing.T) {
	vs := sampleVideos()
	SortBy(vs, "total-reaction")
	if vs[0].Title != "b" {
		t.Errorf("total-reaction: expected b first, got %s", vs[0].Title)
	}
}

func TestSortBy_globalBuzzIndex(t *testing.T) {
	vs := sampleVideos()
	SortBy(vs, "global-buzz-index")
	if vs[0].Title != "c" {
		t.Errorf("global-buzz-index: expected c first, got %s", vs[0].Title)
	}
}

func TestSortBy_pnc(t *testing.T) {
	vs := sampleVideos()
	SortBy(vs, "pnc")
	if vs[0].Title != "b" {
		t.Errorf("pnc: expected b first, got %s", vs[0].Title)
	}
}

func TestSortBy_positiveNegativeCoefficient(t *testing.T) {
	vs := sampleVideos()
	SortBy(vs, "positive-negative-coefficient")
	if vs[0].Title != "b" {
		t.Errorf("positive-negative-coefficient: expected b first, got %s", vs[0].Title)
	}
}

// --- data integrity after sorting ---

// snapshot captures the metric values of each video by title so we can verify
// that sorting reorders without altering any field values.
func snapshot(vs []VideoStatistics) map[string]VideoStatistics {
	m := make(map[string]VideoStatistics, len(vs))
	for _, v := range vs {
		m[v.Title] = v
	}
	return m
}

func assertFieldsUnchanged(t *testing.T, label string, before, after map[string]VideoStatistics) {
	t.Helper()
	for title, b := range before {
		a, ok := after[title]
		if !ok {
			t.Errorf("%s: video %q disappeared after sort", label, title)
			continue
		}
		if b.ViewCount != a.ViewCount || b.LikeCount != a.LikeCount ||
			b.DislikeCount != a.DislikeCount || b.CommentCount != a.CommentCount ||
			b.TotalReaction != a.TotalReaction || b.GlobalBuzzIndex != a.GlobalBuzzIndex ||
			b.TotalInterestingness != a.TotalInterestingness ||
			b.PositiveInterestingness != a.PositiveInterestingness ||
			b.PositiveNegativeCoefficient != a.PositiveNegativeCoefficient {
			t.Errorf("%s: field values changed for video %q after sort", label, title)
		}
	}
}

func TestSortBy_doesNotMutateFields(t *testing.T) {
	modes := []string{
		"total-interest", "positive-interest", "likes",
		"total-reaction", "global-buzz-index", "pnc",
		"positive-negative-coefficient", "",
	}
	for _, mode := range modes {
		vs := sampleVideos()
		before := snapshot(vs)
		SortBy(vs, mode)
		after := snapshot(vs)
		assertFieldsUnchanged(t, "SortBy/"+mode, before, after)
	}
}

func TestApplyStrategy_doesNotMutateFields(t *testing.T) {
	strategies := []string{"viral", "educational", "controversial", "community", "evergreen", "hype"}
	for _, slug := range strategies {
		vs := sampleVideos()
		before := snapshot(vs)
		s := Strategies[slug]
		ApplyStrategy(vs, slug, s.DefaultWeights)
		after := snapshot(vs)
		assertFieldsUnchanged(t, "ApplyStrategy/"+slug, before, after)
	}
}

func TestApplyAllStrategies_doesNotMutateFields(t *testing.T) {
	vs := sampleVideos()
	before := snapshot(vs)
	ApplyAllStrategies(vs)
	after := snapshot(vs)
	assertFieldsUnchanged(t, "ApplyAllStrategies", before, after)
}

func TestApplyAllStrategies_allSlugsPresent(t *testing.T) {
	vs := sampleVideos()
	ApplyAllStrategies(vs)
	for _, v := range vs {
		if len(v.AllScores) != len(StrategyOrder) {
			t.Errorf("video %q: expected %d AllScores entries, got %d", v.Title, len(StrategyOrder), len(v.AllScores))
		}
		for _, slug := range StrategyOrder {
			if _, ok := v.AllScores[slug]; !ok {
				t.Errorf("video %q: missing AllScores[%q]", v.Title, slug)
			}
		}
	}
}

func TestApplyAllStrategies_scoresNonNegative(t *testing.T) {
	vs := sampleVideos()
	ApplyAllStrategies(vs)
	for _, v := range vs {
		for slug, score := range v.AllScores {
			if score < 0 {
				t.Errorf("video %q strategy %q: negative score %f", v.Title, slug, score)
			}
		}
	}
}

func TestApplyAllStrategies_emptySlice(t *testing.T) {
	// must not panic on empty input
	ApplyAllStrategies([]VideoStatistics{})
}

func TestApplyAllStrategies_zeroViews(t *testing.T) {
	// videos with zero views must not produce NaN/Inf scores
	vs := []VideoStatistics{
		{Title: "zero", ViewCount: 0, LikeCount: 0, CommentCount: 0, PublishedAt: time.Now()},
	}
	ApplyAllStrategies(vs)
	for slug, score := range vs[0].AllScores {
		if score != score { // NaN check
			t.Errorf("strategy %q produced NaN score for zero-view video", slug)
		}
	}
}

// --- uploadsPlaylistID ---

func TestUploadsPlaylistID(t *testing.T) {
	tests := []struct {
		channelID string
		want      string
	}{
		{"UCX6OQ3DkcsbYNE6BQcCnHKA", "UUX6OQ3DkcsbYNE6BQcCnHKA"},
		{"UCWeg2Pkate69NFdBeuRFTAw", "UUWeg2Pkate69NFdBeuRFTAw"},
		{"notStartingWithUC", "notStartingWithUC"}, // passthrough
		{"", ""},
	}
	for _, tc := range tests {
		got := uploadsPlaylistID(tc.channelID)
		if got != tc.want {
			t.Errorf("uploadsPlaylistID(%q) = %q, want %q", tc.channelID, got, tc.want)
		}
	}
}
