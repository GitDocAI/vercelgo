package schemas

type ListProjectsResponse struct {
	Projects   []Project  `json:"projects"`
	Pagination Pagination `json:"pagination"`
}

type Project struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AccountID string `json:"accountId"`
}

type CreateProjectRequest struct {
	Name            string `json:"name"`
	BuildCommand    string `json:"buildCommand,omitempty"`
	InstallCommand  string `json:"installCommand,omitempty"`
	DevCommand      string `json:"devCommand,omitempty"`
	Framework       string `json:"framework,omitempty"`
	OutputDirectory string `json:"outputDirectory,omitempty"`
	PublicSource    *bool  `json:"publicSource,omitempty"`
}
