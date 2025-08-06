package vercelgo

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/GitDocAI/vercelgo/config"
	"github.com/GitDocAI/vercelgo/schemas"
	"github.com/GitDocAI/vercelgo/utils"
)

// AddProjectDomain adds a new domain to a specific project.
// It requires the domain name, team ID, project ID or name
// and returns the domain information along with its configuration.
func (c *VercelClient) AddProjectDomain(domainName, teamId, projectIdOrName string) (*schemas.AllDomainWithVerification, error) {
	reqBody := schemas.Domain{
		Name: domainName,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal domain request: %w", err)
	}

	url := fmt.Sprintf("%s/v10/projects/%s/domains?teamId=%s", config.BaseURL, projectIdOrName, teamId)
	_, status, err := utils.DoReq[schemas.DomainInfo](url, bodyBytes, "POST", c.GetHeaders(), false, 15*time.Second)
	if err != nil {
		return nil, fmt.Errorf("error adding domain: %w", err)
	}
	if status != 200 && status != 201 {
		return nil, fmt.Errorf("unexpected status code: %d", status)
	}

	return c.GetProjectDomains(projectIdOrName, teamId, nil)
}

// DeleteProjectDomain removes a domain from a specific project by its name or ID.
// It requires the domain name, project ID or name, and team ID.
func (c *VercelClient) DeleteProjectDomain(domainName, projectIdOrName, teamId string) (*schemas.AllDomainWithVerification, error) {
	if domainName == "" || projectIdOrName == "" || teamId == "" {
		return nil, fmt.Errorf("domainName, projectIdOrName, and teamId are required")
	}

	url := fmt.Sprintf("%s/v9/projects/%s/domains/%s?teamId=%s", config.BaseURL, projectIdOrName, domainName, teamId)

	_, status, err := utils.DoReq[interface{}](url, nil, "DELETE", c.GetHeaders(), false, 15*time.Second)
	if err != nil {
		return nil, fmt.Errorf("error deleting domain: %w", err)
	}
	if status != 200 && status != 204 {
		return nil, fmt.Errorf("unexpected status code: %d", status)
	}

	return c.GetProjectDomains(projectIdOrName, teamId, nil)
}

// GetProjectDomains retrieves all domains associated with a project by ID or name.
// It supports various filtering options including production, target environment, git branch, etc.
func (c *VercelClient) GetProjectDomains(projectIdOrName, teamId string, opts *schemas.Options) (*schemas.AllDomainWithVerification, error) {
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

	domainsWithVerification := make([]schemas.DomainInfoWithVerification, len(response.Domains))
	for i, domain := range response.Domains {
		config, err := c.GetDomainConfig(domain.Name, teamId)
		if err != nil {
			return nil, fmt.Errorf("error getting config for domain %s: %w", domain.Name, err)
		}
		domainsWithVerification[i] = schemas.DomainInfoWithVerification{
			Info:   &domain,
			Config: config,
		}
	}
	return &schemas.AllDomainWithVerification{
		Domains: domainsWithVerification,
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
