package graph

import (
	"github.com/ynqa/ktop/pkg/formats"
	"github.com/ynqa/ktop/pkg/resources"
	"github.com/ynqa/ktop/pkg/ui"
)

type item struct {
	labelHeader string
}

type Contents struct {
	set map[string]item
}

func NewForResources(r resources.Resources) Contents {
	set := map[string]item{}
	for node, noder := range r {
		nodeKey := formats.FormatNodeStateKey(node)
		set[nodeKey] = item{
			labelHeader: formats.FormatLabelHeader(node),
		}
		for pod, podr := range noder.Pods {
			podKey := formats.FormatPodStateKey(node, podr.Namespace, pod)
			set[podKey] = item{
				labelHeader: formats.FormatLabelHeader(pod),
			}
			for container := range podr.Containers {
				containerKey := formats.FormatContainerStateKey(node, podr.Namespace, pod, container)
				set[containerKey] = item{
					labelHeader: formats.FormatLabelHeader(container),
				}
			}
		}
	}
	return Contents{
		set: set,
	}
}

func (c *Contents) Len() int {
	return len(c.set)
}

type Drawer interface {
	Draw(*ui.Graph, Contents, string)
}

type NopDrawer struct{}

func (*NopDrawer) Draw(g *ui.Graph, _ Contents, _ string) {
	g.UpperLimit = 100
	g.Data = append(g.Data, 50)
}

type KubeDrawer struct{}

func (d *KubeDrawer) Draw(g *ui.Graph, c Contents, key string) {
	g.UpperLimit = 100
	g.Data = append(g.Data, 50)
	g.LabelHeader = c.set[key].labelHeader
}
