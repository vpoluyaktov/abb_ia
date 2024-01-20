package github

import (
	"fmt"

	"github.com/hashicorp/go-version"
)

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
