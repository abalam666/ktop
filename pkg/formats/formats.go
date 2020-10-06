package formats

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

const (
	rightArrow  = "▶"
	downArrow   = "▼"
)

func withSpaces(name string, indent int) string {
	return strings.Repeat(" ", len(rightArrow)+indent) + name
}

func withRightArrow(name string, indent int) string {
	return strings.Repeat(" ", indent) + rightArrow + name
}

func withDownArrow(name string, indent int) string {
	return strings.Repeat(" ", indent) + downArrow + name
}

func FormatNodeNameField(name string, childVisible bool) string {
	if childVisible {
		return withDownArrow(name, 0)
	} else {
		return withRightArrow(name, 0)
	}
}

func FormatPodNameField(name string, childVisible bool) string {
	if childVisible {
		return withDownArrow(name, 1)
	} else {
		return withRightArrow(name, 1)
	}
}

func FormatContainerNameField(name string) string {
	return withSpaces(name, 1)
}

func TrimString(name string) string {
	name = strings.TrimLeft(name, " ")
	name = strings.TrimLeft(name, rightArrow)
	name = strings.TrimLeft(name, downArrow)
	return name
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
