package graph

import (
	"fmt"

	"github.com/gizak/termui/v3"
	"github.com/ynqa/widgets/pkg/node"
	corev1 "k8s.io/api/core/v1"

	"github.com/ynqa/ktop/pkg/formats"
	"github.com/ynqa/ktop/pkg/resources"
	"github.com/ynqa/ktop/pkg/ui"
)

type Drawer interface {
	Draw(*ui.Graph, resources.Resources, corev1.ResourceName, []*node.Node)
}

type NopDrawer struct{}

func (*NopDrawer) Draw(*ui.Graph, resources.Resources, corev1.ResourceName, []*node.Node) {}

type KubeDrawer struct{}

func (d *KubeDrawer) Draw(g *ui.Graph, r resources.Resources, typ corev1.ResourceName, nodes []*node.Node) {
	if len(nodes) == 1 {
		node, ok := r.GetNodeResource(nodes[0].Name())
		// fmt.Sprintf("usage (%v) / allocatable (%v) = %v",
		// 	formats.FormatResourceString(typ, node.Usage),
		// 	formats.FormatResourceString(typ, node.Allocatable),
		// 	formats.FormatResourcePercentage(typ, node.Usage, node.Allocatable),
		// )
		if ok {
			g.NodeLimit = ui.Straight{
				Value: float64(formats.FormatResource(typ, node.Allocatable)),
				Label: fmt.Sprintf("allocatable: %v", formats.FormatResourceString(typ, node.Allocatable)),
				Style: termui.NewStyle(termui.ColorRed, termui.ColorClear),
			}
			g.Usage = ui.Curved{
				Values: append(g.Usage.Values, float64(formats.FormatResource(typ, node.Usage))),
				Label:  fmt.Sprintf("usage: %v", formats.FormatResourceString(typ, node.Usage)),
				Style:  termui.NewStyle(termui.ColorClear, termui.ColorClear),
			}
		}
	}
	// else if len(nodes) == 2 {
	// 	node, ok := r[nodes[1].Name()]
	// 	if ok {
	// 		pod, ok := node.Pods[nodes[0].Name()]
	// 		if ok {
	// 			g.UpperLimit = float64(formats.FormatResource(typ, node.Allocatable))
	// 			g.Data = append(g.Data, float64(formats.FormatResource(typ, pod.Usage)))
	// 			g.LabelData = fmt.Sprintf(labelTmpl,
	// 				formats.FormatResourceString(typ, pod.Usage),
	// 				formats.FormatResourceString(typ, node.Allocatable),
	// 				formats.FormatResourcePercentage(typ, pod.Usage, node.Allocatable),
	// 			)
	// 		}
	// 	}
	// } else if len(nodes) == 3 {
	// 	node, ok := r[nodes[2].Name()]
	// 	if ok {
	// 		pod, ok := node.Pods[nodes[1].Name()]
	// 		if ok {
	// 			container, ok := pod.Containers[nodes[0].Name()]
	// 			g.Data = append(g.Data, float64(formats.FormatResource(typ, container.Usage)))
	// 			if ok {
	// 				g.UpperLimit = float64(formats.FormatResource(typ, node.Allocatable))
	// 				g.LabelData = fmt.Sprintf(labelTmpl,
	// 					formats.FormatResourceString(typ, container.Usage),
	// 					formats.FormatResourceString(typ, node.Allocatable),
	// 					formats.FormatResourcePercentage(typ, container.Usage, node.Allocatable),
	// 				)
	// 			}
	// 		}
	// 	}
	// }
}
