package clusterapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	capibm "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta3"
)

func TestRouteModeLoadBalancerIP(t *testing.T) {
	spec := []capibm.LoadBalancerSource{
		{Provision: capibm.LoadBalancerProvision{Name: "infra-loadbalancer"}},
		{Provision: capibm.LoadBalancerProvision{Name: "infra-loadbalancer-egress", RouteMode: capibm.PowerVSLoadBalancerRouteModeEnabled}},
	}

	for _, tc := range []struct {
		name   string
		spec   []capibm.LoadBalancerSource
		status []capibm.LoadBalancerStatus
		want   string
	}{
		{
			name: "route-mode LB IP is chosen over the public LB",
			spec: spec,
			status: []capibm.LoadBalancerStatus{
				{Name: "infra-loadbalancer", PrivateIPs: []string{"10.245.64.4"}},
				{Name: "infra-loadbalancer-egress", PrivateIPs: []string{"10.245.0.5", "10.245.0.6"}},
			},
			want: "10.245.0.5",
		},
		{
			name:   "route-mode LB has no IP yet",
			spec:   spec,
			status: []capibm.LoadBalancerStatus{{Name: "infra-loadbalancer", PrivateIPs: []string{"10.245.64.4"}}, {Name: "infra-loadbalancer-egress"}},
			want:   "",
		},
		{
			name:   "no route-mode LB declared",
			spec:   spec[:1],
			status: []capibm.LoadBalancerStatus{{Name: "infra-loadbalancer", PrivateIPs: []string{"10.245.64.4"}}},
			want:   "",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cluster := &capibm.IBMPowerVSCluster{}
			cluster.Spec.LoadBalancers = tc.spec
			cluster.Status.LoadBalancers = tc.status
			assert.Equal(t, tc.want, routeModeLoadBalancerIP(cluster))
		})
	}
}
