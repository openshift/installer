package config

import "testing"

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
