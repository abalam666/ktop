package state

import (
	"strings"
	"sync"

	corev1 "k8s.io/api/core/v1"

	"github.com/ynqa/ktop/pkg/formats"
	"github.com/ynqa/ktop/pkg/resources"
)

type Kind string

const (
	Node      = Kind("node")
	Pod       = Kind("pod")
	Container = Kind("container")
)

type content struct {
	kind      Kind
	namespace string
	name      string
	usage     corev1.ResourceList
}

type ViewState struct {
	mu              sync.RWMutex
	contents        []*content
	childVisibleSet map[string]struct{}
}

func New() *ViewState {
	return &ViewState{
		contents:        make([]*content, 0),
		childVisibleSet: make(map[string]struct{}),
	}
}

func (v *ViewState) Len() int {
	return len(v.contents)
}

func (v *ViewState) Update(r resources.Resources) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.Len() > 0 {
		v.update(gencontents(r))
		v.contents = gencontents(r)
	} else {
		v.contents = gencontents(r)
	}
}

func (v *ViewState) update(new []*content) {
	// for _, content := range new {

	// }
	// v.contents = new
}

func (v *ViewState) Toggle(idx int) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.contents[idx].kind == Node || v.contents[idx].kind == Pod {
		v.toggleChildVisible(v.contents[idx])
	}
}

func (v *ViewState) toggleChildVisible(c *content) {
	key := genkey(c)
	_, ok := v.childVisibleSet[key]
	if ok {
		delete(v.childVisibleSet, key)
	} else {
		v.childVisibleSet[key] = struct{}{}
	}
}

func (v *ViewState) ToRows() [][]string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	var rows [][]string
	for _, content := range v.contents {
		rows = append(rows, genrow(content, v.childVisible(content)))
	}
	return rows
}

func (v *ViewState) childVisible(c *content) bool {
	_, ok := v.childVisibleSet[genkey(c)]
	return ok
}

func genkey(c *content) string {
	return strings.Join([]string{string(c.kind), c.namespace, c.name}, "_")
}

func genrow(c *content, childVisible bool) []string {
	var row []string
	switch c.kind {
	case Node:
		row = []string{
			formats.FormatNodeNameField(c.name, childVisible),
			formats.FormatResource(corev1.ResourceCPU, c.usage),
			formats.FormatResource(corev1.ResourceMemory, c.usage),
		}
	case Pod:
		row = []string{
			formats.FormatPodNameField(c.name, childVisible),
			formats.FormatResource(corev1.ResourceCPU, c.usage),
			formats.FormatResource(corev1.ResourceMemory, c.usage),
		}
	case Container:
		row = []string{
			formats.FormatContainerNameField(c.name),
			formats.FormatResource(corev1.ResourceCPU, c.usage),
			formats.FormatResource(corev1.ResourceMemory, c.usage),
		}
	}
	return row
}

func gencontents(r resources.Resources) []*content {
	var order []*content
	for _, node := range r.SortedNodes() {
		pods := r.SortedPods(node)
		order = append(order, &content{
			kind:  Node,
			name:  node,
			usage: r[node].Usage,
		})
		for _, pod := range pods {
			containers := r.SortedContainers(node, pod)
			order = append(order, &content{
				kind:  Pod,
				name:  pod,
				usage: r[node].Pods[pod].Usage,
			})
			for _, container := range containers {
				order = append(order, &content{
					kind:  Container,
					name:  container,
					usage: r[node].Pods[pod].Containers[container].Usage,
				})
			}
		}
	}
	return order
}
