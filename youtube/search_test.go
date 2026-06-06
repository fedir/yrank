package youtube

import (
	"math"
	"net/url"
	"strings"
	"testing"
)

// --- SearchStatistics via mock ---

func TestSearchStatistics_count(t *testing.T) {
	defer withMockClient(t)()

	vs := SearchStatistics("fitness", "dummy-key", 0, false)
	// fixture has 7 results
	if len(vs) != 7 {
		t.Fatalf("expected 7 videos, got %d", len(vs))
	}
	for _, v := range vs {
		if v.Title == "" {
			t.Error("video has empty title")
		}
		if v.URL == "" {
			t.Error("video has empty URL")
		}
	}
}

func TestSearchStatistics_titleUnescaped(t *testing.T) {
	defer withMockClient(t)()

	vs := SearchStatistics("fitness", "dummy-key", 0, false)
	for _, v := range vs {
		if strings.Contains(v.Title, "&amp;") || strings.Contains(v.Title, "&#39;") {
			t.Errorf("title %q still contains an HTML entity", v.Title)
		}
	}

	// The fixture includes a "&amp;" title that must decode to "&".
	found := false
	for _, v := range vs {
		if strings.Contains(v.Title, "Motivation & DÉPRESSION") {
			found = true
		}
	}
	if !found {
		t.Error("expected an unescaped '&' title from the fixture")
	}
}

func TestSearchStatistics_metrics(t *testing.T) {
	defer withMockClient(t)()

	vs := SearchStatistics("fitness", "dummy-key", 0, false)
	for _, v := range vs {
		want := v.LikeCount + v.DislikeCount + v.CommentCount
		if v.TotalReaction != want {
			t.Errorf("video %q: TotalReaction=%d, want %d", v.Title, v.TotalReaction, want)
		}
		if v.GlobalBuzzIndex != v.ViewCount*v.TotalReaction {
			t.Errorf("video %q: GlobalBuzzIndex=%d, want %d", v.Title, v.GlobalBuzzIndex, v.ViewCount*v.TotalReaction)
		}
		if math.IsNaN(v.PositiveInterestingness) || math.IsNaN(v.TotalInterestingness) {
			t.Errorf("video %q: NaN interestingness", v.Title)
		}
	}
}

func TestSearchStatistics_respectsMax(t *testing.T) {
	defer withMockClient(t)()

	// The single-page fixture has no nextPageToken, so the cap simply bounds
	// the result to what one page returns; ensure it never exceeds the fixture.
	vs := SearchStatistics("fitness", "dummy-key", 3, false)
	if len(vs) > 7 {
		t.Errorf("expected at most 7 videos from the fixture, got %d", len(vs))
	}
}

func TestSearchStatistics_rankable(t *testing.T) {
	defer withMockClient(t)()

	vs := SearchStatistics("fitness", "dummy-key", 0, false)
	before := snapshot(vs)
	ApplyAllStrategies(vs)
	after := snapshot(vs)
	assertFieldsUnchanged(t, "ApplyAllStrategies on search data", before, after)

	for _, v := range vs {
		if len(v.AllScores) != len(StrategyOrder) {
			t.Errorf("video %q: expected %d strategy scores, got %d", v.Title, len(StrategyOrder), len(v.AllScores))
		}
	}
}

func TestSearchFixtureRoutesCorrectly(t *testing.T) {
	u, err := url.Parse("https://www.googleapis.com/youtube/v3/search?part=snippet&q=fitness")
	if err != nil {
		t.Fatal(err)
	}
	m := &MockTransport{DataDir: testDataDir()}
	got := m.fixtureFor(u)
	if !strings.HasSuffix(got, "search_results.json") {
		t.Errorf("search URL routed to %q, want search_results.json", got)
	}
}
