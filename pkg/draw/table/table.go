package table

import (
	"image"

	"github.com/ynqa/widgets/pkg/node"
	"github.com/ynqa/widgets/pkg/widgets"

	"github.com/ynqa/ktop/pkg/resources"
)

type Drawer interface {
	Draw(*widgets.Table, resources.Resources)
	headers() []string
	widths(image.Rectangle) []int
}

type NopDrawer struct{}

func (d *NopDrawer) Draw(table *widgets.Table, r resources.Resources) {
	table.Headers = d.headers()
	table.Widths = d.widths(table.Inner)
	table.Node = node.New("", []string{"not found: nodes, pods, and containers"})
}

func (*NopDrawer) headers() []string {
	return []string{"message"}
}

func (*NopDrawer) widths(rect image.Rectangle) []int {
	return []int{rect.Dx() - 1}
}

func (*NopDrawer) node(resources.Resources) *node.Node {
	return node.New("", []string{"not found: nodes, pods, and containers"})
}

type KubeDrawer struct{}

func (d *KubeDrawer) Draw(table *widgets.Table, r resources.Resources) {
	table.Headers = d.headers()
	table.Widths = d.widths(table.Inner)
	table.Node = node.ApplyChildVisible(table.Node, r.GetTree())
}

func (*KubeDrawer) headers() []string {
	return []string{"name", "namespace", "usage.cpu", "usage.memory"}
}

func (s *KubeDrawer) widths(rect image.Rectangle) []int {
	widths := []int{rect.Dx() / 2}
	denom := 2 * (len(s.headers()) - 1)
	for i := 1; i < len(s.headers()); i++ {
		widths = append(widths, rect.Dx()/denom)
	}
	return widths
}
