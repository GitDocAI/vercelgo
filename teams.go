package vercelgo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/YHVCorp/vercelgo/config"
	"github.com/YHVCorp/vercelgo/schemas"
	"github.com/YHVCorp/vercelgo/utils"
)

func (c *VercelClient) GetTeam(teamId string) (*schemas.Team, error) {
	response, statusCode, err := utils.DoReq[schemas.Team](fmt.Sprintf("%s/%s", config.TeamsURL, teamId), nil, "GET", c.GetHeaders(), false, time.Second*10)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get team: %s", response)
	}

	return &response, nil
}

func (c *VercelClient) ListTeams(limit int64) ([]schemas.Team, error) {
	response, statusCode, err := utils.DoReq[schemas.ListTeamsResponse](fmt.Sprintf("%s?limit=%d", config.TeamsURL, limit), nil, "GET", c.GetHeaders(), false, time.Second*10)
	if err != nil {
		return nil, fmt.Errorf("failed to list teams: %w", err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list teams: %v", err)
	}

	return response.Teams, nil
}

func (c *VercelClient) CreateTeam(slug, name string) (string, error) {
	team := schemas.Team{
		Name: name,
		Slug: slug,
	}

	body, err := json.Marshal(team)
	if err != nil {
		return "", fmt.Errorf("failed to marshal team: %w", err)
	}

	response, statusCode, err := utils.DoReq[schemas.Team](config.TeamsURL, body, "POST", c.GetHeaders(), false, time.Second*10)
	if err != nil {
		return "", fmt.Errorf("failed to create team: %w", err)
	}

	if statusCode != http.StatusOK {
		return "", fmt.Errorf("failed to create team: %s", response)
	}

	return team.ID, nil
}

func (c *VercelClient) DeleteTeam(teamId string) error {
	_, statusCode, err := utils.DoReq[map[string]interface{}](fmt.Sprintf("%s/%s", config.TeamsURL, teamId), nil, "DELETE", c.GetHeaders(), false, time.Second*10)
	if err != nil {
		return fmt.Errorf("failed to delete team: %v", err)
	}

	if statusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete team: %v", err)
	}

	return nil
}
