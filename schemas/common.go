package schemas

type Pagination struct {
	Count int   `json:"count"`
	Next  int64 `json:"next"`
	Prev  int64 `json:"prev"`
}

type Filter struct {
	Limit  int64  `json:"limit,omitempty"`
	Search string `json:"search,omitempty"`
	TeamID string `json:"teamId,omitempty"`
	Slug   string `json:"slug,omitempty"`
}
