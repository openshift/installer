package client

import (
	"github.com/golang/glog"

	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/labels"
)

// ListPodsWithSelector will return a list of Pods in the given namespace
// filtered by the given label selector.
func (c *Client) ListPodsWithSelector(namespace string, selector *unversioned.LabelSelector) (*v1.PodList, error) {
	sel := labels.FormatLabels(selector.MatchLabels)
	glog.V(4).Infof("[GET List Pods With Selector]: %s, selector: %s", namespace, sel)
	lo := v1.ListOptions{LabelSelector: sel}
	return c.Pods(namespace).List(lo)
}

// GetPod deletes the Pod with the given namespace and name.
func (c *Client) GetPod(namespace, name string) (*v1.Pod, error) {
	glog.V(4).Infof("[GET Pod]: %s:%s", namespace, name)
	return c.Pods(namespace).Get(name)
}

// DeletePod deletes the Pod with the given namespace and name.
func (c *Client) DeletePod(namespace, name string) error {
	glog.V(4).Infof("[DELETE Pod]: %s:%s", namespace, name)
	return c.Pods(namespace).Delete(name, nil)
}
