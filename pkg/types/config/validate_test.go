package config

import (
	"testing"

	"github.com/openshift/installer/pkg/types/config/aws"
	"github.com/openshift/installer/pkg/types/config/libvirt"
)

func TestMissingNodePool(t *testing.T) {
	cases := []struct {
		cluster Cluster
		errs    int
	}{
		{
			cluster: Cluster{},
			errs:    2,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"", "", ""},
				},
			},
			errs: 2,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master"},
				},
			},
			errs: 1,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master"},
				},
				Worker: Worker{
					NodePools: []string{"worker"},
				},
			},
			errs: 0,
		},
	}

	for i, c := range cases {
		var n int
		errs := c.cluster.Validate()
		for _, err := range errs {
			if _, ok := err.(*ErrMissingNodePool); ok {
				n++
			}
		}

		if n != c.errs {
			t.Errorf("test case %d: expected %d missing node pool errors, got %d", i, c.errs, n)
		}
	}
}

func TestMoreThanOneNodePool(t *testing.T) {
	cases := []struct {
		cluster Cluster
		errs    int
	}{
		{
			cluster: Cluster{},
			errs:    0,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master"},
				},
			},
			errs: 0,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master"},
				},
				Worker: Worker{
					NodePools: []string{"worker"},
				},
			},
			errs: 0,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master", "master2"},
				},
				Worker: Worker{
					NodePools: []string{"worker"},
				},
			},
			errs: 1,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master", "master2"},
				},
				Worker: Worker{
					NodePools: []string{"worker", "worker2"},
				},
			},
			errs: 2,
		},
	}

	for i, c := range cases {
		var n int
		errs := c.cluster.Validate()
		for _, err := range errs {
			if _, ok := err.(*ErrMoreThanOneNodePool); ok {
				n++
			}
		}

		if n != c.errs {
			t.Errorf("test case %d: expected %d more-than-one node pool errors, got %d", i, c.errs, n)
		}
	}
}

func TestUnmatchedNodePool(t *testing.T) {
	cases := []struct {
		cluster Cluster
		errs    int
	}{
		{
			cluster: Cluster{},
			errs:    0,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master"},
				},
			},
			errs: 1,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master"},
				},
				Worker: Worker{
					NodePools: []string{"worker"},
				},
			},
			errs: 2,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master", "master2"},
				},
				Worker: Worker{
					NodePools: []string{"worker"},
				},
				NodePools: NodePools{
					{
						Name:  "master",
						Count: 1,
					},
					{
						Name:  "worker",
						Count: 1,
					},
				},
			},
			errs: 1,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master"},
				},
				Worker: Worker{
					NodePools: []string{"worker"},
				},
				NodePools: NodePools{
					{
						Name:  "master",
						Count: 1,
					},
					{
						Name:  "worker",
						Count: 1,
					},
				},
			},
			errs: 0,
		},
	}

	for i, c := range cases {
		var n int
		errs := c.cluster.Validate()
		for _, err := range errs {
			if _, ok := err.(*ErrUnmatchedNodePool); ok {
				n++
			}
		}

		if n != c.errs {
			t.Errorf("test case %d: expected %d unmatched node pool errors, got %d", i, c.errs, n)
		}
	}
}

func TestSharedNodePool(t *testing.T) {
	cases := []struct {
		cluster Cluster
		errs    int
	}{
		{
			cluster: Cluster{},
			errs:    0,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master"},
				},
			},
			errs: 0,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"shared"},
				},
				Worker: Worker{
					NodePools: []string{"shared"},
				},
			},
			errs: 1,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"shared"},
				},
				Worker: Worker{
					NodePools: []string{"shared"},
				},
			},
			errs: 1,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"shared", "shared2"},
				},
				Worker: Worker{
					NodePools: []string{"shared", "shared2"},
				},
			},
			errs: 2,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"shared", "shared2"},
				},
				Worker: Worker{
					NodePools: []string{"shared", "shared2", "shared3"},
				},
			},
			errs: 2,
		},
	}

	for i, c := range cases {
		var n int
		errs := c.cluster.Validate()
		for _, err := range errs {
			if _, ok := err.(*ErrSharedNodePool); ok {
				n++
			}
		}

		if n != c.errs {
			t.Errorf("test case %d: expected %d shared node pool errors, got %d", i, c.errs, n)
		}
	}
}

func TestAWSEndpoints(t *testing.T) {
	cases := []struct {
		cluster Cluster
		err     bool
	}{
		{
			cluster: Cluster{},
			err:     true,
		},
		{
			cluster: defaultCluster,
			err:     false,
		},
		{
			cluster: Cluster{
				AWS: aws.AWS{
					Endpoints: "foo",
				},
			},
			err: true,
		},
		{
			cluster: Cluster{
				AWS: aws.AWS{
					Endpoints: aws.EndpointsAll,
				},
			},
			err: false,
		},
		{
			cluster: Cluster{
				AWS: aws.AWS{
					Endpoints: aws.EndpointsPrivate,
				},
			},
			err: false,
		},
		{
			cluster: Cluster{
				AWS: aws.AWS{
					Endpoints: aws.EndpointsPublic,
				},
			},
			err: false,
		},
	}

	for i, c := range cases {
		if err := c.cluster.validateAWSEndpoints(); (err != nil) != c.err {
			no := "no"
			if c.err {
				no = "an"
			}
			t.Errorf("test case %d: expected %s error, got %v", i, no, err)
		}
	}
}

func TestS3BucketNames(t *testing.T) {
	cases := []struct {
		cluster Cluster
		err     bool
	}{
		{
			cluster: defaultCluster,
			err:     true,
		},
		{
			cluster: Cluster{},
			err:     true,
		},
		{
			cluster: Cluster{
				Name:       "foo",
				BaseDomain: "example.com",
			},
			err: false,
		},
		{
			cluster: Cluster{
				Name:       ".foo",
				BaseDomain: "example.com",
			},
			err: true,
		},
		{
			cluster: Cluster{
				Name:       "foo",
				BaseDomain: "example.com.",
			},
			err: true,
		},
		{
			cluster: Cluster{
				Name:       "foo",
				BaseDomain: "012345678901234567890123456789012345678901234567890123456789.com",
			},
			err: true,
		},
	}

	for i, c := range cases {
		if err := c.cluster.validateS3Bucket(); (err != nil) != c.err {
			no := "no"
			if c.err {
				no = "an"
			}
			t.Errorf("test case %d: expected %s error, got %v", i, no, err)
		}
	}
}

func TestValidateLibvirt(t *testing.T) {
	cases := []struct {
		cluster Cluster
		err     bool
	}{
		{
			cluster: Cluster{},
			err:     true,
		},
		{
			cluster: defaultCluster,
			err:     true,
		},
		{
			cluster: Cluster{
				Libvirt: libvirt.Libvirt{
					Network: libvirt.Network{},
					Image:   "",
					URI:     "",
				},
				Networking: defaultCluster.Networking,
			},
			err: true,
		},
		{
			cluster: Cluster{
				Libvirt: libvirt.Libvirt{
					Network: libvirt.Network{
						Name:    "tectonic",
						IfName:  libvirt.DefaultIfName,
						IPRange: "10.0.1.0/24",
					},
					Image: "file:///foo",
					URI:   "baz",
				},
				Networking: defaultCluster.Networking,
			},
			err: false,
		},
		{
			cluster: Cluster{
				Libvirt: libvirt.Libvirt{
					Network: libvirt.Network{
						Name:    "tectonic",
						IfName:  libvirt.DefaultIfName,
						IPRange: "10.2.1.0/24",
					},
					Image: "file:///foo",
					URI:   "baz",
				},
				Networking: defaultCluster.Networking,
			},
			err: true,
		},
		{
			cluster: Cluster{
				Libvirt: libvirt.Libvirt{
					Network: libvirt.Network{
						Name:    "tectonic",
						IfName:  libvirt.DefaultIfName,
						IPRange: "x",
					},
					Image: "file:///foo",
					URI:   "baz",
				},
				Networking: defaultCluster.Networking,
			},
			err: true,
		},
		{
			cluster: Cluster{
				Libvirt: libvirt.Libvirt{
					Network: libvirt.Network{
						Name:    "tectonic",
						IfName:  libvirt.DefaultIfName,
						IPRange: "192.168.0.1/24",
					},
					Image: "file:///foo",
					URI:   "baz",
				},
				Networking: defaultCluster.Networking,
			},
			err: false,
		},
	}

	for i, c := range cases {
		c.cluster.Platform = PlatformLibvirt
		if err := c.cluster.validateLibvirt(); (err != nil) != c.err {
			no := "no"
			if c.err {
				no = "an"
			}
			t.Errorf("test case %d: expected %s error, got %v", i, no, err)
		}
	}
}

func TestValidateAWS(t *testing.T) {
	d1 := defaultCluster
	d1.Platform = PlatformAWS
	d2 := d1
	d2.Name = "test"
	d2.BaseDomain = "example.com"
	cases := []struct {
		cluster Cluster
		err     bool
	}{
		{
			cluster: Cluster{},
			err:     false,
		},
		{
			cluster: Cluster{
				Platform: PlatformAWS,
			},
			err: true,
		},
		{
			cluster: d1,
			err:     true,
		},
		{
			cluster: d2,
			err:     false,
		},
	}

	for i, c := range cases {
		if err := c.cluster.validateAWS(); (err != nil) != c.err {
			no := "no"
			if c.err {
				no = "an"
			}
			t.Errorf("test case %d: expected %s error, got %v", i, no, err)
		}
	}
}

func TestValidateOverlapWithPodOrServiceCIDR(t *testing.T) {
	cases := []struct {
		cidr    string
		cluster Cluster
		err     bool
	}{
		{
			cidr:    "192.168.0.1/24",
			cluster: Cluster{},
			err:     true,
		},
		{
			cidr:    "192.168.0.1/24",
			cluster: defaultCluster,
			err:     false,
		},
		{
			cidr:    "10.1.0.0/16",
			cluster: defaultCluster,
			err:     false,
		},
		{
			cidr:    "10.2.0.0/16",
			cluster: defaultCluster,
			err:     true,
		},
		{
			cidr: "10.1.0.0/16",
			cluster: Cluster{
				Networking: Networking{
					PodCIDR:     "10.3.0.0/16",
					ServiceCIDR: "10.4.0.0/16",
				},
			},
			err: false,
		},
		{
			cidr: "10.3.0.0/24",
			cluster: Cluster{
				Networking: Networking{
					PodCIDR:     "10.3.0.0/16",
					ServiceCIDR: "10.4.0.0/16",
				},
			},
			err: true,
		},
		{
			cidr: "0.0.0.0/0",
			cluster: Cluster{
				Networking: Networking{
					PodCIDR:     "10.3.0.0/16",
					ServiceCIDR: "10.4.0.0/16",
				},
			},
			err: true,
		},
	}

	for i, c := range cases {
		if err := c.cluster.validateOverlapWithPodOrServiceCIDR(c.cidr, "test"); (err != nil) != c.err {
			no := "no"
			if c.err {
				no = "an"
			}
			t.Errorf("test case %d: expected %s error, got %v", i, no, err)
		}
	}
}
