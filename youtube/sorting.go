package youtube

import "sort"

//SortBy sorts videos by some params
func SortBy(vs []VideoStatistics, sortingColumn string) {
	switch sortingColumn {
	case "total-interest":
		sort.Slice(vs[:], func(i, j int) bool {
			return vs[i].TotalReaction > vs[j].TotalReaction
		})
	case "positive-interest":
		sort.Slice(vs[:], func(i, j int) bool {
			return vs[i].PositiveInterestingness > vs[j].PositiveInterestingness
		})
	default:
		sort.Slice(vs[:], func(i, j int) bool {
			return vs[i].LikeCount > vs[j].LikeCount
		})
	}
}
