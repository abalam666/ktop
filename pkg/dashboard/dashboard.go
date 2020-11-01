package dashboard

import (
	"sync"

	"github.com/gizak/termui/v3"
	"github.com/ynqa/widgets"

	"github.com/ynqa/ktop/pkg/draw/graph"
	"github.com/ynqa/ktop/pkg/draw/table"
	"github.com/ynqa/ktop/pkg/resources"
	"github.com/ynqa/ktop/pkg/ui"
)

type Dashboard struct {
	mu                    sync.RWMutex
	resourceTable         *widgets.ToggleTable
	cpuGraph, memoryGraph *ui.Graph
}

func New() *Dashboard {
	return &Dashboard{
		resourceTable: newTable("RESOURCES"),
		cpuGraph:      newGraph("CPU"),
		memoryGraph:   newGraph("MEMORY"),
	}
}

func newTable(title string) *widgets.ToggleTable {
	table := widgets.NewToggleTable()
	table.Title = title
	table.TitleStyle = termui.NewStyle(termui.ColorClear)
	table.BorderStyle = termui.NewStyle(termui.ColorBlue)
	return table
}

func newGraph(title string) *ui.Graph {
	graph := ui.NewGraph()
	graph.Title = title
	graph.TitleStyle = termui.NewStyle(termui.ColorClear)
	graph.BorderStyle = termui.NewStyle(termui.Color(18))
	graph.LabelNameColor = termui.ColorWhite
	graph.DataColor = termui.ColorGreen
	graph.LimitColor = termui.ColorWhite
	return graph
}

func (d *Dashboard) ResourceTable() *widgets.ToggleTable {
	return d.resourceTable
}

func (d *Dashboard) CPUGraph() *ui.Graph {
	return d.cpuGraph
}

func (d *Dashboard) MemoryGraph() *ui.Graph {
	return d.memoryGraph
}

func (d *Dashboard) Toggle() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.resourceTable.Node.Toggle(d.resourceTable.SelectedRow)
}

func (d *Dashboard) ScrollUp() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.resourceTable.ScrollUp()
	d.cpuGraph.Reset()
	d.memoryGraph.Reset()
}

func (d *Dashboard) ScrollDown() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.resourceTable.ScrollDown()
	d.cpuGraph.Reset()
	d.memoryGraph.Reset()
}

func (d *Dashboard) UpdateTable(
	drawer table.Drawer,
	r resources.Resources,
) {
	d.mu.Lock()
	defer d.mu.Unlock()
	drawer.Draw(d.resourceTable, r)
}

func (d *Dashboard) UpdateCPUGraph(drawer graph.Drawer, c graph.Contents) {
	d.mu.Lock()
	defer d.mu.Unlock()
	// stack := d.resourceTable.Node.Flatten()
	// if 0 < len(stack) && d.resourceTable.SelectedRow < len(stack) {
	// 	drawer.Draw(d.cpuGraph, c, d.resourceTable.Rows[d.resourceTable.SelectedRow].Key)
	// }
}

func (d *Dashboard) UpdateMemoryGraph(drawer graph.Drawer, c graph.Contents) {
	d.mu.Lock()
	defer d.mu.Unlock()
	// if 0 < len(d.resourceTable.Rows) && d.resourceTable.SelectedRow < len(d.resourceTable.Rows) {
	// 	drawer.Draw(d.memoryGraph, c, d.resourceTable.Rows[d.resourceTable.SelectedRow].Key)
	// }
}
