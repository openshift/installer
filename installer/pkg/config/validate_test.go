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
				Master: master{
					NodePools: []string{"", "", ""},
				},
			},
			errs: 3,
		},
		{
			cluster: Cluster{
				Master: master{
					NodePools: []string{"master"},
				},
			},
			errs: 2,
		},
		{
			cluster: Cluster{
				Master: master{
					NodePools: []string{"master"},
				},
				Worker: worker{
					NodePools: []string{"worker"},
				},
				Etcd: etcd{
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
				Master: master{
					NodePools: []string{"master"},
				},
			},
			errs: 0,
		},
		{
			cluster: Cluster{
				Master: master{
					NodePools: []string{"master"},
				},
				Worker: worker{
					NodePools: []string{"worker"},
				},
				Etcd: etcd{
					NodePools: []string{"etcd"},
				},
			},
			errs: 0,
		},
		{
			cluster: Cluster{
				Master: master{
					NodePools: []string{"master", "master2"},
				},
				Worker: worker{
					NodePools: []string{"worker"},
				},
				Etcd: etcd{
					NodePools: []string{"etcd"},
				},
			},
			errs: 1,
		},
		{
			cluster: Cluster{
				Master: master{
					NodePools: []string{"master", "master2"},
				},
				Worker: worker{
					NodePools: []string{"worker", "worker2"},
				},
				Etcd: etcd{
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
				Master: master{
					NodePools: []string{"master"},
				},
			},
			errs: 1,
		},
		{
			cluster: Cluster{
				Master: master{
					NodePools: []string{"master"},
				},
				Worker: worker{
					NodePools: []string{"worker"},
				},
				Etcd: etcd{
					NodePools: []string{"etcd"},
				},
			},
			errs: 3,
		},
		{
			cluster: Cluster{
				Master: master{
					NodePools: []string{"master", "master2"},
				},
				Worker: worker{
					NodePools: []string{"worker"},
				},
				Etcd: etcd{
					NodePools: []string{"etcd"},
				},
				NodePools: nodePools{
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
				Master: master{
					NodePools: []string{"master"},
				},
				Worker: worker{
					NodePools: []string{"worker"},
				},
				Etcd: etcd{
					NodePools: []string{"etcd"},
				},
				NodePools: nodePools{
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
				Master: master{
					NodePools: []string{"master"},
				},
			},
			errs: 0,
		},
		{
			cluster: Cluster{
				Master: master{
					NodePools: []string{"shared"},
				},
				Worker: worker{
					NodePools: []string{"shared"},
				},
				Etcd: etcd{
					NodePools: []string{"etcd"},
				},
			},
			errs: 1,
		},
		{
			cluster: Cluster{
				Master: master{
					NodePools: []string{"shared"},
				},
				Worker: worker{
					NodePools: []string{"shared"},
				},
				Etcd: etcd{
					NodePools: []string{"shared"},
				},
			},
			errs: 1,
		},
		{
			cluster: Cluster{
				Master: master{
					NodePools: []string{"shared", "shared2"},
				},
				Worker: worker{
					NodePools: []string{"shared", "shared2"},
				},
				Etcd: etcd{
					NodePools: []string{"shared"},
				},
			},
			errs: 2,
		},
		{
			cluster: Cluster{
				Master: master{
					NodePools: []string{"shared", "shared2"},
				},
				Worker: worker{
					NodePools: []string{"shared", "shared2", "shared3"},
				},
				Etcd: etcd{
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
		//if c.fields != 0 {
		//if len(sharedErrs) != 1 {
		//t.Errorf("test case %d: expected exactly one shared node pool error, got %d", i, len(sharedErrs))
		//continue
		//}
		//if c.fields != len(sharedErrs[0].fields) {
		//t.Errorf("test case %d: expected shared node pool error for %d fields, got %d", i, c.fields, len(sharedErrs[0].fields))
		//}
		//}
	}
}
