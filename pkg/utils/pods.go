package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

func LoadPodsFromFile(path string) ([]corev1.Pod, error) {
	bytes, err := FetchURL(path)
	if err != nil {
		return nil, err
	}

	var podList corev1.PodList
	err = json.Unmarshal([]byte(bytes), &podList)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return podList.Items, nil
}

func AllPodsSummary(pods []corev1.Pod) string {
	var b strings.Builder
	for _, pod := range pods {
		fmt.Fprintf(&b, "%s/%s on %s: %s\n", pod.Namespace, pod.Name, pod.Spec.NodeName, pod.Status.Phase)
	}
	return b.String()
}

func RunningPodsSummary(pods []corev1.Pod) string {
	var b strings.Builder
	for _, pod := range pods {
		if pod.Status.Phase == corev1.PodRunning {
			fmt.Fprintf(&b, "%s/%s on %s: Running\n", pod.Namespace, pod.Name, pod.Spec.NodeName)
		}
	}
	if b.Len() > 0 {
		return b.String()
	} else {
		return "No pods in Running state."
	}
	return b.String()
}

func CrashLoopBackOffSummary(pods []corev1.Pod) string {
	var b strings.Builder
	for _, pod := range pods {
		for _, cs := range pod.Status.ContainerStatuses {
			if cs.State.Waiting != nil && cs.State.Waiting.Reason == "CrashLoopBackOff" {
				fmt.Fprintf(&b, "%s/%s on %s: %s (%s)\n", pod.Namespace, pod.Name, pod.Spec.NodeName, cs.State.Waiting.Reason, cs.State.Waiting.Message)
			}
		}
	}
	if b.Len() > 0 {
		return b.String()
	} else {
		return "No pods in CrashLoopBackOff state."
	}
	return b.String()
}

func PendingPodsSummary(pods []corev1.Pod) string {
	var b strings.Builder
	for _, pod := range pods {
		if pod.Status.Phase == corev1.PodPending {
			reason := pod.Status.Reason
			if reason == "" {
				// Fallback: Try to get reason from pod conditions
				for _, cond := range pod.Status.Conditions {
					if cond.Type == corev1.PodScheduled && cond.Status == corev1.ConditionFalse {
						reason = cond.Reason
						break
					}
				}
			}
			if reason == "" {
				reason = "Unknown"
			}
			fmt.Fprintf(&b, "%s/%s on %s: Pending (%s)\n", pod.Namespace, pod.Name, pod.Spec.NodeName, reason)
		}
	}
	if b.Len() > 0 {
		return b.String()
	} else {
		return "No pods in Pending state."
	}
	return b.String()
}

func InitStateSummary(pods []corev1.Pod) string {
	var b strings.Builder
	for _, pod := range pods {
		for _, cs := range pod.Status.InitContainerStatuses {
			if cs.State.Waiting != nil && strings.Contains(cs.State.Waiting.Reason, "Init") {
				fmt.Fprintf(&b, "%s/%s on %s: %s\n", pod.Namespace, pod.Name, pod.Spec.NodeName, cs.State.Waiting.Reason)
				break
			}
		}
	}
	if b.Len() > 0 {
		return b.String()
	} else {
		return "No pods in Init state."
	}
	return b.String()
}

func ErrorStateSummary(pods []corev1.Pod) string {
	var b strings.Builder
	for _, pod := range pods {
		for _, cs := range pod.Status.ContainerStatuses {
			if cs.State.Waiting != nil && cs.State.Waiting.Reason == "Error" {
				fmt.Fprintf(&b, "%s/%s: %s\n", pod.Namespace, pod.Name, cs.State.Waiting.Reason)
				break
			}
		}
	}
	if b.Len() > 0 {
		return b.String()
	} else {
		return "No pods in Error state."
	}
	return b.String()
}

func FilterPodsByNamespaceAsString(pods []corev1.Pod, namespace string) string {
	var b strings.Builder
	for _, pod := range pods {
		if pod.Namespace == namespace {
			fmt.Fprintf(&b, "%s/%s on %s: %s\n", pod.Namespace, pod.Name, pod.Spec.NodeName, pod.Status.Phase)
		}
	}
	if b.Len() > 0 {
		return b.String()
	} else {
		return fmt.Sprintf("No pods found in namespace %s.", namespace)
	}
}

func FilterPodsByNodeAsString(pods []corev1.Pod, nodeName string) string {
	var b strings.Builder
	for _, pod := range pods {
		if pod.Spec.NodeName == nodeName {
			fmt.Fprintf(&b, "%s/%s on %s: %s\n", pod.Namespace, pod.Name, pod.Spec.NodeName, pod.Status.Phase)
		}
	}
	if b.Len() > 0 {
		return b.String()
	} else {
		return fmt.Sprintf("No pods found on node %s.", nodeName)
	}
}
