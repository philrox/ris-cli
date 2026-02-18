package model

// SearchResult represents paginated search results from the RIS API.
type SearchResult struct {
	TotalHits int        `json:"total_hits"`
	Page      int        `json:"page"`
	PageSize  int        `json:"page_size"`
	HasMore   bool       `json:"has_more"`
	Documents []Document `json:"documents"`
}
