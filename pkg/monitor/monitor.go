package monitor

import (
	"errors"
	"regexp"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/gizak/termui/v3"

	"github.com/ynqa/ktop/pkg/resources"
	"github.com/ynqa/ktop/pkg/table"
	"github.com/ynqa/ktop/pkg/table/formats"
	"github.com/ynqa/ktop/pkg/table/state"
	"github.com/ynqa/ktop/pkg/ui"
)

type Monitor struct {
	mu sync.RWMutex

	namespace                           string
	clientset                           *kubernetes.Clientset
	metricsclientset                    *versioned.Clientset
	podQuery, containerQuery, nodeQuery *regexp.Regexp

	childVisibleSet       *state.ChildVisibleSet
	ResourceTable         *ui.Table
	CPUGraph, MemoryGraph *ui.Graph
}

func New(
	namespace string,
	clientset *kubernetes.Clientset,
	metricsclientset *versioned.Clientset,
	podQuery, containerQuery, nodeQuery *regexp.Regexp,
) *Monitor {
	return &Monitor{
		namespace:        namespace,
		clientset:        clientset,
		metricsclientset: metricsclientset,
		podQuery:         podQuery,
		containerQuery:   containerQuery,
		nodeQuery:        nodeQuery,

		childVisibleSet: state.New(),
		ResourceTable:   newTable("Resources"),
		CPUGraph:        newGraph("CPU"),
		MemoryGraph:     newGraph("Memory"),
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

func (m *Monitor) Sync() error {
	errCh := make(chan error)
	doneCh := make(chan struct{})
	dataCh := make(chan resources.Resources)

	go func() {
		resources, err := resources.FetchResources(
			m.namespace,
			m.clientset,
			m.metricsclientset,
			m.nodeQuery, m.podQuery, m.containerQuery,
		)
		if err != nil {
			errCh <- err
			return
		}
		dataCh <- resources
	}()

	go func() {
		data, ok := <-dataCh
		if !ok {
			errCh <- errors.New("failed to get resources")
			return
		}

		// change ui states
		m.mu.Lock()
		defer m.mu.Unlock()

		var creator table.ContentsCreator
		if len(data) > 0 {
			creator = &table.KubeResourceContentsCreator{}
		} else {
			creator = &table.NoContentsCreator{}
		}
		contents := creator.Create(data, m.childVisibleSet, m.ResourceTable.Inner)
		m.ResourceTable.Header = contents.Headers
		m.ResourceTable.ColumnWidths = contents.Widths
		m.ResourceTable.Rows = contents.Rows
		doneCh <- struct{}{}
	}()

	for {
		select {
		case <-doneCh:
			return nil
		case err := <-errCh:
			return err
		}
	}
}

func (m *Monitor) ScrollUp() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ResourceTable.ScrollUp()
	m.CPUGraph.Reset()
	m.MemoryGraph.Reset()
}

func (m *Monitor) ScrollDown() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ResourceTable.ScrollDown()
	m.CPUGraph.Reset()
	m.MemoryGraph.Reset()
}

func (m *Monitor) Enter() {
	m.mu.Lock()
	defer m.mu.Unlock()
	name := m.ResourceTable.Rows[m.ResourceTable.SelectedRow][0]
	m.childVisibleSet.Switch(formats.TrimString(name))
}
