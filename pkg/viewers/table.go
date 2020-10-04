package viewers

import (
	"image"

	"github.com/ynqa/ktop/pkg/resources"
)

type Table interface {
	Fields(image.Rectangle, resources.Resources) Fields
}

type Fields struct {
	Title   string
	Headers []string
	Widths  []int
	Rows    [][]string
}

type EmptyTable struct{}

func (*EmptyTable) Fields(rect image.Rectangle, _ resources.Resources) Fields {
	return Fields{
		Title:   "metadata.name",
		Headers: []string{"message"},
		Widths:  []int{rect.Dx() - 1},
		Rows:    [][]string{{"no node, pods, and containers"}},
	}
}
