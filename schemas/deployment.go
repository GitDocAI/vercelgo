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
	Id    string           `json:"id"`
	Url   string           `json:"url"`
	Files []DeploymentFile `json:"files"`
}
