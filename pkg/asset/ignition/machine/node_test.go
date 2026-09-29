package machine

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/powervc"
)

// hostFromPointerConfig is a helper that calls pointerIgnitionConfig and
// returns the host portion of the ignition merge URL.
func hostFromPointerConfig(t *testing.T, ic *types.InstallConfig, role string) string {
	t.Helper()
	cfg := pointerIgnitionConfig(ic, []byte("fake-ca"), role)
	require.NotEmpty(t, cfg.Ignition.Config.Merge, "expected at least one merge entry")
	src := cfg.Ignition.Config.Merge[0].Source
	require.NotNil(t, src, "expected merge source to be set")
	u, err := url.Parse(*src)
	require.NoError(t, err, "unexpected error parsing ignition source URL")
	return u.Host
}

// baseOpenStackIC returns a minimal OpenStack InstallConfig for use in tests.
func baseOpenStackIC(lbType v1.PlatformLoadBalancerType, dnsRecordsType v1.DNSRecordsType) *types.InstallConfig {
	return &types.InstallConfig{
		ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
		BaseDomain: "test-domain",
		Networking: &types.Networking{
			ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("10.0.1.0/24")},
		},
		Platform: types.Platform{
			OpenStack: &openstack.Platform{
				APIVIPs: []string{"1.2.3.4"},
				LoadBalancer: &v1.OpenStackPlatformLoadBalancer{
					Type: lbType,
				},
				DNSRecordsType: dnsRecordsType,
			},
		},
	}
}

// basePowerVCIC returns a minimal PowerVC InstallConfig for use in tests.
func basePowerVCIC(lbType v1.PlatformLoadBalancerType) *types.InstallConfig {
	p := &powervc.Platform{
		Cloud:        "test-cloud",
		APIVIPs:      []string{"5.6.7.8"},
		LoadBalancer: &v1.OpenStackPlatformLoadBalancer{Type: lbType},
	}
	return &types.InstallConfig{
		ObjectMeta: metav1.ObjectMeta{Name: "test-cluster"},
		BaseDomain: "test-domain",
		Networking: &types.Networking{
			ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("10.0.1.0/24")},
		},
		Platform: types.Platform{PowerVC: p},
	}
}

// TestPointerIgnitionConfigOpenStack tests that the ignition host is set
// correctly for OpenStack depending on the LoadBalancer type and DNSRecordsType,
// for both the master and worker roles.
func TestPointerIgnitionConfigOpenStack(t *testing.T) {
	cases := []struct {
		name           string
		lbType         v1.PlatformLoadBalancerType
		dnsRecordsType v1.DNSRecordsType
		expectedHost   string
	}{
		{
			name:           "user-managed LB with external DNS uses api-int FQDN",
			lbType:         v1.LoadBalancerTypeUserManaged,
			dnsRecordsType: v1.DNSRecordsTypeExternal,
			expectedHost:   "api-int.test-cluster.test-domain:22623",
		},
		{
			name:           "openshift-managed LB with internal DNS uses VIP",
			lbType:         v1.LoadBalancerTypeOpenShiftManagedDefault,
			dnsRecordsType: v1.DNSRecordsTypeInternal,
			expectedHost:   "1.2.3.4:22623",
		},
		{
			name:           "user-managed LB with internal DNS uses VIP",
			lbType:         v1.LoadBalancerTypeUserManaged,
			dnsRecordsType: v1.DNSRecordsTypeInternal,
			expectedHost:   "1.2.3.4:22623",
		},
	}

	for _, role := range []string{"master", "worker"} {
		for _, tc := range cases {
			t.Run(role+"/"+tc.name, func(t *testing.T) {
				ic := baseOpenStackIC(tc.lbType, tc.dnsRecordsType)
				host := hostFromPointerConfig(t, ic, role)
				assert.Equal(t, tc.expectedHost, host)
			})
		}
	}
}

// TestPointerIgnitionConfigPowerVC tests that the ignition host is always set
// to the API VIP for PowerVC, regardless of the LoadBalancer type, and for
// both the master and worker roles.
// PowerVC has no DNSRecordsType field, so the externally managed DNS
// condition that unlocks FQDN mode on OpenStack can never be met
// here. The VIP must always be used.
func TestPointerIgnitionConfigPowerVC(t *testing.T) {
	userManaged := v1.LoadBalancerTypeUserManaged
	openshiftManaged := v1.LoadBalancerTypeOpenShiftManagedDefault

	cases := []struct {
		name   string
		lbType v1.PlatformLoadBalancerType
	}{
		{
			name:   "openshift-managed LB uses VIP",
			lbType: openshiftManaged,
		},
		{
			name:   "user-managed LB uses VIP",
			lbType: userManaged,
		},
	}

	for _, role := range []string{"master", "worker"} {
		for _, tc := range cases {
			t.Run(role+"/"+tc.name, func(t *testing.T) {
				ic := basePowerVCIC(tc.lbType)
				host := hostFromPointerConfig(t, ic, role)
				assert.Equal(t, "5.6.7.8:22623", host)
			})
		}
	}
}
