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

type Options struct {
	Production          *string `json:"production,omitempty"`
	Target              *string `json:"target,omitempty"`
	CustomEnvironmentID *string `json:"customEnvironmentId,omitempty"`
	GitBranch           *string `json:"gitBranch,omitempty"`
	Redirects           *string `json:"redirects,omitempty"`
	Redirect            *string `json:"redirect,omitempty"`
	Verified            *string `json:"verified,omitempty"`
	Limit               *int    `json:"limit,omitempty"`
	Since               *int64  `json:"since,omitempty"`
	Until               *int64  `json:"until,omitempty"`
	Order               *string `json:"order,omitempty"`
	Slug                *string `json:"slug,omitempty"`
}
