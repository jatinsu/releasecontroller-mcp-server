package utils

import (
	"encoding/json"
	"fmt"

	configv1 "github.com/openshift/api/config/v1"
)

// Load ClusterVersion object from file
func LoadClusterVersionFromFile(filePath string) (*configv1.ClusterVersion, error) {
	data, err := FetchURL(filePath)
	if err != nil {
		return nil, err
	}

	var cvList configv1.ClusterVersionList
	if err := json.Unmarshal([]byte(data), &cvList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal as ClusterVersionList: %w", err)
	}

	if len(cvList.Items) == 0 {
		return nil, fmt.Errorf("no ClusterVersion object found in list")
	}

	// Only one ClusterVersion should exist
	return &cvList.Items[0], nil
}
