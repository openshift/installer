package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang/glog"

	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/fields"
	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/pkg/util/wait"

	k8stypes "k8s.io/client-go/pkg/kubelet/types"
)

const (
	podEvictionTimeout = 15 * time.Second
)

// GetNode will return the Node object specified by the given name.
func (c *Client) GetNode(name string) (*v1.Node, error) {
	glog.V(4).Infof("[GET Node]: %s", name)
	return c.Core().Nodes().Get(name)
}

// UpdateNode will update the Node object given.
func (c *Client) UpdateNode(node *v1.Node) (*v1.Node, error) {
	glog.V(4).Infof("[UPDATE Node]: %s", node.GetName())
	return c.Core().Nodes().Update(node)
}

// CordonNode will mark the Node 'n' as unschedulable.
func (c *Client) CordonNode(n *v1.Node) (*v1.Node, error) {
	glog.V(4).Infof("[CORDON Node]: %s", n.GetName())
	n.Spec.Unschedulable = true
	nn, err := c.Core().Nodes().Update(n)
	if err != nil {
		return nil, fmt.Errorf("error changing schedulable status of Node %s to %v: %v", n.Name, true, err)
	}
	return nn, nil
}

// UnCordonNode will mark the Node 'n' as schedulable.
func (c *Client) UnCordonNode(n *v1.Node) (*v1.Node, error) {
	glog.V(4).Infof("[UNCORDON Node]: %s", n.GetName())
	n.Spec.Unschedulable = false
	nn, err := c.Core().Nodes().Update(n)
	if err != nil {
		return nil, fmt.Errorf("error changing schedulable status of Node %s to %v: %v", n.Name, false, err)
	}
	return nn, nil
}

// DrainNode will drain all Pods from the given Node, with the exception of
// DaemonSet managed Pods.
func (c *Client) DrainNode(n *v1.Node) error {
	glog.V(4).Infof("[DRAIN Node]: %s", n.GetName())
	pods, err := getPodsForDeletion(c, n.GetName())
	if err != nil {
		return fmt.Errorf("error getting pods for deletion when attempting to drain Node %s: %v", n.GetName(), err)
	}
	for _, p := range pods {
		if err := evictPod(c, p, n); err != nil {
			return err
		}
	}
	return nil
}

func evictPod(c *Client, p v1.Pod, n *v1.Node) error {
	evictionBytes := []byte(`{
  "apiVersion": "policy/v1beta1",
  "kind": "Eviction",
  "metadata": {
    "name": "` + p.Name + `",
    "namespace": "` + p.Namespace + `"
  }
}
`)
	err := wait.Poll(time.Second, podEvictionTimeout, func() (bool, error) {
		res := c.Extensions().
			RESTClient().
			Post().
			AbsPath("api", "v1", "namespaces", p.Namespace, string(v1.ResourcePods), p.Name, "eviction").
			Body(evictionBytes).
			Do()
		var status int
		if res.StatusCode(&status); status == http.StatusTooManyRequests {
			glog.Infof("unable to evict pod %s due to disruption budget", p.Name)
			return false, nil
		}
		if res.Error() != nil {
			return false, fmt.Errorf("error evicting Pod %s for drain on Node %s: %v", p.Name, n.Name, res.Error())
		}
		return true, nil
	})
	if err != nil {
		return fmt.Errorf("error trying to evict Pod %s: %v", p.Name, err)
	}
	return nil
}

func getPodsForDeletion(c *Client, node string) (pods []v1.Pod, err error) {
	pi := c.Core().Pods(api.NamespaceAll)
	podList, err := pi.List(v1.ListOptions{
		FieldSelector: fields.SelectorFromSet(fields.Set{"spec.nodeName": node}).String(),
	})
	if err != nil {
		return pods, err
	}

	for _, pod := range podList.Items {
		// skip mirror pods
		if _, ok := pod.Annotations[k8stypes.ConfigMirrorAnnotationKey]; ok {
			continue
		}

		// unlike kubelet we don't care if you have emptyDir volumes or
		// are not replicated via some controller. sorry.

		// but we do skip daemonset pods, since ds controller will just restart them
		if creatorRef, ok := pod.Annotations[api.CreatedByAnnotation]; ok {
			// decode ref to find kind
			sr := &api.SerializedReference{}
			if err := runtime.DecodeInto(api.Codecs.UniversalDecoder(), []byte(creatorRef), sr); err != nil {
				// really shouldn't happen but at least complain verbosely if it does
				return nil, fmt.Errorf("failed decoding %q annotation on pod %q: %v", api.CreatedByAnnotation, pod.Name, err)
			}
			if sr.Reference.Kind == "DaemonSet" {
				continue
			}

		}

		pods = append(pods, pod)
	}

	return pods, nil
}

func getController(c *Client, sr *api.SerializedReference) (interface{}, error) {
	switch sr.Reference.Kind {
	case "ReplicationController":
		return c.Core().ReplicationControllers(sr.Reference.Namespace).Get(sr.Reference.Name)
	case "DaemonSet":
		return c.Extensions().DaemonSets(sr.Reference.Namespace).Get(sr.Reference.Name)
	case "Job":
		return c.Batch().Jobs(sr.Reference.Namespace).Get(sr.Reference.Name)
	case "ReplicaSet":
		return c.Extensions().ReplicaSets(sr.Reference.Namespace).Get(sr.Reference.Name)
	}
	return nil, fmt.Errorf("unknown controller kind %q", sr.Reference.Kind)
}
