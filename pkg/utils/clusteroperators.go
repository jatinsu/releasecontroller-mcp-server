package utils

import (
	"encoding/json"
	"fmt"

	configv1 "github.com/openshift/api/config/v1"
)

// Load JSON array of ClusterOperators from file
func LoadClusterOperatorsFromFile(filePath string) ([]configv1.ClusterOperator, error) {
	data, err := FetchURL(filePath)
	if err != nil {
		return nil, err
	}
	var coList configv1.ClusterOperatorList
	if err := json.Unmarshal([]byte(data), &coList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal as ClusterOperatorList: %w", err)
	}

	return coList.Items, nil
}
