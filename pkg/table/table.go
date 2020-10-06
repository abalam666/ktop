package table

import (
	"image"

	"github.com/ynqa/ktop/pkg/state"
)

type Contents struct {
	Headers []string
	Widths  []int
}

type Shaper interface {
	Headers() []string
	Widths(image.Rectangle) []int
	Rows(*state.ViewState) [][]string
}

type NopShaper struct{}

func (*NopShaper) Headers() []string {
	return []string{"message"}
}

func (*NopShaper) Widths(rect image.Rectangle) []int {
	return []int{rect.Dx() - 1}
}

func (*NopShaper) Rows(*state.ViewState) [][]string {
	return [][]string{{"not found: nodes, pods, and containers"}}
}

type KubeShaper struct{}

func (*KubeShaper) Headers() []string {
	return []string{"name", "namespace", "usage.cpu", "usage.memory"}
}

func (s *KubeShaper) Widths(rect image.Rectangle) []int {
	widths := []int{rect.Dx() / 2}
	for i := 1; i < len(s.Headers()); i++ {
		denom := 2 * (len(s.Headers()) - 1)
		widths = append(widths, rect.Dx()/denom)
	}
	return widths
}

func (*KubeShaper) Rows(state *state.ViewState) [][]string {
	return state.ToRows()
}
