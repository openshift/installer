// Package libvirt generates Machine objects for libvirt.
package libvirt

import (
	"text/template"

	"github.com/openshift/installer/pkg/types"
)

// MasterConfig is used to generate the master machine list.
type MasterConfig struct {
	ClusterName string
	Instances   []string
	Platform    types.LibvirtPlatform
}

// MasterMachinesTmpl is the template for master machines
var MasterMachinesTmpl = template.Must(template.New("master-machines").Parse(`
{{- $c := . -}}
kind: List
apiVersion: v1
metadata:
  resourceVersion: ""
  selfLink: ""
items:
{{- range $index,$instance := .Instances}}
- apiVersion: cluster.k8s.io/v1alpha1
  kind: Machine
  metadata:
    name: {{$c.ClusterName}}-master-{{$index}}
    namespace: openshift-cluster-api
    labels:
      sigs.k8s.io/cluster-api-cluster: {{$c.ClusterName}}
      sigs.k8s.io/cluster-api-machine-role: master
      sigs.k8s.io/cluster-api-machine-type: master
  spec:
    providerConfig:
      value:
        apiVersion: libvirtproviderconfig/v1alpha1
        kind: LibvirtMachineProviderConfig
        domainMemory: 2048
        domainVcpu: 2
        ignKey: /var/lib/libvirt/images/master-{{$index}}.ign
        volume:
          poolName: default
          baseVolumeID: /var/lib/libvirt/images/coreos_base
        networkInterfaceName: {{$c.Platform.Network.Name}}
        networkInterfaceAddress: {{$c.Platform.Network.IPRange}}
        autostart: false
        uri: {{$c.Platform.URI}}
    versions:
      kubelet: ""
      controlPlane: ""
{{- end }}
`))
