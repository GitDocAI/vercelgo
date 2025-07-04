package schemas

type Pagination struct {
	Count    int64       `json:"count"`
	Next     interface{} `json:"next,omitempty"`
	Previous interface{} `json:"previous,omitempty"`
}
