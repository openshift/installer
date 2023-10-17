/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ipam

import (
	"fmt"
	"net/netip"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apitypes "k8s.io/apimachinery/pkg/types"
	ipamv1 "sigs.k8s.io/cluster-api/exp/ipam/api/v1alpha1"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
)

// ipamDeviceConfig aids and holds state for the process
// of parsing IPAM addresses for a given device.
type ipamDeviceConfig struct {
	DeviceIndex         int
	IPAMAddresses       []*ipamv1.IPAddress
	MACAddress          string
	NetworkSpecGateway4 string
	IPAMConfigGateway4  string
	NetworkSpecGateway6 string
	IPAMConfigGateway6  string
}

func BuildState(ctx context.VMContext, networkStatus []infrav1.NetworkStatus) (map[string]infrav1.NetworkDeviceSpec, error) {
	state := map[string]infrav1.NetworkDeviceSpec{}

	ipamDeviceConfigs, err := buildIPAMDeviceConfigs(ctx, networkStatus)
	if err != nil {
		return state, err
	}

	var errs []error
	for _, ipamDeviceConfig := range ipamDeviceConfigs {
		var addressWithPrefixes []netip.Prefix
		for _, ipamAddress := range ipamDeviceConfig.IPAMAddresses {
			addressWithPrefix, err := parseAddressWithPrefix(ipamAddress)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			if slices.Contains(addressWithPrefixes, addressWithPrefix) {
				errs = append(errs,
					fmt.Errorf("IPAddress %s/%s is a duplicate of another address: %q",
						ipamAddress.Namespace,
						ipamAddress.Name,
						addressWithPrefix))
				continue
			}

			gatewayAddr, err := parseGateway(ipamAddress, addressWithPrefix, ipamDeviceConfig)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			if gatewayAddr != nil {
				if gatewayAddr.Is4() {
					ipamDeviceConfig.IPAMConfigGateway4 = ipamAddress.Spec.Gateway
				} else {
					ipamDeviceConfig.IPAMConfigGateway6 = ipamAddress.Spec.Gateway
				}
			}

			addressWithPrefixes = append(addressWithPrefixes, addressWithPrefix)
		}

		if len(addressWithPrefixes) > 0 {
			state[ipamDeviceConfig.MACAddress] = infrav1.NetworkDeviceSpec{
				IPAddrs:  prefixesAsStrings(addressWithPrefixes),
				Gateway4: ipamDeviceConfig.IPAMConfigGateway4,
				Gateway6: ipamDeviceConfig.IPAMConfigGateway6,
			}
		}
	}

	if len(errs) > 0 {
		var msgs []string
		for _, err := range errs {
			msgs = append(msgs, err.Error())
		}
		msg := strings.Join(msgs, "\n")
		return state, errors.New(msg)
	}
	return state, nil
}

// buildIPAMDeviceConfigs checks that all the IPAddressClaims have been satisfied.
// If each IPAddressClaim has an associated IPAddress, a slice of ipamDeviceConfig
// is returned, one for each device with addressesFromPools.
// If any of the IPAddressClaims do not have an associated IPAddress yet,
// a custom error is returned.
func buildIPAMDeviceConfigs(ctx context.VMContext, networkStatus []infrav1.NetworkStatus) ([]ipamDeviceConfig, error) {
	boundClaims, totalClaims := 0, 0
	ipamDeviceConfigs := []ipamDeviceConfig{}
	for devIdx, networkSpecDevice := range ctx.VSphereVM.Spec.Network.Devices {
		if len(networkStatus) == 0 ||
			len(networkStatus) <= devIdx ||
			networkStatus[devIdx].MACAddr == "" {
			return ipamDeviceConfigs, errors.New("waiting for devices to have MAC address set")
		}

		ipamDeviceConfig := ipamDeviceConfig{
			IPAMAddresses:       []*ipamv1.IPAddress{},
			MACAddress:          networkStatus[devIdx].MACAddr,
			NetworkSpecGateway4: networkSpecDevice.Gateway4,
			NetworkSpecGateway6: networkSpecDevice.Gateway6,
			DeviceIndex:         devIdx,
		}

		for poolRefIdx := range networkSpecDevice.AddressesFromPools {
			totalClaims++
			ipAddrClaimName := util.IPAddressClaimName(ctx.VSphereVM.Name, ipamDeviceConfig.DeviceIndex, poolRefIdx)
			ipAddrClaim, err := getIPAddrClaim(ctx, ipAddrClaimName)
			if err != nil {
				ctx.Logger.Error(err, "error fetching IPAddressClaim", "name", ipAddrClaimName)
				if apierrors.IsNotFound(err) {
					// it would be odd for this to occur, a findorcreate just happened in a previous step
					continue
				}
				return nil, err
			}

			ctx.Logger.V(5).Info("fetched IPAddressClaim", "name", ipAddrClaimName, "namespace", ctx.VSphereVM.Namespace)
			ipAddrName := ipAddrClaim.Status.AddressRef.Name
			if ipAddrName == "" {
				ctx.Logger.V(5).Info("IPAddress not yet bound to IPAddressClaim", "name", ipAddrClaimName, "namespace", ctx.VSphereVM.Namespace)
				continue
			}

			ipAddr := &ipamv1.IPAddress{}
			ipAddrKey := apitypes.NamespacedName{
				Namespace: ctx.VSphereVM.Namespace,
				Name:      ipAddrName,
			}
			if err := ctx.Client.Get(ctx, ipAddrKey, ipAddr); err != nil {
				// because the ref was set on the claim, it is expected this error should not occur
				return nil, err
			}
			ipamDeviceConfig.IPAMAddresses = append(ipamDeviceConfig.IPAMAddresses, ipAddr)
			boundClaims++
		}
		ipamDeviceConfigs = append(ipamDeviceConfigs, ipamDeviceConfig)
	}
	if boundClaims < totalClaims {
		ctx.Logger.Info("waiting for ip address claims to be bound",
			"total claims", totalClaims,
			"claims bound", boundClaims)
		return nil, ErrWaitingForIPAddr
	}
	return ipamDeviceConfigs, nil
}

// getIPAddrClaim fetches an IPAddressClaim from the api with the given name.
func getIPAddrClaim(ctx context.VMContext, ipAddrClaimName string) (*ipamv1.IPAddressClaim, error) {
	ipAddrClaim := &ipamv1.IPAddressClaim{}
	ipAddrClaimKey := apitypes.NamespacedName{
		Namespace: ctx.VSphereVM.Namespace,
		Name:      ipAddrClaimName,
	}

	ctx.Logger.V(5).Info("fetching IPAddressClaim", "name", ipAddrClaimKey.String())
	if err := ctx.Client.Get(ctx, ipAddrClaimKey, ipAddrClaim); err != nil {
		return nil, err
	}
	return ipAddrClaim, nil
}
