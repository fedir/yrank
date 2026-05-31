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
			PositiveNegativeCoefficient: 2.0, GlobalBuzzIndex: 17000,
			PublishedAt: now.Add(-10 * 24 * time.Hour)},
		{Title: "b", ViewCount: 500, LikeCount: 80, DislikeCount: 1, CommentCount: 50,
			TotalReaction: 131, TotalInterestingness: 0.262, PositiveInterestingness: 0.158,
			PositiveNegativeCoefficient: 40.0, GlobalBuzzIndex: 65500,
			PublishedAt: now.Add(-2 * 24 * time.Hour)},
		{Title: "c", ViewCount: 2000, LikeCount: 30, DislikeCount: 20, CommentCount: 5,
			TotalReaction: 55, TotalInterestingness: 0.0275, PositiveInterestingness: 0.005,
			PositiveNegativeCoefficient: 1.43, GlobalBuzzIndex: 110000,
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
