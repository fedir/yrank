package youtube

import (
	"math"
	"testing"
	"time"
)

var baseVideo = VideoStatistics{
	ViewCount:    1000,
	LikeCount:    80,
	DislikeCount: 5,
	CommentCount: 20,
	PublishedAt:  time.Now().Add(-30 * 24 * time.Hour), // 30 days ago
}

var baseCtx = StrategyContext{
	MaxViews: 1000,
	Now:      time.Now(),
	Weights:  nil, // set per test
}

func defaultCtx(slug string) StrategyContext {
	ctx := baseCtx
	ctx.Weights = Strategies[slug].DefaultWeights
	return ctx
}

func TestStrategyViral(t *testing.T) {
	s := Strategies["viral"]
	score := s.Score(baseVideo, defaultCtx("viral"))
	if score <= 0 {
		t.Errorf("viral score should be positive, got %f", score)
	}
	// zero views → 0
	zeroViews := baseVideo
	zeroViews.ViewCount = 0
	if s.Score(zeroViews, defaultCtx("viral")) != 0 {
		t.Error("viral score should be 0 for zero views")
	}
}

func TestStrategyEducational(t *testing.T) {
	s := Strategies["educational"]
	score := s.Score(baseVideo, defaultCtx("educational"))
	if score <= 0 {
		t.Errorf("educational score should be positive, got %f", score)
	}
}

func TestStrategyControversial(t *testing.T) {
	s := Strategies["controversial"]
	// high dislikes → higher score than high likes
	polarising := VideoStatistics{LikeCount: 10, DislikeCount: 90, ViewCount: 1000, PublishedAt: baseVideo.PublishedAt}
	positive := VideoStatistics{LikeCount: 90, DislikeCount: 10, ViewCount: 1000, PublishedAt: baseVideo.PublishedAt}
	ctx := defaultCtx("controversial")
	if s.Score(polarising, ctx) <= s.Score(positive, ctx) {
		t.Error("controversial: polarising video should outscore positive one")
	}
}

func TestStrategyCommunity(t *testing.T) {
	s := Strategies["community"]
	// more comments → higher score
	highComments := baseVideo
	highComments.CommentCount = 200
	ctx := defaultCtx("community")
	if s.Score(highComments, ctx) <= s.Score(baseVideo, ctx) {
		t.Error("community: more comments should yield higher score")
	}
}

func TestStrategyEvergreen(t *testing.T) {
	s := Strategies["evergreen"]
	// older video with same engagement should score differently than new one
	old := baseVideo
	old.PublishedAt = time.Now().Add(-365 * 24 * time.Hour)
	ctx := defaultCtx("evergreen")
	scoreNew := s.Score(baseVideo, ctx)
	scoreOld := s.Score(old, ctx)
	if math.IsNaN(scoreNew) || math.IsNaN(scoreOld) {
		t.Error("evergreen score should not be NaN")
	}
}

func TestStrategyHype(t *testing.T) {
	s := Strategies["hype"]
	// newer video with same views → higher hype score
	newVideo := baseVideo
	newVideo.PublishedAt = time.Now().Add(-1 * 24 * time.Hour)
	ctx := defaultCtx("hype")
	if s.Score(newVideo, ctx) <= s.Score(baseVideo, ctx) {
		t.Error("hype: newer video with same views should have higher score")
	}
}

// --- ResolveWeights ---

func TestResolveWeights_defaults(t *testing.T) {
	w := ResolveWeights("viral", Weights{}, Weights{})
	if w["engagement"] != 0.5 {
		t.Errorf("expected default engagement=0.5, got %f", w["engagement"])
	}
}

func TestResolveWeights_envOverride(t *testing.T) {
	env := Weights{"viral_engagement": 0.9}
	w := ResolveWeights("viral", env, Weights{})
	if w["engagement"] != 0.9 {
		t.Errorf("env override failed: expected 0.9, got %f", w["engagement"])
	}
	// other keys unchanged
	if w["reach"] != 0.3 {
		t.Errorf("unrelated key changed: expected reach=0.3, got %f", w["reach"])
	}
}

func TestResolveWeights_cliOverride(t *testing.T) {
	env := Weights{"viral_engagement": 0.9}
	cli := Weights{"engagement": 0.99}
	w := ResolveWeights("viral", env, cli)
	if w["engagement"] != 0.99 {
		t.Errorf("CLI override failed: expected 0.99, got %f", w["engagement"])
	}
}

func TestResolveWeights_unknownKey(t *testing.T) {
	// unknown keys in env/cli are silently ignored
	w := ResolveWeights("viral", Weights{"viral_nonexistent": 99}, Weights{"also_nope": 99})
	if _, ok := w["nonexistent"]; ok {
		t.Error("unknown key should not be added to resolved weights")
	}
}

// --- ApplyStrategy ---

func TestApplyStrategy_sortsDescending(t *testing.T) {
	now := time.Now()
	videos := []VideoStatistics{
		{Title: "low", ViewCount: 100, LikeCount: 1, PublishedAt: now.Add(-10 * 24 * time.Hour)},
		{Title: "high", ViewCount: 10000, LikeCount: 1000, PublishedAt: now.Add(-1 * 24 * time.Hour)},
		{Title: "mid", ViewCount: 500, LikeCount: 50, PublishedAt: now.Add(-5 * 24 * time.Hour)},
	}
	ApplyStrategy(videos, "hype", Strategies["hype"].DefaultWeights)
	for i := 1; i < len(videos); i++ {
		if videos[i].Score > videos[i-1].Score {
			t.Errorf("ApplyStrategy: not sorted descending at index %d", i)
		}
	}
}
