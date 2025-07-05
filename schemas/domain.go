package schemas

type Domain struct {
	Name       string `json:"name"`
	Method     string `json:"method,omitempty"`
	Verified   bool   `json:"verified,omitempty"`
	CDNEnabled bool   `json:"cdnEnabled,omitempty"`
}

type DomainInfo struct {
	Name                string               `json:"name"`
	ApexName            string               `json:"apexName"`
	ProjectID           string               `json:"projectId"`
	Redirect            string               `json:"redirect"`
	RedirectStatusCode  int                  `json:"redirectStatusCode"`
	GitBranch           string               `json:"gitBranch"`
	CustomEnvironmentID string               `json:"customEnvironmentId"`
	UpdatedAt           int64                `json:"updatedAt"`
	CreatedAt           int64                `json:"createdAt"`
	Verified            bool                 `json:"verified"`
	Verification        []DomainVerification `json:"verification"`
}

type DomainVerification struct {
	Type   string `json:"type"`
	Domain string `json:"domain"`
	Value  string `json:"value"`
	Reason string `json:"reason"`
}

type DomainConfigInfo struct {
	ConfiguredBy       *string                 `json:"configuredBy"` // puede ser null
	Nameservers        []string                `json:"nameservers"`
	ServiceType        string                  `json:"serviceType"`
	CNAMEs             []string                `json:"cnames"`
	AValues            []string                `json:"aValues"`
	Conflicts          []string                `json:"conflicts"`
	AcceptedChallenges []string                `json:"acceptedChallenges"`
	IPStatus           *string                 `json:"ipStatus"` // puede ser null
	Misconfigured      bool                    `json:"misconfigured"`
	RecommendedIPv4    []RecommendedIPv4Entry  `json:"recommendedIPv4"`
	RecommendedCNAME   []RecommendedCNAMEEntry `json:"recommendedCNAME"`
}

type RecommendedIPv4Entry struct {
	Rank  int      `json:"rank"`
	Value []string `json:"value"`
}

type RecommendedCNAMEEntry struct {
	Rank  int    `json:"rank"`
	Value string `json:"value"`
}

type DomainInfoWithVerification struct {
	Info   *DomainInfo       `json:"info"`
	Config *DomainConfigInfo `json:"config"`
}
