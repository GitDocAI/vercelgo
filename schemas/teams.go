package schemas

type Team struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Slug string `json:"slug,omitempty"`
}

type ListTeamsResponse struct {
	Teams      []Team     `json:"teams,omitempty"`
	Pagination Pagination `json:"pagination,omitempty"`
}

type DeleteTeamRequest struct {
	Reasons []Reason `json:"reasons,omitempty"`
}

type Reason struct {
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
}
