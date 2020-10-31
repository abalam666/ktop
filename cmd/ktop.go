package cmd

import (
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/gizak/termui/v3"
	"github.com/spf13/cobra"

	"github.com/ynqa/ktop/pkg/dashboard"
	"github.com/ynqa/ktop/pkg/draw/graph"
	"github.com/ynqa/ktop/pkg/draw/table"
	"github.com/ynqa/ktop/pkg/resources"
)

type ktop struct {
	mu sync.RWMutex

	interval       time.Duration
	nodeQuery      string
	podQuery       string
	containerQuery string

	kubeFlags *genericclioptions.ConfigFlags
}

func New() *cobra.Command {
	ktop := ktop{}
	cmd := &cobra.Command{
		Use:   "ktop",
		Short: "Kubernetes monitoring dashboard on terminal",
		RunE:  ktop.run,
	}
	cmd.Flags().DurationVarP(
		&ktop.interval,
		"interval",
		"i",
		1*time.Second,
		"set interval",
	)
	cmd.Flags().StringVarP(
		&ktop.nodeQuery,
		"node-query",
		"N",
		".*",
		"node query",
	)
	cmd.Flags().StringVarP(
		&ktop.podQuery,
		"pod-query",
		"P",
		".*",
		"pod query",
	)
	cmd.Flags().StringVarP(
		&ktop.containerQuery,
		"container-query",
		"C",
		".*",
		"container query",
	)

	ktop.kubeFlags = genericclioptions.NewConfigFlags(false)
	ktop.kubeFlags.AddFlags(cmd.Flags())
	if *ktop.kubeFlags.Namespace == "" {
		*ktop.kubeFlags.Namespace = "default"
	}

	return cmd
}

func (k *ktop) loop(
	clientset *kubernetes.Clientset,
	metricsclientset *versioned.Clientset,
	podQuery, containerQuery, nodeQuery *regexp.Regexp,
) error {
	// start termui
	if err := termui.Init(); err != nil {
		return err
	}
	defer termui.Close()

	// draw grid
	dashboard := dashboard.New()
	grid := termui.NewGrid()
	grid.Set(
		termui.NewRow(1./2, dashboard.ResourceTable()),
		termui.NewRow(1./4, dashboard.CPUGraph()),
		termui.NewRow(1./4, dashboard.MemoryGraph()),
	)

	panel := func() {
		width, height := termui.TerminalDimensions()
		grid.SetRect(0, 1, width, height-1)
	}
	panel()

	// rendering
	render := func() {
		k.mu.Lock()
		termui.Render(grid)
		k.mu.Unlock()
	}

	errCh := make(chan error)

	tick := time.NewTicker(k.interval)
	recv := make(chan resources.Resources)

	// scheduled to fetch resources from kubernetes metrics server.
	go func() {
		for {
			select {
			case <-tick.C:
				r, err := resources.FetchResources(
					*k.kubeFlags.Namespace,
					clientset,
					metricsclientset,
					podQuery,
					containerQuery,
					nodeQuery,
				)
				if err != nil {
					errCh <- err
					return
				}
				recv <- r
			}
		}
	}()

	tableState := table.NewVisibleSet()
	event := termui.PollEvents()
	doneCh := make(chan struct{})

	go func() {
		for r := range recv {
			// update table:
			go func(r resources.Resources) {
				var drawer table.Drawer
				if len(r) > 0 {
					drawer = &table.KubeDrawer{}
				} else {
					drawer = &table.NopDrawer{}
				}
				dashboard.UpdateTable(drawer, r, tableState)
				render()
			}(r)

			contents := graph.NewForResources(r)
			var drawer graph.Drawer
			if contents.Len() > 0 {
				drawer = &graph.KubeDrawer{}
			} else {
				drawer = &graph.NopDrawer{}
			}

			// update cpu graph:
			go func(drawer graph.Drawer, c graph.Contents) {
				dashboard.UpdateCPUGraph(drawer, c)
				render()
			}(drawer, contents)

			// update memory graph:
			go func(drawer graph.Drawer, c graph.Contents) {
				dashboard.UpdateMemoryGraph(drawer, c)
				render()
			}(drawer, contents)
		}
	}()

	go func() {
		for e := range event {
			switch e.ID {
			case "<Enter>":
				tableState.Toggle(dashboard.CurrentRow().Key)
			case "<Down>":
				dashboard.ScrollDown()
			case "<Up>":
				dashboard.ScrollUp()
			case "q", "<C-c>":
				doneCh <- struct{}{}
				return
			case "<Resize>":
				panel()
			}
			render()
		}
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, os.Interrupt)

	for {
		defer func() {
			close(sig)
			close(recv)
			close(errCh)
			close(doneCh)
		}()
		select {
		case <-sig:
			return nil
		case <-doneCh:
			return nil
		case err := <-errCh:
			if err != nil {
				return err
			}
		}
	}
}

func (k *ktop) run(cmd *cobra.Command, args []string) error {
	// kubernetes clients
	clientset, metricsclientset, err := k.kubeclient()
	if err != nil {
		return err
	}

	// regexp queries
	podQuery, err := regexp.Compile(k.podQuery)
	if err != nil {
		return err
	}
	containerQuery, err := regexp.Compile(k.containerQuery)
	if err != nil {
		return err
	}
	nodeQuery, err := regexp.Compile(k.nodeQuery)
	if err != nil {
		return err
	}

	return k.loop(
		clientset,
		metricsclientset,
		podQuery,
		containerQuery,
		nodeQuery,
	)
}

func (k *ktop) kubeclient() (*kubernetes.Clientset, *versioned.Clientset, error) {
	config, err := k.kubeFlags.ToRESTConfig()
	if err != nil {
		return nil, nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	metricsclientset, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	return clientset, metricsclientset, nil
}
