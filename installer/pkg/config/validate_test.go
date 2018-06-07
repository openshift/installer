package config

import (
	"os"
	"testing"

	"github.com/coreos/tectonic-installer/installer/pkg/config/aws"
)

func TestMissingNodePool(t *testing.T) {
	cases := []struct {
		cluster Cluster
		errs    int
	}{
		{
			cluster: Cluster{},
			errs:    3,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"", "", ""},
				},
			},
			errs: 3,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master"},
				},
			},
			errs: 2,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master"},
				},
				Worker: Worker{
					NodePools: []string{"worker"},
				},
				Etcd: Etcd{
					NodePools: []string{"etcd"},
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
				Etcd: Etcd{
					NodePools: []string{"etcd"},
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
				Etcd: Etcd{
					NodePools: []string{"etcd"},
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
				Etcd: Etcd{
					NodePools: []string{"etcd", "etcd2"},
				},
			},
			errs: 3,
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
				Etcd: Etcd{
					NodePools: []string{"etcd"},
				},
			},
			errs: 3,
		},
		{
			cluster: Cluster{
				Master: Master{
					NodePools: []string{"master", "master2"},
				},
				Worker: Worker{
					NodePools: []string{"worker"},
				},
				Etcd: Etcd{
					NodePools: []string{"etcd"},
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
					{
						Name:  "etcd",
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
				Etcd: Etcd{
					NodePools: []string{"etcd"},
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
					{
						Name:  "etcd",
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
				Etcd: Etcd{
					NodePools: []string{"etcd"},
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
				Etcd: Etcd{
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
				Etcd: Etcd{
					NodePools: []string{"shared"},
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
				Etcd: Etcd{
					NodePools: []string{"shared", "shared3"},
				},
			},
			errs: 3,
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

func TestTNCS3BucketNames(t *testing.T) {
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
				Platform:   PlatformAWS,
			},
			err: true,
		},
	}

	for i, c := range cases {
		if err := c.cluster.validateTNCS3Bucket(); (err != nil) != c.err {
			no := "no"
			if c.err {
				no = "an"
			}
			t.Errorf("test case %d: expected %s error, got %v", i, no, err)
		}
	}
}

func TestValidateIgnitionFiles(t *testing.T) {
	c := Cluster{
		NodePools: NodePools{
			{
				Name:         "error: invalid path",
				IgnitionFile: "do-not-exist.ign",
			},
			{
				Name:         "error: invalid config",
				IgnitionFile: "fixtures/invalid-ign.ign",
			},
			{
				Name: "ok: no field",
			},
			{
				Name:         "ok: empty field",
				IgnitionFile: "",
			},
			{
				Name:         "ok: valid config",
				IgnitionFile: "fixtures/ign.ign",
			},
		},
	}

	errs := c.validateIgnitionFiles()
	if len(errs) != 2 {
		t.Errorf("expected: %d ignition errors, got: %d", 2, len(errs))
	}
	if !os.IsNotExist(errs[0]) {
		t.Errorf("expected: notExistError, got: %v", errs[0])
	}
	if _, ok := errs[1].(*ErrInvalidIgnConfig); !ok {
		t.Errorf("expected: ErrInvalidIgnConfig, got: %v", errs[1])
	}
}

func TestValidateCL(t *testing.T) {
	cases := []struct {
		cluster Cluster
		err     bool
	}{
		{
			cluster: defaultCluster,
			err:     false,
		},
		{
			cluster: Cluster{},
			err:     true,
		},
		{
			cluster: Cluster{
				ContainerLinux: ContainerLinux{
					Channel: ContainerLinuxChannelBeta,
				},
			},
			err: true,
		},
		{
			cluster: Cluster{
				ContainerLinux: ContainerLinux{
					Channel: ContainerLinuxChannelBeta,
					Version: ContainerLinuxVersionLatest,
				},
			},
			err: false,
		},
		{
			cluster: Cluster{
				ContainerLinux: ContainerLinux{
					Channel: "foo",
					Version: ContainerLinuxVersionLatest,
				},
			},
			err: true,
		},
		{
			cluster: Cluster{
				ContainerLinux: ContainerLinux{
					Channel: ContainerLinuxChannelStable,
					Version: "100.99.98",
				},
			},
			err: false,
		},
		{
			cluster: Cluster{
				ContainerLinux: ContainerLinux{
					Channel: ContainerLinuxChannelStable,
					Version: "100..98",
				},
			},
			err: true,
		},
		{
			cluster: Cluster{
				ContainerLinux: ContainerLinux{
					Channel: ContainerLinuxChannelStable,
					Version: "100a99a98",
				},
			},
			err: true,
		},
		{
			cluster: Cluster{
				ContainerLinux: ContainerLinux{
					Channel: ContainerLinuxChannelStable,
					Version: "foo",
				},
			},
			err: true,
		},
	}

	for i, c := range cases {
		if err := c.cluster.validateCL(); (err != nil) != c.err {
			no := "no"
			if c.err {
				no = "an"
			}
			t.Errorf("test case %d: expected %s error, got %v", i, no, err)
		}
	}
}
