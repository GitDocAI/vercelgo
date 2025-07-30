package schemas

type DeploymentFile struct {
	File string `json:"file"`
	Sha  string `json:"sha"`
}

type CreateDeploymentRequest struct {
	Name    string           `json:"name"`
	Project string           `json:"project"`
	Files   []DeploymentFile `json:"files"`
}

type DeploymentResponse struct {
	Id         string           `json:"id"`
	Url        string           `json:"url"`
	Files      []DeploymentFile `json:"files"`
	Status     string           `json:"status"`
	ReadyState string           `json:"readyState"`
	CreatedAt  int64            `json:"createdAt"`
}

type DeploymentStatus struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Url        string `json:"url"`
	Status     string `json:"status"`
	ReadyState string `json:"readyState"`
	CreatedAt  int64  `json:"createdAt"`
}
