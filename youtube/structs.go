package youtube

import "time"

// ChannelByHandle is the API response when resolving a handle to a channel ID.
type ChannelByHandle struct {
	Items []struct {
		ID string `json:"id"`
	} `json:"items"`
}

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
			Description string `json:"description"`
		}
		ContentDetails struct {
			VideoID          string `json:"videoId"`
			VideoPublishedAt string `json:"videoPublishedAt"`
		} `json:"contentDetails"`
	} `json:"items"`
}

// Search is the API response from the search.list endpoint (type=video).
type Search struct {
	NextPageToken string `json:"nextPageToken"`
	Items         []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			Title       string `json:"title"`
			PublishedAt string `json:"publishedAt"`
		} `json:"snippet"`
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
		ContentDetails struct {
			Duration string `json:"duration"`
		} `json:"contentDetails"`
	} `json:"items"`
}

// VideoStatistics statistics of a singular playlist
type VideoStatistics struct {
	Score                       float64            `header:"Score"`
	AllScores                   map[string]float64 `header:"-"`
	Key                         string    `header:"Key"`
	Title                       string    `header:"Title"`
	URL                         string    `header:"URL"`
	PublishedAt                 time.Time `header:"Published at"`
	Duration                    int       `header:"Duration"`
	ViewCount                   int       `header:"View count"`
	LikeCount                   int       `header:"Like count"`
	DislikeCount                int       `header:"Dislike count"`
	CommentCount                int       `header:"Comment count"`
	PositiveInterestingness     float64   `header:"Positive interestingness"`
	PositiveNegativeCoefficient float64   `header:"Positive/negative coefficient"`
	GlobalBuzzIndex             int       `header:"Global buzz index"`
	TotalReaction               int       `header:"Total reaction"`
	TotalInterestingness        float64   `header:"Total interestingness"`
}
