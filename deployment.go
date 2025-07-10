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
func (c *VercelClient) Deploy(projectId, deploymentName, directory, teamId string) (*schemas.ProjectDomainsResponse, error) {
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

	projectDomains, err := c.GetProjectDomains(projectId, teamId, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get project domains: %w", err)
	}

	return projectDomains, nil
}
