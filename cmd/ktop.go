package cmd

import (
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"

	"github.com/gizak/termui/v3"
	"github.com/spf13/cobra"
	"github.com/ynqa/ktop/pkg/monitor"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/metrics/pkg/client/clientset/versioned"
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

func (k *ktop) run(cmd *cobra.Command, args []string) error {
	// setup kubernetes clients
	clientset, metricsclient, err := k.kubeclient()
	if err != nil {
		return err
	}

	// setup queries
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

	// setup monitor
	monitor := monitor.New(
		*k.kubeFlags.Namespace,
		clientset,
		metricsclient,
		podQuery,
		containerQuery,
		nodeQuery,
	)

	// start termui
	if err := termui.Init(); err != nil {
		return err
	}
	defer termui.Close()

	// setup grid
	grid := termui.NewGrid()
	grid.Set(
		termui.NewRow(1/2, monitor.ResourceTable),
		termui.NewRow(1/4, monitor.CPUGraph),
		termui.NewRow(1/4, monitor.MemoryGraph),
	)
	width, height := termui.TerminalDimensions()
	grid.SetRect(0, 0, width, height)

	// poll events
	event := termui.PollEvents()
	tick := time.NewTicker(k.interval)
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, os.Interrupt)

	for {
		select {
		case <-sig:
			return nil
		case <-tick.C:
			if err := monitor.Sync(); err != nil {
				return err
			}

		case e := <-event:
			switch e.ID {
			case "q", "<C-c>":
				return nil
			case "<Resize>":
				width, height := termui.TerminalDimensions()
				grid.SetRect(0, 0, width, height)
			}
		}
		// k.mu.Lock()
		// termui.Render(grid)
		// k.mu.Unlock()
	}
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
