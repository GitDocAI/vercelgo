package schemas

type DeploymentFile struct {
	File string `json:"file"`
	Sha  string `json:"sha"`
}

type CreateDeploymentRequest struct {
	Name    string           `json:"name"`
	Project string           `json:"project"`
	Files   []DeploymentFile `json:"files"`
	Target  string           `json:"target"`
}

type DeploymentResponse struct {
	Id         string           `json:"id"`
	Uid        string           `json:"uid"`
	Url        string           `json:"url"`
	Files      []DeploymentFile `json:"files"`
	Status     string           `json:"status"`
	ReadyState string           `json:"readyState"`
	Target     string           `json:"target"`
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

type DeploymentListResponse struct {
	Pagination  Pagination           `json:"pagination"`
	Deployments []DeploymentResponse `json:"deployments"`
}

type CurrentDeploymentResponse struct {
	Deployment CurrentDeployment `json:"deployment"`
	Domain     CurrentDomain     `json:"domain"`
}

type CurrentDeployment struct {
	CreatedAt              int64   `json:"createdAt"`
	DeploymentHostname     string  `json:"deploymentHostname"`
	Id                     string  `json:"id"`
	Name                   string  `json:"name"`
	ReadyState             string  `json:"readyState"`
	ReadySubstate          string  `json:"readySubstate"`
	Source                 string  `json:"source"`
	TeamId                 string  `json:"teamId"`
	Url                    string  `json:"url"`
	UserId                 string  `json:"userId"`
	ProjectId              string  `json:"projectId"`
	Target                 string  `json:"target"`
	AliasError             *string `json:"aliasError"`
	AliasAssignedAt        int64   `json:"aliasAssignedAt"`
	AliasAssigned          int64   `json:"aliasAssigned"`
	ReadyStateAt           int64   `json:"readyStateAt"`
	BuildingAt             int64   `json:"buildingAt"`
	PreviewCommentsEnabled bool    `json:"previewCommentsEnabled"`
}

type CurrentDomain struct {
	Name      string `json:"name"`
	ApexName  string `json:"apexName"`
	ProjectId string `json:"projectId"`
	UpdatedAt int64  `json:"updatedAt"`
	CreatedAt int64  `json:"createdAt"`
	Verified  bool   `json:"verified"`
}
