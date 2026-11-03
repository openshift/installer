package machines

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/vim25/soap"
)

type stubLookup struct {
	addrs []string
	err   error
}

func (s stubLookup) LookupHost(context.Context, string) ([]string, error) {
	return s.addrs, s.err
}

func TestVSphereCAPIPreflight(t *testing.T) {
	dnsErr := &net.DNSError{Err: "no such host", Name: "vcenter.bogus.invalid", IsNotFound: true}
	soapErr := soap.WrapSoapFault(&soap.Fault{
		Code:   "ServerFaultCode",
		String: "Cannot complete login due to an incorrect user name or password",
	})
	otherErr := errors.New("unable to retrieve compute cluster")

	cases := []struct {
		name         string
		server       string
		lookup       hostLookup
		loadNetworks func(context.Context) error
		wantErr      string
		wantNetworks bool
	}{
		{
			name:   "dns failure returns error instead of success",
			server: "vcenter.bogus.invalid",
			lookup: stubLookup{err: dnsErr},
			loadNetworks: func(context.Context) error {
				t.Fatal("Networks must not run when DNS fails")
				return nil
			},
			wantErr: `unable to resolve vSphere server vcenter.bogus.invalid`,
		},
		{
			name:   "soap auth failure returns error instead of success",
			server: "vcenter.example.com",
			lookup: stubLookup{addrs: []string{"192.0.2.1"}},
			loadNetworks: func(context.Context) error {
				return soapErr
			},
			wantErr:      `authentication failure to vCenter vcenter.example.com`,
			wantNetworks: true,
		},
		{
			name:   "non-soap networks error is returned unchanged",
			server: "vcenter.example.com",
			lookup: stubLookup{addrs: []string{"192.0.2.1"}},
			loadNetworks: func(context.Context) error {
				return otherErr
			},
			wantErr:      `unable to retrieve compute cluster`,
			wantNetworks: true,
		},
		{
			name:   "success loads networks",
			server: "vcenter.example.com",
			lookup: stubLookup{addrs: []string{"192.0.2.1"}},
			loadNetworks: func(context.Context) error {
				return nil
			},
			wantNetworks: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			networksCalled := false
			err := vsphereCAPIPreflight(context.Background(), tc.lookup, tc.server, func(ctx context.Context) error {
				networksCalled = true
				return tc.loadNetworks(ctx)
			})
			if tc.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Regexp(t, tc.wantErr, err)
			}
			assert.Equal(t, tc.wantNetworks, networksCalled)
		})
	}
}

func TestVSphereCAPIPreflightNetworksTimeoutNotNestedUnderDNS(t *testing.T) {
	var gotRemaining time.Duration
	err := vsphereCAPIPreflight(context.Background(), stubLookup{addrs: []string{"192.0.2.1"}}, "vcenter.example.com", func(ctx context.Context) error {
		deadline, ok := ctx.Deadline()
		if !ok {
			t.Fatal("expected Networks context to have a deadline")
		}
		gotRemaining = time.Until(deadline)
		return nil
	})
	assert.NoError(t, err)
	// Nested under the 30s DNS context, remaining would be ~30s. Independent
	// 60s budget must still be close to vsphereNetworksTimeout.
	assert.Greater(t, gotRemaining, 45*time.Second)
	assert.Less(t, gotRemaining, vsphereNetworksTimeout+time.Second)
}
