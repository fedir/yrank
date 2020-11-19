package youtube

import "sort"

//SortBy sorts videos by some params
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
