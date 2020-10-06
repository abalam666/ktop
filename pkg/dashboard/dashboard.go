package dashboard

import (
	"sync"

	"github.com/gizak/termui/v3"

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

// func (d *Dashboard) Sync() error {
// 	errCh := make(chan error)
// 	doneCh := make(chan struct{})
// 	dataCh := make(chan resources.Resources)

// 	go func() {
// 		resources, err := resources.FetchResources(
// 			m.namespace,
// 			m.clientset,
// 			m.metricsclientset,
// 			m.nodeQuery, m.podQuery, m.containerQuery,
// 		)
// 		if err != nil {
// 			errCh <- err
// 			return
// 		}
// 		dataCh <- resources
// 	}()

// 	go func() {
// 		data, ok := <-dataCh
// 		if !ok {
// 			errCh <- errors.New("failed to get resources")
// 			return
// 		}

// 		// change ui from here!
// 		m.mu.Lock()
// 		defer m.mu.Unlock()

// 		// change ui: table section
// 		var creator table.ContentsCreator
// 		if len(data) > 0 {
// 			creator = &table.KubeResourceContentsCreator{}
// 		} else {
// 			creator = &table.NoContentsCreator{}
// 		}
// 		contents := creator.Create(data, m.childVisibleSet, m.ResourceTable.Inner)
// 		m.ResourceTable.Header = contents.Headers
// 		m.ResourceTable.ColumnWidths = contents.Widths
// 		m.ResourceTable.Rows = contents.Rows

// 		// TODO: change ui: graph sections
// 		doneCh <- struct{}{}
// 	}()

// 	for {
// 		select {
// 		case <-doneCh:
// 			return nil
// 		case err := <-errCh:
// 			return err
// 		}
// 	}
// }

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

func (d *Dashboard) SwitchChildVisible() {
	d.mu.Lock()
	defer d.mu.Unlock()
	// name := d.ResourceTable.Rows[d.ResourceTable.SelectedRow][0]
	// d.childVisibleSet.Switch(formats.TrimString(name))
}

func (d *Dashboard) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	// d.childVisibleSet.Reset()
	// d.ResourceTable.SelectedRow = 0
}
