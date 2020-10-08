package table

import (
	"image"
	"sync"

	corev1 "k8s.io/api/core/v1"

	"github.com/ynqa/ktop/pkg/formats"
	"github.com/ynqa/ktop/pkg/resources"
	"github.com/ynqa/ktop/pkg/ui"
)

type VisibleSet struct {
	mu  sync.RWMutex
	set map[string]struct{}
}

func NewVisibleSet() *VisibleSet {
	return &VisibleSet{
		set: make(map[string]struct{}),
	}
}

func (s *VisibleSet) Contains(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.set[name]
	return ok
}

func (s *VisibleSet) Toggle(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.set[name]
	if ok {
		delete(s.set, name)
	} else {
		s.set[name] = struct{}{}
	}
}

type Drawer interface {
	Draw(*ui.Table, resources.Resources, *VisibleSet)
	headers() []string
	widths(image.Rectangle) []int
	rows(resources.Resources, *VisibleSet) []ui.Row
}

type NopDrawer struct{}

func (d *NopDrawer) Draw(table *ui.Table, r resources.Resources, state *VisibleSet) {
	table.Headers = d.headers()
	table.Widths = d.widths(table.Inner)
	table.Rows = d.rows(r, state)
}

func (*NopDrawer) headers() []string {
	return []string{"message"}
}

func (*NopDrawer) widths(rect image.Rectangle) []int {
	return []int{rect.Dx() - 1}
}

func (*NopDrawer) rows(resources.Resources, *VisibleSet) []ui.Row {
	return []ui.Row{
		{
			Elems: []string{"not found: nodes, pods, and containers"},
		},
	}
}

type KubeDrawer struct{}

func (d *KubeDrawer) Draw(table *ui.Table, r resources.Resources, state *VisibleSet) {
	table.Headers = d.headers()
	table.Widths = d.widths(table.Inner)
	table.Rows = d.rows(r, state)
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

func (*KubeDrawer) rows(r resources.Resources, state *VisibleSet) []ui.Row {
	var rows []ui.Row
	for _, node := range r.SortedNodes() {
		usage := r[node].Usage
		nodeKey := formats.FormatNodeStateKey(node)
		childVisible := state.Contains(nodeKey)
		rows = append(rows, ui.Row{
			Key: nodeKey,
			Elems: []string{
				formats.FormatNodeNameField(node, childVisible),
				"",
				formats.FormatResource(corev1.ResourceCPU, usage),
				formats.FormatResource(corev1.ResourceMemory, usage),
			},
		})
		if childVisible {
			for _, pod := range r.SortedPods(node) {
				ns := r[node].Pods[pod].Namespace
				podKey := formats.FormatPodStateKey(node, ns, pod)
				usage := r[node].Pods[pod].Usage
				childVisible := state.Contains(podKey)
				rows = append(rows, ui.Row{
					Key: podKey,
					Elems: []string{
						formats.FormatPodNameField(pod, childVisible),
						ns,
						formats.FormatResource(corev1.ResourceCPU, usage),
						formats.FormatResource(corev1.ResourceMemory, usage),
					},
				})
				if childVisible {
					for _, container := range r.SortedContainers(node, pod) {
						ns := r[node].Pods[pod].Namespace
						containerKey := formats.FormatContainerStateKey(node, ns, pod, container)
						usage = r[node].Pods[pod].Containers[container].Usage
						rows = append(rows, ui.Row{
							Key: containerKey,
							Elems: []string{
								formats.FormatContainerNameField(container),
								ns,
								formats.FormatResource(corev1.ResourceCPU, usage),
								formats.FormatResource(corev1.ResourceMemory, usage),
							},
						})
					}
				}
			}
		}
	}
	return rows
}
