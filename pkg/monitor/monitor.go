package monitor

import (
	"errors"
	"regexp"

	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/gizak/termui/v3"

	"github.com/ynqa/ktop/pkg/resources"
	"github.com/ynqa/ktop/pkg/ui"
	"github.com/ynqa/ktop/pkg/viewers"
)

type Monitor struct {
	namespace                           string
	clientset                           *kubernetes.Clientset
	metricsclientset                    *versioned.Clientset
	podQuery, containerQuery, nodeQuery *regexp.Regexp

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
		clientset:        clientset,
		metricsclientset: metricsclientset,
		podQuery:         podQuery,
		containerQuery:   containerQuery,
		nodeQuery:        nodeQuery,

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

		var viewer viewers.Table
		if len(data) > 0 {
			viewer = &viewers.ResourceTable{}
		} else {
			viewer = &viewers.EmptyTable{}
		}
		fields := viewer.Fields(m.ResourceTable.Inner, data)
		m.ResourceTable.Header = fields.Headers
		m.ResourceTable.ColumnWidths = fields.Widths
		m.ResourceTable.Rows = fields.Rows
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
