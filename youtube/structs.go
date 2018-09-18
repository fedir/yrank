package youtube

// Playlist datastructure for JSON unmarshalling
type Playlist struct {
	PageInfo struct {
		TotalResults int `json:"totalResults"`
	}
	Items []struct {
		ContentDetails struct {
			VideoID string `json:"videoId"`
		} `json:"contentDetails"`
	} `json:"items"`
}

// Video datastructure for JSON unmarshalling and future ranking
type Video struct {
	Items []struct {
		Statistics []struct {
			ViewCount    int `json:"viewCount"`
			LikeCount    int `json:"likeCount"`
			DislikeCount int `json:"dislikeCount"`
		} `json:"statistics"`
	} `json:"items"`
}
