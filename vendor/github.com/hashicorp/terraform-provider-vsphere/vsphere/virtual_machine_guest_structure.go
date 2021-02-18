package vsphere

import (
	"log"
	"net"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/govmomi/vim25/types"
)

// schemaVirtualMachineGuestInfo returns schema items for the relevant parts of
// GuestInfo that vsphere_virtual_machine tracks (mostly guest information).
func schemaVirtualMachineGuestInfo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"default_ip_address": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The IP address selected by Terraform to be used for the provisioner.",
		},
		"guest_ip_addresses": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "The current list of IP addresses on this virtual machine.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

type int32arr []int32

func (a int32arr) Len() int           { return len(a) }
func (a int32arr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a int32arr) Less(i, j int) bool { return a[i] < a[j] }

// buildAndSelectGuestIPs builds a list of IP addresses known to VMware tools.
// From this list, it selects the first IP address it seems that's associated
// with a default gateway - first IPv4, and then IPv6 if criteria can't be
// satisfied - and sets that as the default_ip_address and also the IP address
// used for provisioning. The full list of IP addresses is saved to
// guest_ip_addresses.
func buildAndSelectGuestIPs(d *schema.ResourceData, guest types.GuestInfo) error {
	log.Printf("[DEBUG] %s: Checking guest networking state", resourceVSphereVirtualMachineIDString(d))
	var v4primary, v6primary, v4gw, v6gw net.IP
	var v4net2addrs, v6net2addrs map[int32][]string
	var deviceConfigIds int32arr

	// Fetch gateways first.
	for _, s := range guest.IpStack {
		if s.IpRouteConfig != nil {
			for _, r := range s.IpRouteConfig.IpRoute {
				switch r.Network {
				case "0.0.0.0":
					v4gw = net.ParseIP(r.Gateway.IpAddress)
				case "::":
					v6gw = net.ParseIP(r.Gateway.IpAddress)
				}
			}
		}
	}

	addrs := make([]string, 0)
	v4net2addrs = make(map[int32][]string)
	v6net2addrs = make(map[int32][]string)

	// Now fetch all IP addresses, checking at the same time to see if the IP
	// address is eligible to be a primary IP address.
	for _, n := range guest.Net {
		if n.IpConfig != nil {
			deviceConfigIds = append(deviceConfigIds, n.DeviceConfigId)
			v4net2addrs[n.DeviceConfigId] = make([]string, 0)
			v6net2addrs[n.DeviceConfigId] = make([]string, 0)
			for _, addr := range n.IpConfig.IpAddress {
				ip := net.ParseIP(addr.IpAddress)
				var mask net.IPMask
				if ip.To4() != nil {
					v4net2addrs[n.DeviceConfigId] = append(v4net2addrs[n.DeviceConfigId], addr.IpAddress)
					mask = net.CIDRMask(int(addr.PrefixLength), 32)
					if ip.Mask(mask).Equal(v4gw.Mask(mask)) && v4primary == nil {
						v4primary = ip
					}
				} else {
					v6net2addrs[n.DeviceConfigId] = append(v6net2addrs[n.DeviceConfigId], addr.IpAddress)
					mask = net.CIDRMask(int(addr.PrefixLength), 128)
					if ip.Mask(mask).Equal(v6gw.Mask(mask)) && v6primary == nil {
						v6primary = ip
					}
				}
			}
		}
	}

	sort.Sort(deviceConfigIds)
	for _, deviceConfigId := range deviceConfigIds {
		addrs = append(addrs, v4net2addrs[deviceConfigId]...)
		addrs = append(addrs, v6net2addrs[deviceConfigId]...)
	}

	// Fall back to the IpAddress property in GuestInfo directly when the
	// IpStack and Net properties are not populated. This generally means that
	// an older version of VMTools is in use.
	if len(addrs) < 1 && guest.IpAddress != "" {
		addrs = append(addrs, guest.IpAddress)
	}

	if len(addrs) < 1 {
		// No IP addresses were discovered. This more than likely means that the VM
		// is powered off, or VMware tools is not installed. We can return here,
		// setting the empty set of addresses to avoid spurious diffs.
		log.Printf("[DEBUG] %s: No IP addresses found in guest state", resourceVSphereVirtualMachineIDString(d))
		return d.Set("guest_ip_addresses", addrs)
	}
	var primary string
	switch {
	case v4primary != nil:
		primary = v4primary.String()
	case v6primary != nil:
		primary = v6primary.String()
	default:
		primary = addrs[0]
	}
	log.Printf("[DEBUG] %s: Primary IP address: %s", resourceVSphereVirtualMachineIDString(d), primary)
	d.Set("default_ip_address", primary)
	log.Printf("[DEBUG] %s: All IP addresses: %s", resourceVSphereVirtualMachineIDString(d), strings.Join(addrs, ","))
	if err := d.Set("guest_ip_addresses", addrs); err != nil {
		return err
	}
	d.SetConnInfo(map[string]string{
		"type": "ssh",
		"host": primary,
	})

	return nil
}
