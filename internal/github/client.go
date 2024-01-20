package github

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GithubClient struct {
	repoOwner string
	repoName  string
}

func NewClient(repoOwner string, repoName string) (*GithubClient) {
	c := &GithubClient{repoOwner: repoOwner, repoName: repoName}
	return c
}

type Release struct {
	TagName string `json:"tag_name"`
}

func (c *GithubClient) GetLatestVersion() (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", c.repoOwner, c.repoName)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch latest release: %s", resp.Status)
	}

	var release Release
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return "", err
	}

	return release.TagName, nil
}

