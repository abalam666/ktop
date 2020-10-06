package graph

import (
	"sync"

	"github.com/ynqa/ktop/pkg/formats"
	"github.com/ynqa/ktop/pkg/resources"
	"github.com/ynqa/ktop/pkg/ui"
)

type item struct {
	labelHeader string
}

type VisibleSet struct {
	mu  sync.RWMutex
	set map[string]item
}

func NewVisibleSet() *VisibleSet {
	return &VisibleSet{
		set: make(map[string]item),
	}
}

func (s *VisibleSet) Add(r resources.Resources) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for node, noder := range r {
		nodeKey := formats.FormatNodeStateKey(node)
		s.set[nodeKey] = item{
			labelHeader: node,
		}
		for pod, podr := range noder.Pods {
			podKey := formats.FormatPodStateKey(node, podr.Namespace, pod)
			s.set[podKey] = item{
				labelHeader: pod,
			}
			for container := range podr.Containers {
				containerKey := formats.FormatContainerStateKey(node, podr.Namespace, pod, container)
				s.set[containerKey] = item{
					labelHeader: container,
				}
			}
		}
	}
}

func (s *VisibleSet) Pick(key string) item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.set[key]
}

type Drawer interface {
	Draw(*ui.Graph, string)
}

type NopDrawer struct{}

func (*NopDrawer) Draw(g *ui.Graph, _ string) {
	g.UpperLimit = 100
	g.Data = append(g.Data, 50)
}

type KubeDrawer struct{}

func (*KubeDrawer) Draw(g *ui.Graph, _ string) {
	g.UpperLimit = 100
	g.Data = append(g.Data, 50)
}
