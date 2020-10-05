package table

import (
	"image"

	corev1 "k8s.io/api/core/v1"

	"github.com/ynqa/ktop/pkg/resources"
	"github.com/ynqa/ktop/pkg/table/formats"
	"github.com/ynqa/ktop/pkg/table/state"
)

type Contents struct {
	Headers []string
	Widths  []int
	Rows    [][]string
}

type ContentsCreator interface {
	Create(resources.Resources, *state.ChildVisibleSet, image.Rectangle) Contents
}

type NoContentsCreator struct{}

func (*NoContentsCreator) Create(_ resources.Resources, _ *state.ChildVisibleSet, rect image.Rectangle) Contents {
	return Contents{
		Headers: []string{"message"},
		Widths:  []int{rect.Dx() - 1},
		Rows:    [][]string{{"no node, pods, and containers"}},
	}
}

type KubeResourceContentsCreator struct{}

func (*KubeResourceContentsCreator) Create(data resources.Resources, childVisibleSet *state.ChildVisibleSet, rect image.Rectangle) Contents {
	headers := []string{"metadata.name", "usage.cpu", "usage.memory"}

	// estimate width for columns
	widths := []int{rect.Dx() / 2}
	for i := 1; i < len(headers); i++ {
		denom := 2 * (len(headers) - 1)
		widths = append(widths, rect.Dx()/denom)
	}

	// generate rows to view resources
	var (
		rows         [][]string
		childVisible bool
	)
	for _, node := range data.SortedNodes() {
		usage := data[node].Usage
		childVisible = childVisibleSet.Contains(node)
		rows = append(rows, []string{
			formats.FormatNodeNameField(node, childVisible),
			formats.FormatResource(corev1.ResourceCPU, usage),
			formats.FormatResource(corev1.ResourceMemory, usage),
		})
		if childVisible {
			for _, pod := range data.SortedPods(node) {
				usage := data[node].Pods[pod].Usage
				childVisible = childVisibleSet.Contains(pod)
				rows = append(rows, []string{
					formats.FormatPodNameField(pod, childVisible),
					formats.FormatResource(corev1.ResourceCPU, usage),
					formats.FormatResource(corev1.ResourceMemory, usage),
				})
				if childVisible {
					for _, container := range data.SortedContainers(node, pod) {
						usage = data[node].Pods[pod].Containers[container].Usage
						rows = append(rows, []string{
							formats.FormatContainerNameField(container),
							formats.FormatResource(corev1.ResourceCPU, usage),
							formats.FormatResource(corev1.ResourceMemory, usage),
						})
					}
				}
			}
		}
	}

	return Contents{
		Headers: headers,
		Widths:  widths,
		Rows:    rows,
	}
}
