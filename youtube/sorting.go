package youtube

import (
	"log"
	"sort"
	"time"
)

// ApplyStrategy scores every video using the named strategy and resolved
// weights, then sorts the slice descending by Score.
func ApplyStrategy(vs []VideoStatistics, slug string, weights Weights) {
	s, ok := Strategies[slug]
	if !ok {
		log.Fatalf("unknown strategy %q", slug)
	}
	maxViews := 0
	for _, v := range vs {
		if v.ViewCount > maxViews {
			maxViews = v.ViewCount
		}
	}
	ctx := StrategyContext{
		MaxViews: maxViews,
		Now:      time.Now(),
		Weights:  weights,
	}
	for i := range vs {
		vs[i].Score = s.Score(vs[i], ctx)
	}
	sort.Slice(vs, func(i, j int) bool {
		return vs[i].Score > vs[j].Score
	})
}

// StrategyOrder defines the canonical column order for -strategy all.
var StrategyOrder = []string{"viral", "educational", "controversial", "community", "evergreen", "hype"}

// ApplyAllStrategies scores every video with every strategy using default
// weights, stores results in AllScores, then sorts by total interestingness.
func ApplyAllStrategies(vs []VideoStatistics) {
	maxViews := 0
	for _, v := range vs {
		if v.ViewCount > maxViews {
			maxViews = v.ViewCount
		}
	}
	now := time.Now()
	for i := range vs {
		vs[i].AllScores = make(map[string]float64, len(StrategyOrder))
		for _, slug := range StrategyOrder {
			s := Strategies[slug]
			ctx := StrategyContext{MaxViews: maxViews, Now: now, Weights: s.DefaultWeights}
			vs[i].AllScores[slug] = s.Score(vs[i], ctx)
		}
	}
	SortBy(vs, "total-interest")
}

// SortBy sorts videos by some params
func SortBy(vs []VideoStatistics, sortingColumn string) {
	switch sortingColumn {
	case "total-reaction":
		sort.Slice(vs[:], func(i, j int) bool {
			return vs[i].TotalReaction > vs[j].TotalReaction
		})
	case "positive-interest":
		sort.Slice(vs[:], func(i, j int) bool {
			return vs[i].PositiveInterestingness > vs[j].PositiveInterestingness
		})
	case "pnc":
		fallthrough
	case "positive-negative-coefficient":
		sort.Slice(vs[:], func(i, j int) bool {
			return vs[i].PositiveNegativeCoefficient > vs[j].PositiveNegativeCoefficient
		})
	case "global-buzz-index":
		sort.Slice(vs[:], func(i, j int) bool {
			return vs[i].GlobalBuzzIndex > vs[j].GlobalBuzzIndex
		})
	case "likes":
		sort.Slice(vs[:], func(i, j int) bool {
			return vs[i].LikeCount > vs[j].LikeCount
		})
	case "total-interest":
		fallthrough
	default:
		sort.Slice(vs[:], func(i, j int) bool {
			return vs[i].TotalInterestingness > vs[j].TotalInterestingness
		})
	}
}
