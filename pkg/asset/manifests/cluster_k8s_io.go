package manifests

import (
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	clusterv1a1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

// This file was originally in pkg/assets/machines, but is now in
// /manifests due to an import loop.

// ClusterK8sIO generates the `Cluster.machine.openshift.io/v1beta1` object.
type ClusterK8sIO struct {
	Raw []byte
}

var _ asset.Asset = (*ClusterK8sIO)(nil)

// Name returns a human friendly name for the ClusterK8sIO Asset.
func (c *ClusterK8sIO) Name() string {
	return "Cluster.machine.openshift.io/v1beta1"
}

// Dependencies returns all of the dependencies directly needed by the
// ClusterK8sIO asset
func (c *ClusterK8sIO) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.ClusterID{},
		&Networking{},
	}
}

// Generate generates the Worker asset.
func (c *ClusterK8sIO) Generate(dependencies asset.Parents) error {
	clusterID := &installconfig.ClusterID{}
	net := &Networking{}
	dependencies.Get(clusterID, net)
	clusterNet, err := net.ClusterNetwork()
	if err != nil {
		return errors.Wrap(err, "failed to generate ClusterNetworkingConfig")
	}

	cluster := clusterv1a1.Cluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machine.openshift.io/v1beta1",
			Kind:       "Cluster",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID,
			Namespace: "openshift-machine-api",
		},
		Spec: clusterv1a1.ClusterSpec{
			ClusterNetwork: *clusterNet,
		},
	}

	c.Raw, err = yaml.Marshal(cluster)
	return err
}
