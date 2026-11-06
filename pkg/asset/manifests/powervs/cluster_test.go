package powervs

import (
	"testing"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
	capibm "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta3"

	"github.com/openshift/installer/pkg/types"
)

func TestBuildLoadBalancers(t *testing.T) {
	const infraID = "test-abcde"

	type lbExpect struct {
		name      string
		lbType    capibm.LoadBalancerType
		routeMode capibm.PowerVSLoadBalancerRouteMode
	}

	cases := []struct {
		name          string
		useDHCPServer bool
		publish       types.PublishingStrategy
		expected      []lbExpect
	}{
		{
			name:          "DHCPServer external: public LB and internal LB, no route-mode NLB",
			useDHCPServer: true,
			publish:       types.ExternalPublishingStrategy,
			expected: []lbExpect{
				{name: infraID + "-loadbalancer", lbType: capibm.LoadBalancerTypePublic},
				{name: infraID + "-loadbalancer-int", lbType: capibm.LoadBalancerTypePrivate},
			},
		},
		{
			name:    "DHCPSubnet external: public API LB first, then route-mode egress NLB, then internal LB",
			publish: types.ExternalPublishingStrategy,
			expected: []lbExpect{
				{name: infraID + "-loadbalancer", lbType: capibm.LoadBalancerTypePublic},
				{name: infraID + "-loadbalancer-egress", lbType: capibm.LoadBalancerTypePrivate, routeMode: capibm.PowerVSLoadBalancerRouteModeEnabled},
				{name: infraID + "-loadbalancer-int", lbType: capibm.LoadBalancerTypePrivate},
			},
		},
		{
			name:    "DHCPSubnet internal: route-mode NLB keeps the -loadbalancer name, no public LB",
			publish: types.InternalPublishingStrategy,
			expected: []lbExpect{
				{name: infraID + "-loadbalancer", lbType: capibm.LoadBalancerTypePrivate, routeMode: capibm.PowerVSLoadBalancerRouteModeEnabled},
				{name: infraID + "-loadbalancer-int", lbType: capibm.LoadBalancerTypePrivate},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildLoadBalancers(infraID, tc.useDHCPServer, tc.publish)

			actual := make([]lbExpect, 0, len(got))
			for _, lb := range got {
				assert.Equal(t, capibm.SourceTypeProvision, lb.Type)
				actual = append(actual, lbExpect{
					name:      lb.Provision.Name,
					lbType:    lb.Provision.Type,
					routeMode: lb.Provision.RouteMode,
				})

				// A route-mode LB must be a private network-fixed NLB; public LBs must not use route mode.
				if lb.Provision.RouteMode == capibm.PowerVSLoadBalancerRouteModeEnabled {
					assert.Equal(t, capibm.LoadBalancerTypePrivate, lb.Provision.Type)
					assert.Equal(t, capibm.PowerVSLoadBalancerProfileNetworkFixed, lb.Provision.Profile)
				}
			}
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestDHCPServerSnatPolicy(t *testing.T) {
	cases := []struct {
		name     string
		ic       *types.InstallConfig
		expected capibm.DHCPSnatPolicy
	}{
		{
			name:     "external, no mirrors: SNAT enabled",
			ic:       &types.InstallConfig{Publish: types.ExternalPublishingStrategy},
			expected: capibm.DHCPSnatPolicyEnabled,
		},
		{
			name: "external with mirrors: SNAT enabled",
			ic: &types.InstallConfig{
				Publish:            types.ExternalPublishingStrategy,
				ImageDigestSources: []types.ImageDigestSource{{Source: "quay.io/a", Mirrors: []string{"mirror.example.com/a"}}},
			},
			expected: capibm.DHCPSnatPolicyEnabled,
		},
		{
			name:     "internal, no mirrors: SNAT enabled",
			ic:       &types.InstallConfig{Publish: types.InternalPublishingStrategy},
			expected: capibm.DHCPSnatPolicyEnabled,
		},
		{
			name: "internal with ImageDigestSources: SNAT disabled",
			ic: &types.InstallConfig{
				Publish:            types.InternalPublishingStrategy,
				ImageDigestSources: []types.ImageDigestSource{{Source: "quay.io/a", Mirrors: []string{"mirror.example.com/a"}}},
			},
			expected: capibm.DHCPSnatPolicyDisabled,
		},
		{
			name: "internal with DeprecatedImageContentSources: SNAT disabled",
			ic: &types.InstallConfig{
				Publish:                       types.InternalPublishingStrategy,
				DeprecatedImageContentSources: []types.ImageContentSource{{Source: "quay.io/a", Mirrors: []string{"mirror.example.com/a"}}},
			},
			expected: capibm.DHCPSnatPolicyDisabled,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, dhcpServerSnatPolicy(tc.ic))
		})
	}
}

func TestExistingSubnetSource(t *testing.T) {
	id := "subnet-id"
	name := "subnet-name"

	got := existingSubnetSource(vpcv1.Subnet{ID: &id, Name: &name})

	assert.Equal(t, capibm.SourceTypeReference, got.Type)
	// CAPIBM requires exactly one of id or name on a ResourceIdentifier.
	assert.Equal(t, id, got.Reference.ID)
	assert.Empty(t, got.Reference.Name)
}

func TestBuildTransitGatewaySource(t *testing.T) {
	const infraID = "test-abcde"

	cases := []struct {
		name         string
		tgName       string
		isID         bool
		expectType   capibm.SourceType
		expectRef    capibm.ResourceIdentifier
		expectProv   capibm.TransitGatewayProvision
		expectNoProv bool
	}{
		{
			name:         "existing gateway ID is referenced by ID",
			tgName:       "0123abcd-ef01-2345-6789-abcdef012345",
			isID:         true,
			expectType:   capibm.SourceTypeReference,
			expectRef:    capibm.ResourceIdentifier{ID: "0123abcd-ef01-2345-6789-abcdef012345"},
			expectNoProv: true,
		},
		{
			name:       "blank name provisions a new gateway named after the infra ID",
			tgName:     "",
			expectType: capibm.SourceTypeProvision,
			expectProv: capibm.TransitGatewayProvision{Name: infraID + "-tg"},
		},
		{
			name:         "non-ID name is referenced by name, not provisioned",
			tgName:       "my-existing-tg",
			expectType:   capibm.SourceTypeReference,
			expectRef:    capibm.ResourceIdentifier{Name: "my-existing-tg"},
			expectNoProv: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildTransitGatewaySource(tc.tgName, tc.isID, infraID)

			assert.Equal(t, tc.expectType, got.Type)
			assert.Equal(t, tc.expectRef, got.Reference)
			if tc.expectNoProv {
				assert.Equal(t, capibm.TransitGatewayProvision{}, got.Provision)
			} else {
				assert.Equal(t, tc.expectProv, got.Provision)
			}
			// The VPC and PowerVS connections are always provisioned.
			assert.Equal(t, capibm.SourceTypeProvision, got.VPCConnection.Type)
			assert.Equal(t, capibm.SourceTypeProvision, got.PowerVSConnection.Type)
		})
	}
}

func TestPlanVPCEgress(t *testing.T) {
	const (
		infraID = "test-abcde"
		zone    = "us-south-1"
	)
	subnetIn := func(z string, attached bool) vpcv1.Subnet {
		s := vpcv1.Subnet{Zone: &vpcv1.ZoneReference{Name: ptr.To(z)}}
		if attached {
			s.PublicGateway = &vpcv1.PublicGatewayReference{ID: ptr.To("pgw-x")}
		}
		return s
	}
	pgw := func(id, name, z string) vpcv1.PublicGateway {
		return vpcv1.PublicGateway{ID: &id, Name: &name, Zone: &vpcv1.ZoneReference{Name: &z}}
	}
	rt := func(name string, tgIngress bool) vpcv1.RoutingTable {
		return vpcv1.RoutingTable{ID: ptr.To(name + "-id"), Name: ptr.To(name), RouteTransitGatewayIngress: ptr.To(tgIngress)}
	}

	cases := []struct {
		name       string
		subnets    []vpcv1.Subnet
		pgws       []vpcv1.PublicGateway
		rts        []vpcv1.RoutingTable
		expectErr  bool
		expectPGW  capibm.SourceType // type of the declared gateway
		expectAtt  capibm.ResourceIdentifier
		expectRT   bool
		expectWarn int
	}{
		{
			name:      "clean VPC: provision both",
			subnets:   []vpcv1.Subnet{subnetIn(zone, false)},
			rts:       []vpcv1.RoutingTable{rt("default", false)},
			expectPGW: capibm.SourceTypeProvision,
			expectAtt: capibm.ResourceIdentifier{Name: infraID + "-pgw"},
			expectRT:  true,
		},
		{
			name:      "user gateway in zone is referenced, not duplicated, and attached to bare subnets",
			subnets:   []vpcv1.Subnet{subnetIn(zone, false)},
			pgws:      []vpcv1.PublicGateway{pgw("pgw-user", "user-pgw", zone)},
			expectPGW: capibm.SourceTypeReference,
			expectAtt: capibm.ResourceIdentifier{ID: "pgw-user"},
			expectRT:  true,
			// reuse-and-attach warning
			expectWarn: 1,
		},
		{
			name:      "user gateway already attached to the zone subnet: reference without a warning",
			subnets:   []vpcv1.Subnet{subnetIn(zone, true)},
			pgws:      []vpcv1.PublicGateway{pgw("pgw-user", "user-pgw", zone)},
			expectPGW: capibm.SourceTypeReference,
			expectAtt: capibm.ResourceIdentifier{ID: "pgw-user"},
			expectRT:  true,
		},
		{
			name:      "gateway in another zone does not count",
			subnets:   []vpcv1.Subnet{subnetIn(zone, false)},
			pgws:      []vpcv1.PublicGateway{pgw("pgw-user", "user-pgw", "us-south-2")},
			expectPGW: capibm.SourceTypeProvision,
			expectAtt: capibm.ResourceIdentifier{Name: infraID + "-pgw"},
			expectRT:  true,
		},
		{
			name:      "gateway left by an earlier attempt keeps being provisioned by name",
			subnets:   []vpcv1.Subnet{subnetIn(zone, true)},
			pgws:      []vpcv1.PublicGateway{pgw("pgw-old", infraID+"-pgw", zone)},
			expectPGW: capibm.SourceTypeProvision,
			expectAtt: capibm.ResourceIdentifier{Name: infraID + "-pgw"},
			expectRT:  true,
		},
		{
			name:      "no gateway and no subnet in the zone: would orphan a gateway",
			subnets:   []vpcv1.Subnet{subnetIn("us-south-2", false)},
			expectErr: true,
		},
		{
			name:       "user routing table already owns Transit Gateway ingress: do not create ours",
			subnets:    []vpcv1.Subnet{subnetIn(zone, false)},
			rts:        []vpcv1.RoutingTable{rt("user-rt", true)},
			expectPGW:  capibm.SourceTypeProvision,
			expectAtt:  capibm.ResourceIdentifier{Name: infraID + "-pgw"},
			expectRT:   false,
			expectWarn: 1,
		},
		{
			name:      "routing table left by an earlier attempt keeps being provisioned by name",
			subnets:   []vpcv1.Subnet{subnetIn(zone, false)},
			rts:       []vpcv1.RoutingTable{rt(infraID+"-rt", true)},
			expectPGW: capibm.SourceTypeProvision,
			expectAtt: capibm.ResourceIdentifier{Name: infraID + "-pgw"},
			expectRT:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := planVPCEgress(infraID, zone, tc.subnets, tc.pgws, tc.rts)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			if assert.Len(t, got.publicGateways, 1) {
				assert.Equal(t, tc.expectPGW, got.publicGateways[0].Type)
				assert.Equal(t, zone, got.publicGateways[0].Zone)
			}
			assert.Equal(t, tc.expectAtt, got.attach)
			assert.Equal(t, tc.expectRT, len(got.routingTables) == 1)
			assert.Len(t, got.warnings, tc.expectWarn)
		})
	}
}

func TestNewVPCEgressPlan(t *testing.T) {
	got := newVPCEgressPlan("test-abcde", "us-south-1")

	if assert.Len(t, got.publicGateways, 1) {
		assert.Equal(t, capibm.SourceTypeProvision, got.publicGateways[0].Type)
		assert.Equal(t, "test-abcde-pgw", got.publicGateways[0].Provision.Name)
	}
	assert.Equal(t, capibm.ResourceIdentifier{Name: "test-abcde-pgw"}, got.attach)
	if assert.Len(t, got.routingTables, 1) {
		p := got.routingTables[0].Provision
		assert.Equal(t, "test-abcde-rt", p.Name)
		assert.Equal(t, capibm.RoutingTableIngressPolicyEnabled, p.RouteTransitGatewayIngress)
		if assert.Len(t, p.Routes, 1) {
			assert.Equal(t, "us-south-1", p.Routes[0].Zone)
			assert.Equal(t, "0.0.0.0/0", p.Routes[0].Destination)
		}
	}
}
