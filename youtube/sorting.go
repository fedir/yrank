package youtube

import "sort"

// SortByLikes sorts videos by it's likes
func SortByLikes(vs []VideoStatistics) {
	sort.Slice(vs[:], func(i, j int) bool {
		return vs[i].LikeCount > vs[j].LikeCount
	})
}

// SortByPositiveInterestingness sorts videos by it's positive (ignoring negative) interestingness
func SortByPositiveInterestingness(vs []VideoStatistics) {
	sort.Slice(vs[:], func(i, j int) bool {
		return vs[i].PositiveInterestingness > vs[j].PositiveInterestingness
	})
}

// SortByTotalInterestingness sorts videos by it's total interestingness
func SortByTotalInterestingness(vs []VideoStatistics) {
	sort.Slice(vs[:], func(i, j int) bool {
		return vs[i].TotalInterestingness > vs[j].TotalInterestingness
	})
}
