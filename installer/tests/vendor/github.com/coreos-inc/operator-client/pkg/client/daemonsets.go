package client

import (
	"fmt"
	"time"

	"github.com/golang/glog"

	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/errors"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
	v1beta1extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/pkg/util/wait"
)

const (
	daemonsetPodUpdateTimeout = 2 * time.Minute
)

// CreateDaemonSet will return the DaemonSet object for the given namespace and name.
func (c *Client) CreateDaemonSet(ds *v1beta1extensions.DaemonSet) (*v1beta1extensions.DaemonSet, error) {
	glog.V(4).Infof("[CREATE DaemonSet]: %s:%s", ds.GetNamespace(), ds.GetName())
	return c.Extensions().DaemonSets(ds.GetNamespace()).Create(ds)
}

// GetDaemonSet will return the DaemonSet object for the given namespace and name.
func (c *Client) GetDaemonSet(namespace, name string) (*v1beta1extensions.DaemonSet, error) {
	glog.V(4).Infof("[GET DaemonSet]: %s:%s", namespace, name)
	return c.Extensions().DaemonSets(namespace).Get(name)
}

// UpdateDaemonSetObject will update the given DaemonSet object.
func (c *Client) UpdateDaemonSetObject(ds *v1beta1extensions.DaemonSet) (*v1beta1extensions.DaemonSet, error) {
	glog.V(4).Infof("[UPDATE DaemonSet Object]: %s:%s", ds.GetNamespace(), ds.GetName())
	return c.Extensions().DaemonSets(ds.GetNamespace()).Update(ds)
}

// DaemonSetRollingUpdate will perform a rolling update on the given DaemonSet.
// It will first update the DaemonSet object via the API Server, ane then it will
// go through each Pod managed by the Daemonset and delete it, waiting for a new
// Pod to get scheduled and ran before moving onto the next Pod. It repeats this
// process until all Pods have been updated.
//
// TODO This method only supports image updates, not manifest updates, and will only update
// the container image in the Pod that has the same name as the DaemonSet object.
// Should make this more flexible to allow specifying container in Pod to update.
func (c *Client) DaemonSetRollingUpdate(ds *v1beta1extensions.DaemonSet, opts UpdateOpts) (bool, error) {
	glog.V(4).Infof("[ROLLING UPDATE DaemonSet]: %s:%s", ds.GetNamespace(), ds.GetName())
	if opts.Migrations != nil {
		if err := opts.Migrations.RunBeforeMigrations(c, ds.GetNamespace(), ds.GetName()); err != nil {
			return false, err
		}
	}

	unavailable, err := c.UnavailablePodsForDaemonSet(ds)
	if err != nil {
		return false, err
	}

	if unavailable >= opts.MaxUnavailable {
		return false, fmt.Errorf("unavailable pods (%d) reached maximum (%d)", unavailable, opts.MaxUnavailable)
	}

	ds, err = c.UpdateDaemonSetObject(ds)
	if err != nil {
		return false, fmt.Errorf("unable to update DaemonSet %s: %v", ds.GetName(), err)
	}

	// Updating the DaemonSet does not update pods, so here we
	// get all pods, delete them, and then wait for them to come back
	// with the correct version running.
	pl, err := c.ListPodsWithSelector(ds.GetNamespace(), ds.Spec.Selector)
	if err != nil {
		return false, err
	}
	updated := false
	for _, p := range pl.Items {
		podUpdated, err := deletePodAndWait(c, &p, ds)
		if err != nil {
			return false, err
		}
		if !updated && podUpdated {
			updated = true
		}
	}
	if opts.Migrations != nil {
		if err := opts.Migrations.RunAfterMigrations(c, ds.GetNamespace(), ds.GetName()); err != nil {
			return false, err
		}
	}
	return updated, nil
}

// AvailablePodsForDaemonSet will return the number of Pods that are ready and
// available for the given DaemonSet.
func (c *Client) AvailablePodsForDaemonSet(ds *v1beta1extensions.DaemonSet) (int, error) {
	glog.V(4).Infof("[GET Available Pods DaemonSet]: %s:%s", ds.GetNamespace(), ds.GetName())
	pl, err := c.ListPodsWithSelector(ds.GetNamespace(), ds.Spec.Selector)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, p := range pl.Items {
		var apiPod api.Pod
		if err := v1.Convert_v1_Pod_To_api_Pod(&p, &apiPod, nil); err != nil {
			return 0, err
		}
		available := api.IsPodAvailable(&apiPod, 5, unversioned.Now())
		if available {
			count++
		}
	}
	return count, nil
}

// UnavailablePodsForDaemonSet returns the count of unavailable Pods for the given
// DaemonSet object.
func (c *Client) UnavailablePodsForDaemonSet(ds *v1beta1extensions.DaemonSet) (int, error) {
	glog.V(4).Infof("[GET Available Pods DaemonSet]: %s:%s", ds.GetNamespace(), ds.GetName())
	pl, err := c.ListPodsWithSelector(ds.GetNamespace(), ds.Spec.Selector)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, p := range pl.Items {
		var apiPod api.Pod
		if err := v1.Convert_v1_Pod_To_api_Pod(&p, &apiPod, nil); err != nil {
			return 0, err
		}
		available := api.IsPodAvailable(&apiPod, 5, unversioned.Now())
		if !available {
			count++
		}
	}
	return count, nil
}

// NumberOfDesiredPodsForDaemonSet returns the number of Pods the DaemonSet should
// be running.
func (c *Client) NumberOfDesiredPodsForDaemonSet(ds *v1beta1extensions.DaemonSet) (int, error) {
	glog.V(4).Infof("[GET Available Pods DaemonSet]: %s:%s", ds.GetNamespace(), ds.GetName())
	ds, err := c.GetDaemonSet(ds.GetNamespace(), ds.GetName())
	if err != nil {
		return 0, err
	}
	return int(ds.Status.DesiredNumberScheduled), nil
}

// deletePodAndWait will delete the given Pod and wait until another Pod has been
// scheduled in its place. It will return true if the Pod is deleted and updated.
func deletePodAndWait(client *Client, p *v1.Pod, ds *v1beta1extensions.DaemonSet) (bool, error) {
	// Delete old DS Pod.
	glog.V(4).Infof("Deleting pod %s", p.Name)
	err := client.DeletePod(p.GetNamespace(), p.GetName())
	if err != nil {
		return false, err
	}
	glog.V(4).Infof("Deleted pod %s", p.GetName())

	// Wait for all pods to be available before moving on.
	if err = wait.Poll(time.Second, daemonsetPodUpdateTimeout, func() (bool, error) {
		// Wait for Pod to no longer exist.
		_, err = client.GetPod(p.GetNamespace(), p.GetName())
		if !errors.IsNotFound(err) {
			return false, nil
		}
		var available int
		available, err = client.AvailablePodsForDaemonSet(ds)
		if err != nil {
			return false, err
		}
		var desired int
		desired, err = client.NumberOfDesiredPodsForDaemonSet(ds)
		if err != nil {
			return false, err
		}
		if available == desired {
			return true, nil
		}
		glog.V(4).Infof("waiting for new DaemonSet %s pods to be available: %d of %d", ds.GetName(), available, desired)
		return false, nil
	}); err != nil {
		return false, fmt.Errorf("error waiting for DaemonSet %s Pod update: %v", ds.GetName(), err)
	}
	return true, nil
}
