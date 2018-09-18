package youtube

// Playlist datastructure for JSON unmarshalling
type Playlist struct {
	PageInfo struct {
		TotalResults int `json:"totalResults"`
	}
	Items []struct {
		Snippet struct {
			Title       string `json:"title"`
			Description string `json:"Description"`
		}
		ContentDetails struct {
			VideoID string `json:"videoId"`
		} `json:"contentDetails"`
	} `json:"items"`
}

// Video datastructure for JSON unmarshalling and future ranking
type Video struct {
	Items []struct {
		ID         string `json:"id"`
		Statistics struct {
			ViewCount    string `json:"viewCount"`
			LikeCount    string `json:"likeCount"`
			DislikeCount string `json:"dislikeCount"`
			CommentCount string `json:"commentCount"`
		} `json:"statistics"`
	} `json:"items"`
}

// VideoStatistics statistics of a singular playlist
type VideoStatistics struct {
	Key          string `header:"Key"`
	Title        string `header:"Title"`
	URL          string `header:"URL"`
	ViewCount    string `header:"View count"`
	LikeCount    string `header:"Like count"`
	DislikeCount string `header:"Dislike count"`
	CommentCount string `header:"Comment count"`
}
