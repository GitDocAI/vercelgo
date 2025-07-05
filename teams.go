package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/YHVCorp/vercelgo/config"
	"github.com/YHVCorp/vercelgo/schemas"
	"github.com/YHVCorp/vercelgo/utils"
)

// Get information for the Team specified by the teamId parameter.
func (c *VercelClient) GetTeam(teamId string) (*schemas.Team, error) {
	response, statusCode, err := utils.DoReq[schemas.Team](
		fmt.Sprintf("%s/v1/teams/%s", config.BaseURL, teamId),
		nil, "GET", c.GetHeaders(), false, time.Second*10,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get team: %s", response)
	}

	return &response, nil
}

// Get a paginated list of all the Teams the authenticated User is a member of.
func (c *VercelClient) ListTeams(filter *schemas.Filter) ([]schemas.Team, error) {
	url := fmt.Sprintf("%s/v1/teams", config.BaseURL)
	if filter != nil && filter.Limit > 0 {
		url = fmt.Sprintf("%s?limit=%d", url, filter.Limit)
	}

	response, statusCode, err := utils.DoReq[schemas.ListTeamsResponse](url, nil, "GET", c.GetHeaders(), false, time.Second*10)
	if err != nil {
		return nil, fmt.Errorf("failed to list teams: %w", err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list teams: %v", err)
	}

	return response.Teams, nil
}

// Create a new Team under your account.
// You need to send a POST request with the desired Team slug, and optionally the Team name.
func (c *VercelClient) CreateTeam(slug, name string) (string, error) {
	team := schemas.Team{
		Name: name,
		Slug: slug,
	}

	body, err := json.Marshal(team)
	if err != nil {
		return "", fmt.Errorf("failed to marshal team: %w", err)
	}

	response, statusCode, err := utils.DoReq[schemas.Team](fmt.Sprintf("%s/v1/teams", config.BaseURL), body, "POST", c.GetHeaders(), false, time.Second*10)
	if err != nil {
		return "", fmt.Errorf("failed to create team: %w", err)
	}

	if statusCode != http.StatusOK {
		return "", fmt.Errorf("failed to create team: %s", response)
	}

	return team.ID, nil
}

// Update the information of a Team specified by the teamId parameter.
// The name and slug parameters are optional.
func (c *VercelClient) UpdateTeam(teamId, name, slug string) (*schemas.Team, error) {
	team := schemas.Team{}
	if name != "" {
		team.Name = name
	}
	if slug != "" {
		team.Slug = slug
	}

	body, err := json.Marshal(team)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal team: %w", err)
	}

	response, statusCode, err := utils.DoReq[schemas.Team](
		fmt.Sprintf("%s/v1/teams/%s", config.BaseURL, teamId),
		body, "PATCH", c.GetHeaders(), false, time.Second*10,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update team: %w", err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update team: %s", response)
	}

	return &response, nil
}

// Delete a team under your account.
// You need to send teamId as required parameter.
// An optional array of reasons for deletion may also be sent.
func (c *VercelClient) DeleteTeam(teamId string, reasons []schemas.Reason) error {
	body, err := json.Marshal(schemas.DeleteTeamRequest{Reasons: reasons})
	if err != nil {
		return fmt.Errorf("failed to marshal delete team request: %w", err)
	}
	_, statusCode, err := utils.DoReq[map[string]interface{}](
		fmt.Sprintf("%s/v1/teams/%s", config.BaseURL, teamId),
		body, "DELETE", c.GetHeaders(), false, time.Second*10,
	)
	if err != nil {
		return fmt.Errorf("failed to delete team: %v", err)
	}

	if statusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete team: %v", err)
	}

	return nil
}
