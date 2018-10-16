package machines

import (
	"bytes"
	"text/template"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
)

// ClusterK8sIO generates the `Cluster.cluster.k8s.io/v1alpha1` object.
type ClusterK8sIO struct {
	Raw []byte
}

var _ asset.Asset = (*ClusterK8sIO)(nil)

// Name returns a human friendly name for the ClusterK8sIO Asset.
func (c *ClusterK8sIO) Name() string {
	return "Cluster.cluster.k8s.io/v1alpha1"
}

// Dependencies returns all of the dependencies directly needed by the
// ClusterK8sIO asset
func (c *ClusterK8sIO) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
	}
}

// Generate generates the Worker asset.
func (c *ClusterK8sIO) Generate(dependencies asset.Parents) error {
	installconfig := &installconfig.InstallConfig{}
	dependencies.Get(installconfig)

	c.Raw = clusterK8sIO(installconfig.Config)
	return nil
}

var clusterK8sIOTmpl = template.Must(template.New("cluster").Parse(`
apiVersion: "cluster.k8s.io/v1alpha1"
kind: Cluster
metadata:
  name: {{.Name}}
  namespace: openshift-cluster-api
spec:
  clusterNetwork:
    services:
      cidrBlocks:
      - {{.ServiceCIDR}}
    pods:
      cidrBlocks:
      - {{.PodCIDR}}
    serviceDomain: unused
`))

func clusterK8sIO(ic *types.InstallConfig) []byte {
	templateData := struct {
		Name        string
		ServiceCIDR string
		PodCIDR     string
	}{
		Name:        ic.ObjectMeta.Name,
		ServiceCIDR: ic.Networking.ServiceCIDR.String(),
		PodCIDR:     ic.Networking.PodCIDR.String(),
	}
	buf := &bytes.Buffer{}
	if err := clusterK8sIOTmpl.Execute(buf, templateData); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
