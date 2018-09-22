package youtube

import "sort"

func SortByLikes(vs []VideoStatistics) {
	sort.Slice(vs[:], func(i, j int) bool {
		return vs[i].LikeCount > vs[j].LikeCount
	})
}

func SortByPositiveInterestingness(vs []VideoStatistics) {
	sort.Slice(vs[:], func(i, j int) bool {
		return vs[i].PositiveInterestingness > vs[j].PositiveInterestingness
	})
}

func SortByTotalInterestingness(vs []VideoStatistics) {
	sort.Slice(vs[:], func(i, j int) bool {
		return vs[i].TotalInterestingness > vs[j].TotalInterestingness
	})
}
