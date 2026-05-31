package youtube

import (
	"math"
	"strings"
	"time"
)

// Weights maps a weight key to its multiplier.
type Weights map[string]float64

// StrategyContext holds dataset-level values needed for normalisation.
type StrategyContext struct {
	MaxViews int
	Now      time.Time
	Weights  Weights
}

// Strategy defines a named scoring approach.
type Strategy struct {
	Name           string
	Description    string
	DefaultWeights Weights
	Score          func(v VideoStatistics, ctx StrategyContext) float64
}

// Strategies is the registry of all available evaluation strategies.
var Strategies = map[string]Strategy{
	"viral": {
		Name:        "Viral",
		Description: "Rewards high engagement rate on a large audience base (algo/trending lens)",
		DefaultWeights: Weights{
			"engagement": 0.5,
			"reach":      0.3,
			"comments":   0.2,
		},
		Score: func(v VideoStatistics, ctx StrategyContext) float64 {
			if v.ViewCount == 0 {
				return 0
			}
			engagement := float64(v.LikeCount+v.DislikeCount) / float64(v.ViewCount)
			reach := 0.0
			if ctx.MaxViews > 0 {
				reach = float64(v.ViewCount) / float64(ctx.MaxViews)
			}
			comments := float64(v.CommentCount) / float64(v.ViewCount)
			return ctx.Weights["engagement"]*engagement +
				ctx.Weights["reach"]*reach +
				ctx.Weights["comments"]*comments
		},
	},

	"educational": {
		Name:        "Educational",
		Description: "Rewards likes and discussion; deprioritises recency (tutorial/reference lens)",
		DefaultWeights: Weights{
			"likes":    0.6,
			"comments": 0.3,
			"recency":  0.1,
		},
		Score: func(v VideoStatistics, ctx StrategyContext) float64 {
			if v.ViewCount == 0 {
				return 0
			}
			likes := float64(v.LikeCount) / float64(v.ViewCount)
			comments := float64(v.CommentCount) / float64(v.ViewCount)
			ageDays := ctx.Now.Sub(v.PublishedAt).Hours() / 24
			if ageDays < 1 {
				ageDays = 1
			}
			recency := 1.0 / ageDays
			return ctx.Weights["likes"]*likes +
				ctx.Weights["comments"]*comments +
				ctx.Weights["recency"]*recency
		},
	},

	"controversial": {
		Name:        "Controversial",
		Description: "High dislike ratio on a large reaction base (debate/polarising lens)",
		DefaultWeights: Weights{
			"ratio":  1.0,
			"volume": 1.0,
		},
		Score: func(v VideoStatistics, ctx StrategyContext) float64 {
			reactions := float64(v.LikeCount + v.DislikeCount)
			ratio := float64(v.DislikeCount+1) / float64(v.LikeCount+1)
			volume := math.Log1p(reactions)
			return ctx.Weights["ratio"]*ratio * ctx.Weights["volume"]*volume
		},
	},

	"community": {
		Name:        "Community",
		Description: "Comments-first; sentiment secondary (fan/community-building lens)",
		DefaultWeights: Weights{
			"comments":  0.5,
			"sentiment": 0.5,
		},
		Score: func(v VideoStatistics, ctx StrategyContext) float64 {
			if v.ViewCount == 0 {
				return 0
			}
			comments := float64(v.CommentCount) / float64(v.ViewCount)
			sentiment := float64(v.LikeCount) / float64(1+v.DislikeCount)
			normSentiment := sentiment / (1 + sentiment)
			return ctx.Weights["comments"]*comments +
				ctx.Weights["sentiment"]*normSentiment
		},
	},

	"evergreen": {
		Name:        "Evergreen",
		Description: "Steady engagement per day of life (long-tail/SEO lens)",
		DefaultWeights: Weights{
			"engagement": 0.5,
			"age":        0.5,
		},
		Score: func(v VideoStatistics, ctx StrategyContext) float64 {
			ageDays := ctx.Now.Sub(v.PublishedAt).Hours() / 24
			if ageDays < 1 {
				ageDays = 1
			}
			engagement := float64(v.LikeCount+v.CommentCount) / ageDays
			ageScore := 1.0 / ageDays
			return ctx.Weights["engagement"]*engagement +
				ctx.Weights["age"]*ageScore
		},
	},

	"hype": {
		Name:        "Hype",
		Description: "Pure view velocity — views per day since publication (launch/premiere lens)",
		DefaultWeights: Weights{
			"velocity": 1.0,
		},
		Score: func(v VideoStatistics, ctx StrategyContext) float64 {
			ageDays := ctx.Now.Sub(v.PublishedAt).Hours() / 24
			if ageDays < 1 {
				ageDays = 1
			}
			return ctx.Weights["velocity"] * float64(v.ViewCount) / ageDays
		},
	},
}

// ResolveWeights merges weights in priority order:
// strategy defaults → env overrides → CLI overrides.
func ResolveWeights(slug string, env Weights, cli Weights) Weights {
	s, ok := Strategies[slug]
	if !ok {
		return Weights{}
	}
	resolved := make(Weights, len(s.DefaultWeights))
	for k, v := range s.DefaultWeights {
		resolved[k] = v
	}
	prefix := strings.ToLower(slug) + "_"
	for k, v := range env {
		key := strings.TrimPrefix(strings.ToLower(k), prefix)
		if _, exists := resolved[key]; exists {
			resolved[key] = v
		}
	}
	for k, v := range cli {
		if _, exists := resolved[k]; exists {
			resolved[k] = v
		}
	}
	return resolved
}
