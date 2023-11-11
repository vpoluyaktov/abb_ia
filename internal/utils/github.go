package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-version"
)

type Release struct {
	TagName string `json:"tag_name"`
}

func GetLatestVersion(owner string, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

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

func CompareVersions(version1 string, version2 string) (int, error) {
	v1, err := version.NewVersion(version1)
	if err != nil {
		return 0, fmt.Errorf("invalid version format: %s", version1)
	}

	v2, err := version.NewVersion(version2)
	if err != nil {
		return 0, fmt.Errorf("invalid version format: %s", version2)
	}

	return v1.Compare(v2), nil
}
