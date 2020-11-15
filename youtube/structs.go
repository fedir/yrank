package youtube

import "time"

// Channel datastructure for JSON unmarshalling
type Channel struct {
	NextPageToken string `json:"nextPageToken"`
	PageInfo      struct {
		TotalResults int `json:"totalResults"`
	}
	Items []struct {
		PlaylistID string `json:"id"`
	} `json:"items"`
}

// Playlist datastructure for JSON unmarshalling
type Playlist struct {
	NextPageToken string `json:"nextPageToken"`
	PageInfo      struct {
		TotalResults int `json:"totalResults"`
	}
	Items []struct {
		Snippet struct {
			Title       string `json:"title"`
			Description string `json:"Description"`
		}
		ContentDetails struct {
			VideoID          string `json:"videoId"`
			VideoPublishedAt string `json:"videoPublishedAt"`
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
	Key                     string    `header:"Key"`
	Title                   string    `header:"Title"`
	URL                     string    `header:"URL"`
	PublishedAt             time.Time `header:"Published at"`
	ViewCount               int       `header:"View count"`
	LikeCount               int       `header:"Like count"`
	DislikeCount            int       `header:"Dislike count"`
	CommentCount            int       `header:"Comment count"`
	PositiveInterestingness float64   `header:"Positive interestingness"`
	GlobalBuzzIndex         int       `header:"Global buzz index"`
	TotalReaction           int       `header:"Total reaction"`
	TotalInterestingness    float64   `header:"Total reaction"`
}
