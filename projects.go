package vercelgo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GitDocAI/vercelgo/config"
	"github.com/GitDocAI/vercelgo/schemas"
	"github.com/GitDocAI/vercelgo/utils"
)

// Allows to create a new project with the provided configuration.
// It only requires the project name and team ID but more configuration can be provided to override the defaults.
func (c *VercelClient) CreateProject(payload schemas.CreateProjectRequest, teamId string) (*schemas.Project, error) {
	if payload.Name == "" {
		return nil, fmt.Errorf("project name is required")
	}
	if teamId == "" {
		return nil, fmt.Errorf("teamId is required")
	}
	if payload.Framework == "" {
		payload.Framework = "nextjs"
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create project request: %w", err)
	}

	url := fmt.Sprintf("%s/v11/projects?teamId=%s", config.BaseURL, teamId)

	response, status, err := utils.DoReq[schemas.Project](url, body, "POST", c.GetHeaders(), false, 15*time.Second)
	if err != nil {
		return nil, fmt.Errorf("create project error: %w", err)
	}

	if status != http.StatusOK && status != http.StatusCreated {
		return nil, fmt.Errorf("failed to create project: status %d", status)
	}

	return &response, nil
}

// Update the fields of a project using either its name or id
func (c *VercelClient) UpdateProject(projectIdOrName string, payload schemas.CreateProjectRequest, teamId string) (*schemas.Project, error) {
	if projectIdOrName == "" {
		return nil, fmt.Errorf("projectIdOrName is required")
	}
	if teamId == "" {
		return nil, fmt.Errorf("teamId is required")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update project request: %w", err)
	}

	url := fmt.Sprintf("%s/v9/projects/%s?teamId=%s", config.BaseURL, projectIdOrName, teamId)

	response, status, err := utils.DoReq[schemas.Project](url, body, "PATCH", c.GetHeaders(), false, 15*time.Second)
	if err != nil {
		return nil, fmt.Errorf("update project error: %w", err)
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("failed to update project: status %d", status)
	}

	return &response, nil
}

// Delete a specific project by passing either the project id or name
func (c *VercelClient) DeleteProject(projectIdOrName string, teamId string) error {
	if projectIdOrName == "" {
		return fmt.Errorf("projectIdOrName is required")
	}
	if teamId == "" {
		return fmt.Errorf("teamId is required")
	}

	url := fmt.Sprintf("%s/v9/projects/%s?teamId=%s", config.BaseURL, projectIdOrName, teamId)

	_, status, err := utils.DoReq[interface{}](url, nil, "DELETE", c.GetHeaders(), false, 15*time.Second)
	if err != nil {
		return fmt.Errorf("delete project error: %w", err)
	}

	if status != http.StatusNoContent {
		return fmt.Errorf("failed to delete project: status %d", status)
	}

	return nil
}
