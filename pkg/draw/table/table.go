package table

import (
	"image"

	"github.com/ynqa/widgets"
	"github.com/ynqa/widgets/node"

	corev1 "k8s.io/api/core/v1"

	"github.com/ynqa/ktop/pkg/formats"
	"github.com/ynqa/ktop/pkg/resources"
)

type Drawer interface {
	Draw(*widgets.ToggleTable, resources.Resources)
	headers() []string
	widths(image.Rectangle) []int
	node(resources.Resources) *node.Node
}

type NopDrawer struct{}

func (d *NopDrawer) Draw(table *widgets.ToggleTable, r resources.Resources) {
	table.Headers = d.headers()
	table.Widths = d.widths(table.Inner)
	table.Node = d.node(r)
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

func (d *KubeDrawer) Draw(table *widgets.ToggleTable, r resources.Resources) {
	table.Headers = d.headers()
	table.Widths = d.widths(table.Inner)
	table.Node = node.ApplyChildVisible(table.Node, d.node(r))
}

func (*KubeDrawer) headers() []string {
	return []string{"name", "namespace", "usage.cpu", "usage.memory"}
}

func (s *KubeDrawer) widths(rect image.Rectangle) []int {
	widths := []int{rect.Dx() / 2}
	for i := 1; i < len(s.headers()); i++ {
		denom := 2 * (len(s.headers()) - 1)
		widths = append(widths, rect.Dx()/denom)
	}
	return widths
}

func (*KubeDrawer) node(r resources.Resources) *node.Node {
	tree := node.Root()
	for _, nd := range r.SortedNodes() {
		usage := r[nd].Usage
		cursorNode := node.New(nd, []string{
			nd,
			"",
			formats.FormatResource(corev1.ResourceCPU, usage),
			formats.FormatResource(corev1.ResourceMemory, usage),
		})
		tree.Append(cursorNode)

		for _, pd := range r.SortedPods(nd) {
			usage := r[nd].Pods[pd].Usage
			cursorPod := node.New(pd, []string{
				pd,
				r[nd].Pods[pd].Namespace,
				formats.FormatResource(corev1.ResourceCPU, usage),
				formats.FormatResource(corev1.ResourceMemory, usage),
			})
			cursorNode.Append(cursorPod)

			for _, ct := range r.SortedContainers(nd, pd) {
				usage = r[nd].Pods[pd].Containers[ct].Usage
				cursorPod.Append(node.Leaf(ct, []string{
					ct,
					"",
					formats.FormatResource(corev1.ResourceCPU, usage),
					formats.FormatResource(corev1.ResourceMemory, usage),
				}))
			}
		}
	}
	return tree
}
