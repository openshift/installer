package ovirtclient

import (
	"fmt"
	"net"
)

func (m *mockClient) WaitForVMIPAddresses(
	id VMID,
	params VMIPSearchParams,
	retries ...RetryStrategy,
) (result map[string][]net.IP, err error) {
	return waitForIPAddresses(id, params, retries, m.logger, m)
}

func (o *oVirtClient) WaitForVMIPAddresses(
	id VMID,
	params VMIPSearchParams,
	retries ...RetryStrategy,
) (map[string][]net.IP, error) {
	return waitForIPAddresses(id, params, retries, o.logger, o)
}

var errNoIPAddressesReportedYet = newError(EPending, "no IP addresses reported yet")

func waitForIPAddresses(
	id VMID,
	params VMIPSearchParams,
	retries []RetryStrategy,
	logger Logger,
	client Client,
) (result map[string][]net.IP, err error) {
	retries = defaultRetries(retries, defaultLongTimeouts(client))
	hasNICs := false
	err = retry(
		fmt.Sprintf("waiting for IP addresses on VM %s", id),
		logger,
		retries,
		func() error {
			result, err = client.GetVMIPAddresses(id, params, retries...)
			if err != nil {
				return err
			}
			if len(result) == 0 {
				if !hasNICs {
					// We check if the VM has network interfaces and warn the user if that's not the case.
					hasNICs, err = noIPAddrCheckIfVMHasNICs(client, id, retries)
					if err != nil {
						// If a specific error was returned, return that, otherwise fall back on the normal
						// EPending below.
						return err
					}
				}
				return errNoIPAddressesReportedYet
			}
			return nil
		},
	)
	return result, err
}

func noIPAddrCheckIfVMHasNICs(client Client, id VMID, retries []RetryStrategy) (bool, error) {
	networkInterfaces, err := client.ListNICs(id, retries...)
	if err != nil {
		return false, errNoIPAddressesReportedYet
	}
	if len(networkInterfaces) == 0 {
		return false, wrap(
			fmt.Errorf("VM has no network interfaces"),
			EPending,
			"no IP addresses reported yet",
		)
	}
	return true, errNoIPAddressesReportedYet
}
