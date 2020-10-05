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
	Create(resources.Resources, *state.VisibleSet, image.Rectangle) Contents
}

type NopCreator struct{}

func (*NopCreator) Create(_ resources.Resources, _ *state.VisibleSet, rect image.Rectangle) Contents {
	return Contents{
		Headers: []string{"message"},
		Widths:  []int{rect.Dx() - 1},
		Rows:    [][]string{{"no node, pods, and containers"}},
	}
}

type Creator struct{}

func (*Creator) Create(data resources.Resources, _ *state.VisibleSet, rect image.Rectangle) Contents {
	headers := []string{"metadata.name", "usage.cpu", "usage.memory"}

	widths := []int{rect.Dx() / 2}
	for i := 1; i < len(headers); i++ {
		denom := 2 * (len(headers) - 1)
		widths = append(widths, rect.Dx()/denom)
	}

	// generate rows
	var rows [][]string
	for _, node := range data.SortedNodes() {
		usage := data[node].Usage
		rows = append(rows, []string{
			formats.FormatNodeNameField(node),
			formats.FormatResource(corev1.ResourceCPU, usage),
			formats.FormatResource(corev1.ResourceMemory, usage),
		})
		for _, pod := range data.SortedPods(node) {
			usage := data[node].Pods[pod].Usage
			rows = append(rows, []string{
				formats.FormatPodNameField(pod),
				formats.FormatResource(corev1.ResourceCPU, usage),
				formats.FormatResource(corev1.ResourceMemory, usage),
			})
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

	return Contents{
		Headers: headers,
		Widths:  widths,
		Rows:    rows,
	}
}
