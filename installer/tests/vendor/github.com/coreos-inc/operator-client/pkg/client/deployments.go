package client

import (
	"time"

	"github.com/golang/glog"

	v1beta1extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/pkg/util/wait"
)

const (
	deploymentUpdateTimeout = 5 * time.Minute
)

// GetDeployment will return the Deployment object specified by the given
// namespace and name.
func (c *Client) GetDeployment(namespace, name string) (*v1beta1extensions.Deployment, error) {
	glog.V(4).Infof("[GET Deployment]: %s:%s", namespace, name)
	return c.Extensions().Deployments(namespace).Get(name)
}

// CreateDeployment will create the given Deployment object.
func (c *Client) CreateDeployment(dep *v1beta1extensions.Deployment) (*v1beta1extensions.Deployment, error) {
	glog.V(4).Infof("[CREATE Deployment]: %s:%s", dep.GetNamespace(), dep.GetName())
	return c.Extensions().Deployments(dep.GetNamespace()).Create(dep)
}

// UpdateDeployment will update the given Deployment object.
// Note, this does not wait for the Deployment rollout, only updates
// the Deployment object via the API Server.
func (c *Client) UpdateDeployment(dep *v1beta1extensions.Deployment) (*v1beta1extensions.Deployment, error) {
	glog.V(4).Infof("[UPDATE Deployment]: %s:%s", dep.GetNamespace(), dep.GetName())
	return c.Extensions().Deployments(dep.GetNamespace()).Update(dep)
}

// DeploymentRollingUpdate will perform a rolling update on the given Deployment.
// The return values are whether the Deployment was updated, and an error.
func (c *Client) DeploymentRollingUpdate(dep *v1beta1extensions.Deployment, opts UpdateOpts) (bool, error) {
	glog.V(4).Infof("[ROLLING UPDATE Deployment]: %s:%s", dep.GetNamespace(), dep.GetName())
	if opts.Migrations != nil {
		if err := opts.Migrations.RunBeforeMigrations(c, dep.GetNamespace(), dep.GetName()); err != nil {
			return false, err
		}
	}

	var err error
	oldGeneration := dep.Status.ObservedGeneration
	oldResourceVersion := dep.GetResourceVersion()
	// Update the image in the specific container we care about (should match name
	// of the deployment itself, per convention).
	dep, err = c.UpdateDeployment(dep)
	if err != nil {
		return false, err
	}
	updated := dep.GetResourceVersion() != oldResourceVersion
	if !updated {
		return false, nil
	}
	err = c.waitForDeploymentRollout(dep, oldGeneration)
	if err != nil {
		return false, err
	}
	if opts.Migrations != nil {
		if err := opts.Migrations.RunAfterMigrations(c, dep.GetNamespace(), dep.GetName()); err != nil {
			return false, err
		}
	}
	return updated, nil
}

// waitForDeploymentRollout will wait until the Deployment has finshed
// a rolling update. It does this by getting the new ReplicaSet for the
// Deployment, and waiting until all of its Pods are available.
func (c *Client) waitForDeploymentRollout(dep *v1beta1extensions.Deployment, oldGeneration int64) error {
	return wait.Poll(time.Second, deploymentUpdateTimeout, func() (bool, error) {
		d, err := c.GetDeployment(dep.GetNamespace(), dep.GetName())
		if err != nil {
			return false, err
		}
		if d.Status.ObservedGeneration > oldGeneration &&
			d.Status.UpdatedReplicas == d.Status.Replicas &&
			d.Status.UnavailableReplicas == 0 {
			return true, nil
		}
		return false, nil
	})
}
