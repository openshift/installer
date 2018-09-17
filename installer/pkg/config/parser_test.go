package config

import (
	"testing"

	"github.com/coreos/tectonic-config/config/tectonic-network"
	"github.com/openshift/installer/installer/pkg/config/aws"
	"github.com/openshift/installer/installer/pkg/config/libvirt"
	"github.com/stretchr/testify/assert"
)

func TestParseInstallConfig(t *testing.T) {
	data := []byte(`admin:
  email: test-email
  password: test-password
  sshKey: test-sshkey
baseDomain: test-domain
clusterID: test-cluster-id
machines:
  - name: master
    replicas: 3
  - name: worker
    replicas: 2
    platform:
      aws:
        type: m4-large        
metadata:
  name: test-cluster-name
networking:
  podCIDR: 10.2.0.0/16
  serviceCIDR: 10.3.0.0/16
  type: flannel
pullSecret: '{"auths": {}}'
`)

	actual, err := ParseConfig(data)
	if err != nil {
		t.Fatal(err)
	}
	actual.EC2AMIOverride = ""

	expected := &Cluster{
		Name: "test-cluster-name",
		Admin: Admin{
			Email:    "test-email",
			Password: "test-password",
			SSHKey:   "test-sshkey",
		},
		BaseDomain: "test-domain",
		CA: CA{
			RootCAKeyAlg: "RSA",
		},
		Internal: Internal{
			ClusterID: "test-cluster-id",
		},
		Networking: Networking{
			Type:        tectonicnetwork.NetworkFlannel,
			PodCIDR:     "10.2.0.0/16",
			ServiceCIDR: "10.3.0.0/16",
		},
		NodePools: []NodePool{
			{
				Name:  "master",
				Count: 3,
			},
			{
				Name:  "worker",
				Count: 2,
			},
		},
		AWS: aws.AWS{
			Endpoints:    aws.EndpointsAll,
			Region:       aws.DefaultRegion,
			Profile:      aws.DefaultProfile,
			VPCCIDRBlock: "10.0.0.0/16",
			Worker: aws.Worker{
				EC2Type: "m4-large",
			},
		},
		Libvirt: libvirt.Libvirt{
			Network: libvirt.Network{
				IfName: libvirt.DefaultIfName,
			},
		},
		PullSecret: "{\"auths\": {}}",
	}

	assert.Equal(t, expected, actual)
}
