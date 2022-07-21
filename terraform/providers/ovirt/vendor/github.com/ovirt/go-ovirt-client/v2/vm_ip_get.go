package ovirtclient

import (
	"fmt"
	"net"
	"regexp"

	ovirtsdk "github.com/ovirt/go-ovirt"
)

func (m *mockClient) GetVMIPAddresses(id VMID, params VMIPSearchParams, _ ...RetryStrategy) (map[string][]net.IP, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.vms[id]; !ok {
		return nil, newError(ENotFound, "VM %s not found", id)
	}

	source := m.vmIPs[id]

	return filterReportedIPList(source, params), nil
}

func (o *oVirtClient) GetVMIPAddresses(id VMID, params VMIPSearchParams, retries ...RetryStrategy) (result map[string][]net.IP, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	result = map[string][]net.IP{}
	err = retry(
		fmt.Sprintf("getting IP addresses for VM %s", id),
		o.logger,
		retries,
		func() error {
			reportedDevicesResponse, err := o.conn.SystemService().VmsService().VmService(string(id)).ReportedDevicesService().List().Send()

			reportedDevices, ok := reportedDevicesResponse.ReportedDevice()
			if !ok {
				return newFieldNotFound("reported devices response", "reported device")
			}

			for _, reportedDevice := range reportedDevices.Slice() {
				reportedDeviceType, ok := reportedDevice.Type()
				if !ok {
					continue
				}
				reportedDeviceName, ok := reportedDevice.Name()
				if !ok {
					continue
				}
				reportedDeviceIPs, ok := reportedDevice.Ips()
				if !ok {
					continue
				}
				if reportedDeviceType == ovirtsdk.REPORTEDDEVICETYPE_NETWORK {
					result[reportedDeviceName] = []net.IP{}
					for _, reportedDeviceIP := range reportedDeviceIPs.Slice() {
						ip, ok := reportedDeviceIP.Address()
						if !ok {
							continue
						}
						result[reportedDeviceName] = append(result[reportedDeviceName], net.ParseIP(ip))
					}
				}
			}

			return err
		})
	return filterReportedIPList(result, params), err
}

func filterReportedIPList(source map[string][]net.IP, params VMIPSearchParams) map[string][]net.IP {
	if source == nil {
		return nil
	}
	if params == nil {
		params = NewVMIPSearchParams()
	}

	result := map[string][]net.IP{}
	includedRanges := params.GetIncludedRanges()
	excludedRanges := params.GetExcludedRanges()
	includedInterfaceNames := params.GetIncludedInterfaces()
	excludedInterfaceNames := params.GetExcludedInterfaces()
	includedInterfacePatterns := params.GetIncludedInterfacePatterns()
	excludedInterfacePatterns := params.GetExcludedInterfacePatterns()

	for interf, ips := range source {
		if excludeInterface(
			interf,
			excludedInterfaceNames,
			excludedInterfacePatterns,
			includedInterfaceNames,
			includedInterfacePatterns,
		) {
			continue
		}

		ipList := gatherInterfaceIPs(ips, excludedRanges, includedRanges)
		if len(ipList) > 0 {
			result[interf] = ipList
		}
	}
	return result
}

func gatherInterfaceIPs(ips []net.IP, excludedRanges []net.IPNet, includedRanges []net.IPNet) []net.IP {
	var ipList []net.IP
	for _, ip := range ips {
		valid := true
		for _, excludedRange := range excludedRanges {
			if excludedRange.Contains(ip) {
				valid = false
				break
			}
		}
		if valid && len(includedRanges) > 0 {
			valid = false
			for _, includedRange := range includedRanges {
				if includedRange.Contains(ip) {
					valid = true
					break
				}
			}
		}
		if valid {
			ipList = append(ipList, ip)
		}
	}
	return ipList
}

func excludeInterface(
	interf string,
	excludedInterfaceNames []string,
	excludedInterfacePatterns []*regexp.Regexp,
	includedInterfaceNames []string,
	includedInterfacePatterns []*regexp.Regexp,
) bool {
	if interfaceIsExcludedByName(
		excludedInterfaceNames,
		interf,
	) || interfaceIsExcludedByPattern(
		excludedInterfacePatterns,
		interf,
	) {
		return true
	}
	includedByName := isInIncludedInterfaceNames(includedInterfaceNames, interf)
	includedByPattern := isInIncludedInterfacePatterns(includedInterfacePatterns, interf)
	return !((len(includedInterfaceNames) == 0 && len(includedInterfacePatterns) == 0) || includedByName || includedByPattern)
}

func isInIncludedInterfacePatterns(patterns []*regexp.Regexp, interf string) bool {
	for _, pattern := range patterns {
		if pattern.MatchString(interf) {
			return true
		}
	}
	return false
}

func isInIncludedInterfaceNames(includedInterfaceNames []string, interf string) bool {
	for _, interf2 := range includedInterfaceNames {
		if interf2 == interf {
			return true
		}
	}
	return false
}

func interfaceIsExcludedByPattern(excludedInterfacePatterns []*regexp.Regexp, interf string) bool {
	skip := false
	for _, p := range excludedInterfacePatterns {
		if p.MatchString(interf) {
			skip = true
			break
		}
	}
	return skip
}

func interfaceIsExcludedByName(excludedInterfaceNames []string, interf string) bool {
	skip := false
	for _, n := range excludedInterfaceNames {
		if n == interf {
			skip = true
			break
		}
	}
	return skip
}
