package vercelgo

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/YHVCorp/vercelgo/config"
	"github.com/YHVCorp/vercelgo/schemas"
	"github.com/YHVCorp/vercelgo/utils"
)

// AddProjectDomain adds a new domain to a specific project.
// It requires the domain name, team ID, project ID or name
// and returns the domain information along with its configuration.
func (c *VercelClient) AddProjectDomain(domainName, teamId, projectIdOrName string) (*schemas.DomainInfoWithVerification, error) {
	reqBody := schemas.Domain{
		Name: domainName,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal domain request: %w", err)
	}

	url := fmt.Sprintf("%s/v10/projects/%s/domains?teamId=%s", config.BaseURL, projectIdOrName, teamId)
	domainInfo, status, err := utils.DoReq[schemas.DomainInfo](url, bodyBytes, "POST", c.GetHeaders(), false, 15*time.Second)
	if err != nil {
		return nil, fmt.Errorf("error adding domain: %w", err)
	}
	if status != 200 && status != 201 {
		return nil, fmt.Errorf("unexpected status code: %d", status)
	}

	domainConfigInfo, err := c.GetDomainConfig(domainName, teamId)
	if err != nil {
		return nil, fmt.Errorf("error getting domain config: %w", err)
	}

	return &schemas.DomainInfoWithVerification{
		Info:   &domainInfo,
		Config: domainConfigInfo,
	}, nil
}

// GetDomainConfig retrieves the configuration details of a domain by its name and team ID.
// It returns a DomainConfigInfo struct containing the configuration details.
func (c *VercelClient) GetDomainConfig(domainName, teamId string) (*schemas.DomainConfigInfo, error) {
	url := fmt.Sprintf("%s/v6/domains/%s/config?teamId=%s", config.BaseURL, domainName, teamId)

	response, status, err := utils.DoReq[schemas.DomainConfigInfo](url, nil, "GET", c.GetHeaders(), false, 15*time.Second)
	if err != nil {
		return nil, fmt.Errorf("error getting domain config: %w", err)
	}
	if status != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", status)
	}

	return &response, nil
}

// VerifyDomain retrieves the domain information and its configuration for a specific project.
func (c *VercelClient) VerifyDomain(domainName, projectIdOrName, teamId string) (*schemas.DomainInfoWithVerification, error) {
	url := fmt.Sprintf("%s/v9/projects/%s/domains/%s?teamId=%s", config.BaseURL, projectIdOrName, domainName, teamId)

	info, status, err := utils.DoReq[schemas.DomainInfo](url, nil, "GET", c.GetHeaders(), false, 15*time.Second)
	if err != nil {
		return nil, fmt.Errorf("error getting domain info: %w", err)
	}
	if status != 200 && status != 201 {
		return nil, fmt.Errorf("unexpected status code: %d", status)
	}

	configInfo, err := c.GetDomainConfig(domainName, teamId)
	if err != nil {
		return nil, fmt.Errorf("error getting domain config: %w", err)
	}

	return &schemas.DomainInfoWithVerification{
		Info:   &info,
		Config: configInfo,
	}, nil
}

// DeleteProjectDomain removes a domain from a specific project by its name or ID.
// It requires the domain name, project ID or name, and team ID.
func (c *VercelClient) DeleteProjectDomain(domainName, projectIdOrName, teamId string) error {
	if domainName == "" || projectIdOrName == "" || teamId == "" {
		return fmt.Errorf("domainName, projectIdOrName, and teamId are required")
	}

	url := fmt.Sprintf("%s/v9/projects/%s/domains/%s?teamId=%s", config.BaseURL, projectIdOrName, domainName, teamId)

	_, status, err := utils.DoReq[interface{}](url, nil, "DELETE", c.GetHeaders(), false, 15*time.Second)
	if err != nil {
		return fmt.Errorf("error deleting domain: %w", err)
	}
	if status != 200 && status != 204 {
		return fmt.Errorf("unexpected status code: %d", status)
	}

	return nil
}

// GetProjectDomains retrieves all domains associated with a project by ID or name.
// It supports various filtering options including production, target environment, git branch, etc.
func (c *VercelClient) GetProjectDomains(projectIdOrName, teamId string, opts *schemas.Options) (*schemas.ProjectDomainsResponse, error) {
	if projectIdOrName == "" || teamId == "" {
		return nil, fmt.Errorf("projectIdOrName and teamId are required")
	}

	url := fmt.Sprintf("%s/v9/projects/%s/domains", config.BaseURL, projectIdOrName)
	params := utils.BuildProjectDomainsParams(teamId, opts)
	if params != "" {
		url += "?" + params
	}

	response, status, err := utils.DoReq[schemas.ProjectDomainsResponse](url, nil, "GET", c.GetHeaders(), false, 15*time.Second)
	if err != nil {
		return nil, fmt.Errorf("error getting project domains: %w", err)
	}
	if status != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", status)
	}

	return &response, nil
}
