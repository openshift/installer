// Package libvirt generates Machine objects for libvirt.
package libvirt

import (
	"text/template"

	"github.com/openshift/installer/pkg/types"
)

// Config is used to generate the machine.
type Config struct {
	ClusterName string
	Replicas    int64
	Platform    types.LibvirtPlatform
}

// WorkerMachineSetTmpl is template for worker machineset.
var WorkerMachineSetTmpl = template.Must(template.New("worker-machineset").Parse(`
apiVersion: cluster.k8s.io/v1alpha1
kind: MachineSet
metadata:
  name: {{.ClusterName}}-worker-0
  namespace: openshift-cluster-api
  labels:
    sigs.k8s.io/cluster-api-cluster: {{.ClusterName}}
    sigs.k8s.io/cluster-api-machine-role: worker
    sigs.k8s.io/cluster-api-machine-type: worker
spec:
  replicas: {{.Replicas}}
  selector:
    matchLabels:
      sigs.k8s.io/cluster-api-machineset: worker
      sigs.k8s.io/cluster-api-cluster: {{.ClusterName}}
      sigs.k8s.io/cluster-api-machine-role: worker
      sigs.k8s.io/cluster-api-machine-type: worker
  template:
    metadata:
      labels:
        sigs.k8s.io/cluster-api-machineset: worker
        sigs.k8s.io/cluster-api-cluster: {{.ClusterName}}
        sigs.k8s.io/cluster-api-machine-role: worker
        sigs.k8s.io/cluster-api-machine-type: worker
    spec:
      providerConfig:
        value:
          apiVersion: libvirtproviderconfig/v1alpha1
          kind: LibvirtMachineProviderConfig
          domainMemory: 2048
          domainVcpu: 2
          ignKey: /var/lib/libvirt/images/worker.ign
          volume:
            poolName: default
            baseVolumeID: /var/lib/libvirt/images/coreos_base
          networkInterfaceName: {{.Platform.Network.Name}}
          networkInterfaceAddress: {{.Platform.Network.IPRange}}
          autostart: false
          uri: {{.Platform.URI}}
      versions:
        kubelet: ""
        controlPlane: ""
`))
