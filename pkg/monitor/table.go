package monitor

import (
	"image"

	"github.com/ynqa/ktop/pkg/resources"
)

type viewer interface {
	fields(resources.Resources, image.Rectangle) fields
}

type fields struct {
	Headers []string
	Widths  []int
	Rows    [][]string
}

type emptyviewer struct{}

func (*emptyviewer) fields(
	_ resources.Resources,
	rect image.Rectangle,
) fields {
	return fields{
		Headers: []string{"message"},
		Widths:  []int{rect.Dx() - 1},
		Rows:    [][]string{{"no node, pods, and containers"}},
	}
}

type simpleviewer struct{}

func (*simpleviewer) fields(
	data resources.Resources,
	rect image.Rectangle,
) fields {
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
	return fields{
		Headers: headers,
		Widths:  widths,
		Rows:    rows,
	}
}
