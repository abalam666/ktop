package dashboard

import (
	"sync"

	"github.com/gizak/termui/v3"

	"github.com/ynqa/ktop/pkg/resources"
	"github.com/ynqa/ktop/pkg/table"
	"github.com/ynqa/ktop/pkg/ui"
)

type Dashboard struct {
	mu                    sync.RWMutex
	ResourceTable         *ui.Table
	CPUGraph, MemoryGraph *ui.Graph
}

func New() *Dashboard {
	return &Dashboard{
		ResourceTable: newTable("Resources"),
		CPUGraph:      newGraph("CPU"),
		MemoryGraph:   newGraph("Memory"),
	}
}

func newTable(title string) *ui.Table {
	table := ui.NewTable()
	table.Title = title
	table.TitleStyle = termui.NewStyle(termui.ColorWhite, termui.ColorClear, termui.ModifierBold)
	table.Cursor = true
	table.BorderStyle = termui.NewStyle(termui.ColorBlue)
	table.CursorColor = termui.ColorYellow
	return table
}

func newGraph(title string) *ui.Graph {
	graph := ui.NewGraph()
	graph.Title = title
	graph.TitleStyle = termui.NewStyle(termui.ColorWhite, termui.ColorClear, termui.ModifierBold)
	graph.BorderStyle = termui.NewStyle(termui.ColorBlue)
	graph.LabelNameColor = termui.ColorWhite
	graph.DataColor = termui.ColorGreen
	graph.LimitColor = termui.ColorWhite
	return graph
}

func (d *Dashboard) ScrollUp() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.ResourceTable.ScrollUp()
	d.CPUGraph.Reset()
	d.MemoryGraph.Reset()
}

func (d *Dashboard) ScrollDown() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.ResourceTable.ScrollDown()
	d.CPUGraph.Reset()
	d.MemoryGraph.Reset()
}

func (d *Dashboard) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.ResourceTable.SelectedRow = 0
}

func (d *Dashboard) UpdateTable(
	shaper table.Shaper,
	r resources.Resources,
	state *table.VisibleSet,
) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.ResourceTable.Header = shaper.Headers()
	d.ResourceTable.ColumnWidths = shaper.Widths(d.ResourceTable.Inner)
	d.ResourceTable.Rows = shaper.Rows(r, state)
}
