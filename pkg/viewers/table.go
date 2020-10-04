package viewers

import (
	"image"

	"github.com/ynqa/ktop/pkg/resources"
)

type Table interface {
	Fields(image.Rectangle, resources.Resources) Fields
}

type Fields struct {
	Headers []string
	Widths  []int
	Rows    [][]string
}

type EmptyTable struct{}

func (*EmptyTable) Fields(rect image.Rectangle, _ resources.Resources) Fields {
	return Fields{
		Headers: []string{"message"},
		Widths:  []int{rect.Dx() - 1},
		Rows:    [][]string{{"no node, pods, and containers"}},
	}
}

type ResourceTable struct{}

func (*ResourceTable) Fields(rect image.Rectangle, resources resources.Resources) Fields {
	headers := []string{
		"metadata.name", "usage.cpu", "usage.memory",
	}
	base := rect.Dx() - 1*len(headers)
	widths := []int{
		base / 2, base / 4, base / 4,
	}
	var rows [][]string
	nodes := resources.SortedNodes()
	for _, node := range nodes {
		usage := resources[node].Usage
		rows = append(rows, []string{
			node, usage.Cpu().String(), usage.Memory().String(),
		})
	}
	return Fields{
		Headers: headers,
		Widths:  widths,
		Rows:    rows,
	}
}
