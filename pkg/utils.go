package pkg

import (
	v1 "k8s.io/api/core/v1"
	"strings"
)

const NodeOrderAnnotation = "scheduler.alpha.kubernetes.io/node-order"

func GetNodeOrderFromPod(pod *v1.Pod) []string {
	if pod == nil {
		return nil
	}

	v, ok := pod.Annotations[NodeOrderAnnotation]
	if !ok || v == "" {
		return nil
	}

	parts := strings.Split(v, ",")
	var nodes []string
	for _, n := range parts {
		n = strings.TrimSpace(n)
		if n != "" {
			nodes = append(nodes, n)
		}
	}
	return nodes
}
