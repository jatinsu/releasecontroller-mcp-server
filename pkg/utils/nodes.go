package utils

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	corev1 "k8s.io/api/core/v1"
)

func LoadNodesFromFile(path string) ([]corev1.Node, error) {
	bytes, err := FetchURL(path)
	if err != nil {
		return nil, err
	}

	var nodeList corev1.NodeList
	err = json.Unmarshal([]byte(bytes), &nodeList)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nodeList.Items, nil
}

func FindNodeByName(nodes []corev1.Node, name string) (*corev1.Node, error) {
	for _, node := range nodes {
		if node.Name == name {
			return &node, nil
		}
	}
	return nil, fmt.Errorf("node %s not found", name)
}

// GetNodeInfoString safely extracts and returns NodeInfo from a Node object as a formatted string
func GetNodeInfoString(node *corev1.Node) string {
	if node == nil {
		return "Node is nil"
	}

	info := node.Status.NodeInfo
	// If NodeInfo is empty (all fields are zero values), this indicates it wasn't populated.
	if info.MachineID == "" && info.SystemUUID == "" && info.KernelVersion == "" && info.OSImage == "" {
		return fmt.Sprintf("NodeInfo is not available for node %s", node.Name)
	}

	return fmt.Sprintf(
		`Node Info for %s:
  Architecture: %s
  Boot ID: %s
  Container Runtime Version: %s
  Kernel Version: %s
  KubeProxy Version: %s
  Kubelet Version: %s
  Machine ID: %s
  Operating System: %s
  OS Image: %s
  System UUID: %s`,
		node.Name,
		safe(info.Architecture),
		safe(info.BootID),
		safe(info.ContainerRuntimeVersion),
		safe(info.KernelVersion),
		safe(info.KubeProxyVersion),
		safe(info.KubeletVersion),
		safe(info.MachineID),
		safe(info.OperatingSystem),
		safe(info.OSImage),
		safe(info.SystemUUID),
	)
}

// safe returns a placeholder if the string is empty
func safe(field string) string {
	if field == "" {
		return "<not available>"
	}
	return field
}

// GetNodeLabelsString returns all labels from node metadata as a string
func GetNodeLabelsString(node *corev1.Node) string {
	if node == nil {
		return "Node is nil"
	}

	if len(node.Labels) == 0 {
		return fmt.Sprintf("No labels found for node %s", node.Name)
	}

	var keys []string
	for k := range node.Labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	output := fmt.Sprintf("Labels for node %s:\n", node.Name)
	for _, k := range keys {
		output += fmt.Sprintf("  %s: %s\n", k, node.Labels[k])
	}
	return output
}

// GetNodeAnnotationsString returns all annotations from node metadata as a string
func GetNodeAnnotationsString(node *corev1.Node) string {
	if node == nil {
		return "Node is nil"
	}

	if len(node.Annotations) == 0 {
		return fmt.Sprintf("No annotations found for node %s", node.Name)
	}

	var keys []string
	for k := range node.Annotations {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	output := fmt.Sprintf("Annotations for node %s:\n", node.Name)
	for _, k := range keys {
		output += fmt.Sprintf("  %s: %s\n", k, node.Annotations[k])
	}
	return output
}

// GetNodeConditionsString returns node condition information as a formatted string
func GetNodeConditionsString(node *corev1.Node) string {
	if node == nil {
		return "Node is nil"
	}

	if len(node.Status.Conditions) == 0 {
		return fmt.Sprintf("No conditions found for node %s", node.Name)
	}

	output := fmt.Sprintf("Conditions for node %s:\n", node.Name)
	for _, cond := range node.Status.Conditions {
		output += fmt.Sprintf("  - Type: %s\n", cond.Type)
		output += fmt.Sprintf("    Status: %s\n", cond.Status)
		output += fmt.Sprintf("    Reason: %s\n", cond.Reason)
		output += fmt.Sprintf("    Message: %s\n", cond.Message)
		output += fmt.Sprintf("    LastTransitionTime: %s\n", cond.LastTransitionTime.Format(time.RFC3339))
	}
	return output
}
