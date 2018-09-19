package youtube

import "sort"

func sortByLikes(vs []VideoStatistics) {
	sort.Slice(vs[:], func(i, j int) bool {
		return vs[i].LikeCount > vs[j].LikeCount
	})
}
