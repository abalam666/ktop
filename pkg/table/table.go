package table

import (
	"image"

	corev1 "k8s.io/api/core/v1"

	"github.com/ynqa/ktop/pkg/formats"
	"github.com/ynqa/ktop/pkg/resources"
	"github.com/ynqa/ktop/pkg/state"
	"github.com/ynqa/ktop/pkg/ui"
)

type Contents struct {
	Headers []string
	Widths  []int
}

type Shaper interface {
	Headers() []string
	Widths(image.Rectangle) []int
	Rows(resources.Resources, *state.TableVisibleSet) []ui.Row
}

type NopShaper struct{}

func (*NopShaper) Headers() []string {
	return []string{"message"}
}

func (*NopShaper) Widths(rect image.Rectangle) []int {
	return []int{rect.Dx() - 1}
}

func (*NopShaper) Rows(resources.Resources, *state.TableVisibleSet) []ui.Row {
	return []ui.Row{
		{
			Elems: []string{"not found: nodes, pods, and containers"},
		},
	}
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

func (*KubeShaper) Rows(r resources.Resources, state *state.TableVisibleSet) []ui.Row {
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
