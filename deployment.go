package vercelgo

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/YHVCorp/vercelgo/config"
	"github.com/YHVCorp/vercelgo/schemas"
	"github.com/YHVCorp/vercelgo/utils"
)

// Deploy uploads files to Vercel from a directory and creates a deployment for the specified project.
func (c *VercelClient) Deploy(projectId, deploymentName, directory, teamId string) (*schemas.AllDomainWithVerification, error) {
	files := []schemas.DeploymentFile{}
	err := filepath.WalkDir(directory, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %q: %v", path, err)
		}
		name := d.Name()
		if d.IsDir() {
			if name == "node_modules" || name == ".git" || name == ".next" || strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			return nil
		} else {
			if strings.HasPrefix(name, ".") {
				return nil
			}
		}

		relPath, err := filepath.Rel(directory, path)
		if err != nil {
			return fmt.Errorf("error getting relative path for %q: %v", path, err)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file %q: %v", path, err)
		}

		hashBytes := sha1.Sum(content)
		hash := hex.EncodeToString(hashBytes[:])

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/files?teamId=%s", config.BaseURL, teamId), bytes.NewReader(content))
		if err != nil {
			return fmt.Errorf("error creating request for %q: %v", path, err)
		}
		req.Header.Set("Authorization", "Bearer "+c.Token)
		req.Header.Set("x-vercel-digest", hash)
		req.Header.Set("Content-Length", fmt.Sprintf("%d", len(content)))
		req.Header.Set("Content-Type", "application/octet-stream")

		client := &http.Client{Timeout: 15 * time.Second}
		res, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error uploading file %q: %v", path, err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(res.Body)
			return fmt.Errorf("upload failed (%d): %s", res.StatusCode, string(body))
		}

		files = append(files, schemas.DeploymentFile{
			File: relPath,
			Sha:  hash,
		})
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed walking files: %w", err)
	}

	deploymentReq := schemas.CreateDeploymentRequest{
		Name:    deploymentName,
		Project: projectId,
		Files:   files,
	}

	body, err := json.Marshal(deploymentReq)
	if err != nil {
		return nil, fmt.Errorf("marshal deployment error: %w", err)
	}

	_, status, err := utils.DoReq[schemas.DeploymentResponse](fmt.Sprintf("%s/v13/deployments?teamId=%s", config.BaseURL, teamId), body, "POST", c.GetHeaders(), false, 30*time.Second)
	if err != nil {
		return nil, fmt.Errorf("create deployment error: %w", err)
	}
	if status != http.StatusOK && status != http.StatusCreated {
		return nil, fmt.Errorf("deployment failed with status %d", status)
	}

	allDomains, err := c.GetProjectDomains(projectId, teamId, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get project domains: %w", err)
	}

	return allDomains, nil
}

// GetDeploymentStatus gets the status of a specific deployment by its ID and team ID.
func (c *VercelClient) GetDeploymentStatus(deploymentId, teamId string) (*schemas.DeploymentStatus, error) {
	deploymentStatus, status, err := utils.DoReq[schemas.DeploymentStatus](
		fmt.Sprintf("%s/v13/deployments/%s?teamId=%s", config.BaseURL, deploymentId, teamId),
		nil,
		"GET",
		c.GetHeaders(),
		false,
		30*time.Second,
	)
	if err != nil {
		return nil, fmt.Errorf("get deployment status error: %w", err)
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("failed to get deployment status with code %d", status)
	}

	return &deploymentStatus, nil
}

func (c *VercelClient) WaitForDeployment(deploymentId, teamId string, timeout time.Duration) (*schemas.DeploymentStatus, error) {
	if timeout == 0 {
		timeout = 10 * time.Minute
	}

	startTime := time.Now()
	checkInterval := 5 * time.Second

	for {
		if time.Since(startTime) > timeout {
			return nil, fmt.Errorf("deployment monitoring timed out after %v", timeout)
		}

		status, err := c.GetDeploymentStatus(deploymentId, teamId)
		if err != nil {
			return nil, fmt.Errorf("error checking deployment status: %w", err)
		}

		switch status.ReadyState {
		case "READY":
			return status, nil
		case "ERROR":
			return status, fmt.Errorf("deployment failed with state: %s", status.ReadyState)
		case "CANCELED":
			return status, fmt.Errorf("deployment was canceled")
		}

		time.Sleep(checkInterval)
	}
}
