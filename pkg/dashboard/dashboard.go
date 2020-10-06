package dashboard

import (
	"sync"

	"github.com/gizak/termui/v3"

	"github.com/ynqa/ktop/pkg/graph"
	"github.com/ynqa/ktop/pkg/resources"
	"github.com/ynqa/ktop/pkg/table"
	"github.com/ynqa/ktop/pkg/ui"
)

type Dashboard struct {
	mu                    sync.RWMutex
	resourcetable         *ui.Table
	cpugraph, memorygraph *ui.Graph
}

func New() *Dashboard {
	return &Dashboard{
		resourcetable: newTable("Resources"),
		cpugraph:      newGraph("CPU"),
		memorygraph:   newGraph("Memory"),
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

func (d *Dashboard) ResourceTable() *ui.Table {
	return d.resourcetable
}

func (d *Dashboard) CPUGraph() *ui.Graph {
	return d.cpugraph
}

func (d *Dashboard) MemoryGraph() *ui.Graph {
	return d.memorygraph
}

func (d *Dashboard) CurrentRow() ui.Row {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.resourcetable.Rows[d.resourcetable.SelectedRow]
}

func (d *Dashboard) ScrollUp() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.resourcetable.ScrollUp()
	d.cpugraph.Reset()
	d.memorygraph.Reset()
}

func (d *Dashboard) ScrollDown() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.resourcetable.ScrollDown()
	d.cpugraph.Reset()
	d.memorygraph.Reset()
}

func (d *Dashboard) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.resourcetable.SelectedRow = 0
}

func (d *Dashboard) UpdateTable(
	shaper table.Shaper,
	r resources.Resources,
	state *table.VisibleSet,
) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.resourcetable.Header = shaper.Headers()
	d.resourcetable.ColumnWidths = shaper.Widths(d.resourcetable.Inner)
	d.resourcetable.Rows = shaper.Rows(r, state)
}

func (d *Dashboard) UpdateCPUGraph(
	shaper graph.Shaper,
	r resources.Resources,
) {
	d.mu.Lock()
	defer d.mu.Unlock()
}

func (d *Dashboard) UpdateMemoryGraph(
	shaper graph.Shaper,
	r resources.Resources,
) {
	d.mu.Lock()
	defer d.mu.Unlock()
}
