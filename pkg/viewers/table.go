package viewers

import (
	"image"

	"github.com/ynqa/ktop/pkg/resources"
)

type Table interface {
	Fields(resources.Resources, image.Rectangle) Fields
}

type Fields struct {
	Headers []string
	Widths  []int
	Rows    [][]string
}

type EmptyTable struct{}

func (*EmptyTable) Fields(
	_ resources.Resources,
	rect image.Rectangle,
) Fields {
	return Fields{
		Headers: []string{"message"},
		Widths:  []int{rect.Dx() - 1},
		Rows:    [][]string{{"no node, pods, and containers"}},
	}
}

type ResourceTable struct{}

func (*ResourceTable) Fields(
	data resources.Resources,
	rect image.Rectangle,
) Fields {
	headers := []string{
		"metadata.name", "usage.cpu", "usage.memory",
	}
	base := rect.Dx() - 1*len(headers)
	widths := []int{
		base / 2, base / 4, base / 4,
	}
	var rows [][]string
	for _, node := range data.SortedNodes() {
		usage := data[node].Usage
		rows = append(rows, []string{
			node, usage.Cpu().String(), usage.Memory().String(),
		})
		for _, pod := range data.SortedPods(node) {
			rows = append(rows, []string{
				pod, usage.Cpu().String(), usage.Memory().String(),
			})
			for _, container := range data.SortedContainers(node, pod) {
				usage = data[node].Pods[pod].Containers[container].Usage
				rows = append(rows, []string{
					container, usage.Cpu().String(), usage.Memory().String(),
				})
			}
		}
	}
	return Fields{
		Headers: headers,
		Widths:  widths,
		Rows:    rows,
	}
}
