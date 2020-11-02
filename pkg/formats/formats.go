package formats

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func FormatResourcePercentage(typ corev1.ResourceName, usage, available corev1.ResourceList) string {
	var u, a int64
	switch typ {
	case corev1.ResourceCPU:
		u = usage.Cpu().MilliValue()
		a = available.Cpu().MilliValue()
	case corev1.ResourceMemory:
		u = usage.Memory().MilliValue()
		a = available.Memory().MilliValue()
	}
	return fmt.Sprintf("%.1f%%", (float64(u)/float64(a))*100)
}

func FormatResource(typ corev1.ResourceName, list corev1.ResourceList) int64 {
	switch typ {
	case corev1.ResourceCPU:
		return list.Cpu().MilliValue()
	case corev1.ResourceMemory:
		return list.Memory().ScaledValue(resource.Mega)
	default:
		return 0
	}
}

func FormatResourceString(typ corev1.ResourceName, list corev1.ResourceList) string {
	switch typ {
	case corev1.ResourceCPU:
		return fmt.Sprintf("%vm", FormatResource(typ, list))
	case corev1.ResourceMemory:
		return fmt.Sprintf("%vMi", FormatResource(typ, list))
	default:
		return "-"
	}
}
