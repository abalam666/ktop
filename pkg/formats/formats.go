package formats

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

func FormatLabelHeader(name string) string {
	return fmt.Sprintf("name: %s", name)
}

func FormatResource(name corev1.ResourceName, list corev1.ResourceList) string {
	r, ok := list[name]
	switch {
	case ok && name == corev1.ResourceCPU:
		return fmt.Sprintf("%vm", r.MilliValue())
	case ok && name == corev1.ResourceMemory:
		return fmt.Sprintf("%vMi", r.Value()/(1024*1024))
	default:
		return "-"
	}
}
