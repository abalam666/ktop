package resources

import (
	"context"
	"regexp"
	"sort"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

type Resources map[string]*NodeResource

func (r Resources) SortedNodes() []string {
	var res []string
	for k := range r {
		res = append(res, k)
	}
	sort.Strings(res)
	return res
}

func (r Resources) SortedPods(node string) []string {
	var res []string
	for k := range r[node].Pods {
		res = append(res, k)
	}
	sort.Strings(res)
	return res
}

func (r Resources) SortedContainers(node, pod string) []string {
	var res []string
	for k := range r[node].Pods[pod].Containers {
		res = append(res, k)
	}
	sort.Strings(res)
	return res
}

// - https://kubernetes.io/docs/tasks/administer-cluster/reserve-compute-resources/#node-allocatable
// - https://github.com/kubernetes/kubectl/blob/v0.18.8/pkg/describe/describe.go#L3521-L3602
type NodeResource struct {
	Pods        map[string]*PodResource
	Capacity    corev1.ResourceList
	Allocatable corev1.ResourceList
	Usage       corev1.ResourceList
}

type PodResource struct {
	Namespace  string
	Containers map[string]*ContainerResource
	Usage      corev1.ResourceList
}

type ContainerResource struct {
	Usage corev1.ResourceList
}

func FetchResources(
	namespace string,
	clientset *kubernetes.Clientset,
	metricsclientset *versioned.Clientset,
	nodeQuery, podQuery, containerQuery *regexp.Regexp,
) (Resources, error) {

	data := make(Resources)

	// get nodes and their metrics
	nodes, err := clientset.CoreV1().Nodes().List(
		context.Background(),
		metav1.ListOptions{
			LabelSelector: labels.Everything().String(),
		},
	)
	if err != nil {
		return nil, err
	}
	nodeMetrics, err := metricsclientset.MetricsV1beta1().
		NodeMetricses().List(
		context.Background(),
		metav1.ListOptions{
			LabelSelector: labels.Everything().String(),
		},
	)
	if err != nil {
		return nil, err
	}

	nodeMetrics.Items = matchNodes(nodeQuery, nodeMetrics.Items)
	for _, metric := range nodeMetrics.Items {
		nodeStatus := getNodeStatus(metric.Name, nodes.Items)
		data[metric.Name] = &NodeResource{
			Pods:        make(map[string]*PodResource),
			Capacity:    nodeStatus.Capacity,
			Allocatable: nodeStatus.Allocatable,
			Usage:       metric.Usage.DeepCopy(),
		}
	}

	// get pods and their metrics
	pods, err := clientset.CoreV1().Pods(namespace).List(
		context.Background(),
		metav1.ListOptions{
			LabelSelector: labels.Everything().String(),
		},
	)
	if err != nil {
		return nil, err
	}
	podMetrics, err := metricsclientset.MetricsV1beta1().
		PodMetricses(namespace).List(
		context.Background(),
		metav1.ListOptions{
			LabelSelector: labels.Everything().String(),
		},
	)
	if err != nil {
		return nil, err
	}

	podMetrics.Items = matchPods(podQuery, podMetrics.Items)
	for _, podMetric := range podMetrics.Items {
		node := getAssignedNode(podMetric.Name, pods.Items)
		// it sometimes cannot get nodes because of the filters by `matchNodes`.
		if _, ok := data[node]; !ok {
			continue
		}

		// investigate all pods (without filtering) to aggregate the usages of their own containers.
		var cpuperpod, memoryperpod resource.Quantity
		for _, containerMetric := range podMetric.Containers {
			cpuperpod.Add(containerMetric.Usage.Cpu().DeepCopy())
			memoryperpod.Add(containerMetric.Usage.Memory().DeepCopy())
		}
		data[node].Pods[podMetric.Name] = &PodResource{
			Namespace:  podMetric.Namespace,
			Containers: make(map[string]*ContainerResource),
			Usage: corev1.ResourceList{
				corev1.ResourceCPU:    cpuperpod,
				corev1.ResourceMemory: memoryperpod,
			},
		}

		// but to search them by a given query is required on viewing.
		podMetric.Containers = matchContainers(containerQuery, podMetric.Containers)
		for _, containerMetric := range podMetric.Containers {
			data[node].Pods[podMetric.Name].Containers[containerMetric.Name] = &ContainerResource{
				Usage: containerMetric.Usage.DeepCopy(),
			}
		}
	}

	return data, nil
}

func matchNodes(query *regexp.Regexp, nodes []metricsv1beta1.NodeMetrics) []metricsv1beta1.NodeMetrics {
	var res []metricsv1beta1.NodeMetrics
	for _, node := range nodes {
		if query.MatchString(node.Name) {
			res = append(res, node)
		}
	}
	return res
}

func getNodeStatus(target string, nodes []corev1.Node) corev1.NodeStatus {
	for _, node := range nodes {
		if target == node.Name {
			return node.Status
		}
	}
	return corev1.NodeStatus{}
}

func matchPods(query *regexp.Regexp, pods []metricsv1beta1.PodMetrics) []metricsv1beta1.PodMetrics {
	var res []metricsv1beta1.PodMetrics
	for _, pod := range pods {
		if query.MatchString(pod.Name) {
			res = append(res, pod)
		}
	}
	return res
}

func getAssignedNode(target string, pods []corev1.Pod) string {
	for _, pod := range pods {
		if target == pod.Name {
			return pod.Spec.NodeName
		}
	}
	return ""
}

func matchContainers(query *regexp.Regexp, containers []metricsv1beta1.ContainerMetrics) []metricsv1beta1.ContainerMetrics {
	var res []metricsv1beta1.ContainerMetrics
	for _, container := range containers {
		if query.MatchString(container.Name) {
			res = append(res, container)
		}
	}
	return res
}
